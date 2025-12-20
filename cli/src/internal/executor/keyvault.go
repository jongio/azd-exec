package executor

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

// Key Vault reference patterns
var (
	// Pattern: @Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/name[/version])
	kvRefSecretURIPattern = regexp.MustCompile(`^@Microsoft\.KeyVault\(SecretUri=(.+)\)$`)
	
	// Pattern: @Microsoft.KeyVault(VaultName=vault;SecretName=name[;SecretVersion=version])
	kvRefVaultNamePattern = regexp.MustCompile(`^@Microsoft\.KeyVault\(VaultName=([^;]+);SecretName=([^;)]+)(?:;SecretVersion=([^;)]+))?\)$`)
)

// KeyVaultResolver resolves Key Vault references in environment variables.
type KeyVaultResolver struct {
	credential *azidentity.DefaultAzureCredential
	clients    map[string]*azsecrets.Client
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
	return strings.HasPrefix(value, "@Microsoft.KeyVault(") && strings.HasSuffix(value, ")")
}

// ResolveReference resolves a single Key Vault reference to its secret value.
func (r *KeyVaultResolver) ResolveReference(ctx context.Context, reference string) (string, error) {
	// Try SecretUri pattern first
	if matches := kvRefSecretURIPattern.FindStringSubmatch(reference); matches != nil {
		secretURI := matches[1]
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

	return "", fmt.Errorf("invalid Key Vault reference format: %s", reference)
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
func (r *KeyVaultResolver) ResolveEnvironmentVariables(ctx context.Context, envVars []string) ([]string, error) {
	resolved := make([]string, 0, len(envVars))

	for _, envVar := range envVars {
		// Split into key=value
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			resolved = append(resolved, envVar)
			continue
		}

		key := parts[0]
		value := parts[1]

		// Check if value is a Key Vault reference
		if IsKeyVaultReference(value) {
			secretValue, err := r.ResolveReference(ctx, value)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve Key Vault reference for %s: %w", key, err)
			}
			resolved = append(resolved, fmt.Sprintf("%s=%s", key, secretValue))
		} else {
			resolved = append(resolved, envVar)
		}
	}

	return resolved, nil
}
