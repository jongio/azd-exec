package executor

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

// Key Vault reference patterns.
var (
	// Pattern: @Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/name[/version]).
	kvRefSecretURIPattern = regexp.MustCompile(`^@Microsoft\.KeyVault\(SecretUri=(.+)\)$`)

	// Pattern: @Microsoft.KeyVault(VaultName=vault;SecretName=name[;SecretVersion=version]).
	kvRefVaultNamePattern = regexp.MustCompile(`^@Microsoft\.KeyVault\(VaultName=([^;]+);SecretName=([^;)]+)(?:;SecretVersion=([^;)]+))?\)$`)

	// Pattern: akvs://<guid>/<vault>/<secret>[/<version>]
	// Note: guid is informational; <vault> is used to construct https://<vault>.vault.azure.net
	kvRefAzdAkvsPattern = regexp.MustCompile(`^akvs://([^/]+)/([^/]+)/([^/]+)(?:/([^/]+))?$`)
)

func normalizeKeyVaultReferenceValue(value string) string {
	normalized := strings.TrimSpace(value)
	if len(normalized) < 2 {
		return normalized
	}

	first := normalized[0]
	last := normalized[len(normalized)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		// Strip wrapper quotes only when they wrap the entire value.
		normalized = strings.TrimSpace(normalized[1 : len(normalized)-1])
	}

	return normalized
}

// KeyVaultResolver resolves Key Vault references in environment variables.
type KeyVaultResolver struct {
	credential *azidentity.DefaultAzureCredential
	clients    map[string]*azsecrets.Client
	mu         sync.RWMutex // protects clients map
}

// ResolveEnvironmentOptions controls how environment variable resolution behaves.
type ResolveEnvironmentOptions struct {
	// StopOnError causes ResolveEnvironmentVariables to stop at the first resolution error.
	// When true, it returns a non-nil error and does not return partial results.
	StopOnError bool
}

// KeyVaultResolutionWarning represents a non-fatal resolution failure for a single environment variable.
// The secret value is never included.
type KeyVaultResolutionWarning struct {
	Key string
	Err error
}

// NewKeyVaultResolver creates a new Key Vault resolver.
func NewKeyVaultResolver() (*KeyVaultResolver, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	return &KeyVaultResolver{
		credential: cred,
		clients:    make(map[string]*azsecrets.Client),
	}, nil
}

// IsKeyVaultReference checks if a value is a Key Vault reference.
func IsKeyVaultReference(value string) bool {
	normalized := normalizeKeyVaultReferenceValue(value)

	if strings.HasPrefix(normalized, "akvs://") {
		return true
	}

	return strings.HasPrefix(normalized, "@Microsoft.KeyVault(") && strings.HasSuffix(normalized, ")")
}

// ResolveReference resolves a single Key Vault reference to its secret value.
func (r *KeyVaultResolver) ResolveReference(ctx context.Context, reference string) (string, error) {
	reference = normalizeKeyVaultReferenceValue(reference)

	// Try SecretUri pattern first
	if matches := kvRefSecretURIPattern.FindStringSubmatch(reference); matches != nil {
		secretURI := strings.TrimSpace(matches[1])
		return r.resolveBySecretURI(ctx, secretURI)
	}

	// Try VaultName pattern
	if matches := kvRefVaultNamePattern.FindStringSubmatch(reference); matches != nil {
		vaultName := matches[1]
		secretName := matches[2]
		version := ""
		if len(matches) > 3 {
			version = matches[3]
		}
		return r.resolveByVaultNameAndSecret(ctx, vaultName, secretName, version)
	}

	// Try azd akvs format: akvs://<guid>/<vault>/<secret>[/<version>]
	if strings.HasPrefix(reference, "akvs://") {
		if !kvRefAzdAkvsPattern.MatchString(reference) {
			return "", fmt.Errorf("invalid akvs URI format")
		}

		guid, vaultName, secretName, version, err := parseAzdAkvsURI(reference)
		_ = guid // informational only
		if err != nil {
			return "", err
		}
		return r.resolveByVaultNameAndSecret(ctx, vaultName, secretName, version)
	}

	// Avoid returning the full reference value (could include sensitive context).
	return "", fmt.Errorf("invalid Key Vault reference format")
}

func parseAzdAkvsURI(akvsURI string) (guid, vaultName, secretName, version string, err error) {
	parsed, err := url.Parse(akvsURI)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid akvs URI: %w", err)
	}

	if parsed.Scheme != "akvs" {
		return "", "", "", "", fmt.Errorf("invalid akvs URI scheme")
	}

	guid = parsed.Host
	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", "", "", fmt.Errorf("invalid akvs URI path")
	}
	if len(pathParts) > 3 {
		return "", "", "", "", fmt.Errorf("invalid akvs URI path")
	}

	vaultName = pathParts[0]
	secretName = pathParts[1]
	if len(pathParts) == 3 {
		version = pathParts[2]
	}

	if guid == "" || vaultName == "" || secretName == "" {
		return "", "", "", "", fmt.Errorf("invalid akvs URI")
	}

	return guid, vaultName, secretName, version, nil
}

// resolveBySecretURI resolves a secret using a full secret URI.
func (r *KeyVaultResolver) resolveBySecretURI(ctx context.Context, secretURI string) (string, error) {
	// Parse the secret URI to extract vault URL and secret details
	parsedURI, err := url.Parse(secretURI)
	if err != nil {
		return "", fmt.Errorf("invalid secret URI: %w", err)
	}

	// Extract vault URL (scheme + host)
	vaultURL := fmt.Sprintf("%s://%s", parsedURI.Scheme, parsedURI.Host)

	// Get or create client for this vault
	client, err := r.getClient(vaultURL)
	if err != nil {
		return "", err
	}

	// Parse path to get secret name and optional version
	// Path format: /secrets/{secret-name}[/{version}]
	pathParts := strings.Split(strings.Trim(parsedURI.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "secrets" {
		return "", fmt.Errorf("invalid secret URI path: %s", parsedURI.Path)
	}

	secretName := pathParts[1]
	version := ""
	if len(pathParts) > 2 {
		version = pathParts[2]
	}

	// Fetch the secret
	return r.getSecretValue(ctx, client, secretName, version)
}

// resolveByVaultNameAndSecret resolves a secret using vault name and secret name.
func (r *KeyVaultResolver) resolveByVaultNameAndSecret(ctx context.Context, vaultName, secretName, version string) (string, error) {
	// Construct vault URL
	vaultURL := fmt.Sprintf("https://%s.vault.azure.net", vaultName)

	// Get or create client for this vault
	client, err := r.getClient(vaultURL)
	if err != nil {
		return "", err
	}

	// Fetch the secret
	return r.getSecretValue(ctx, client, secretName, version)
}

// getClient gets or creates a Key Vault client for the specified vault URL.
func (r *KeyVaultResolver) getClient(vaultURL string) (*azsecrets.Client, error) {
	// Check if client exists (read lock)
	r.mu.RLock()
	if client, ok := r.clients[vaultURL]; ok {
		r.mu.RUnlock()
		return client, nil
	}
	r.mu.RUnlock()

	// Create new client (write lock)
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine may have created it)
	if client, ok := r.clients[vaultURL]; ok {
		return client, nil
	}

	client, err := azsecrets.NewClient(vaultURL, r.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Vault client for %s: %w", vaultURL, err)
	}

	r.clients[vaultURL] = client
	return client, nil
}

// getSecretValue fetches the actual secret value from Key Vault.
func (r *KeyVaultResolver) getSecretValue(ctx context.Context, client *azsecrets.Client, secretName, version string) (string, error) {
	var resp azsecrets.GetSecretResponse
	var err error

	if version != "" {
		resp, err = client.GetSecret(ctx, secretName, version, nil)
	} else {
		resp, err = client.GetSecret(ctx, secretName, "", nil)
	}

	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("secret %s has no value", secretName)
	}

	return *resp.Value, nil
}

// ResolveEnvironmentVariables resolves all Key Vault references in the provided environment variables.
func (r *KeyVaultResolver) ResolveEnvironmentVariables(ctx context.Context, envVars []string, options ResolveEnvironmentOptions) ([]string, []KeyVaultResolutionWarning, error) {
	return resolveEnvironmentVariables(ctx, envVars, r.ResolveReference, options)
}

func resolveEnvironmentVariables(
	ctx context.Context,
	envVars []string,
	resolve func(context.Context, string) (string, error),
	options ResolveEnvironmentOptions,
) ([]string, []KeyVaultResolutionWarning, error) {
	resolved := make([]string, 0, len(envVars))
	warnings := make([]KeyVaultResolutionWarning, 0)

	for _, envVar := range envVars {
		// Split into key=value
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			resolved = append(resolved, envVar)
			continue
		}

		key := parts[0]
		value := parts[1]

		if !IsKeyVaultReference(value) {
			resolved = append(resolved, envVar)
			continue
		}

		secretValue, err := resolve(ctx, value)
		if err != nil {
			warnings = append(warnings, KeyVaultResolutionWarning{Key: key, Err: err})
			if options.StopOnError {
				return nil, warnings, fmt.Errorf("failed to resolve Key Vault reference for %s: %w", key, err)
			}
			// Keep original (unresolved) reference value and continue.
			resolved = append(resolved, envVar)
			continue
		}

		resolved = append(resolved, fmt.Sprintf("%s=%s", key, secretValue))
	}

	return resolved, warnings, nil
}
