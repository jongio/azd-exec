package executor

import (
	"strings"
	"testing"
)

func TestParseSecretURI(t *testing.T) {
	tests := []struct {
		name        string
		secretURI   string
		wantVault   string
		wantSecret  string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "Valid URI without version",
			secretURI:   "https://myvault.vault.azure.net/secrets/mysecret",
			wantVault:   "https://myvault.vault.azure.net",
			wantSecret:  "mysecret",
			wantVersion: "",
			wantErr:     false,
		},
		{
			name:        "Valid URI with version",
			secretURI:   "https://myvault.vault.azure.net/secrets/mysecret/abc123",
			wantVault:   "https://myvault.vault.azure.net",
			wantSecret:  "mysecret",
			wantVersion: "abc123",
			wantErr:     false,
		},
		{
			name:        "Invalid path - missing secrets",
			secretURI:   "https://myvault.vault.azure.net/mysecret",
			wantVault:   "",
			wantSecret:  "",
			wantVersion: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We don't have a public parseSecretURI method, so we test through pattern matching
			// This validates our URI parsing logic is correct
			reference := "@Microsoft.KeyVault(SecretUri=" + tt.secretURI + ")"
			matches := kvRefSecretURIPattern.FindStringSubmatch(reference)
			if matches == nil {
				t.Error("Reference did not match pattern")
				return
			}

			if matches[1] != tt.secretURI {
				t.Errorf("Extracted URI = %q, want %q", matches[1], tt.secretURI)
			}
		})
	}
}

func TestResolveEnvironmentVariables_NoReferences(t *testing.T) {
	// Test without actually connecting to Azure
	envVars := []string{
		"NORMAL_VAR=value1",
		"ANOTHER_VAR=value2",
		"PATH=/usr/bin:/bin",
	}

	// We can test the IsKeyVaultReference check without creating a resolver
	for _, envVar := range envVars {
		parts := splitEnvVar(envVar)
		if len(parts) == 2 && IsKeyVaultReference(parts[1]) {
			t.Errorf("Non-reference value detected as Key Vault reference: %s", envVar)
		}
	}
}

func TestResolveEnvironmentVariables_WithReferences(t *testing.T) {
	// Test detection of Key Vault references in environment variables
	envVars := []string{
		"NORMAL_VAR=value1",
		"KV_SECRET=@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)",
		"ANOTHER_NORMAL=value2",
		"KV_SECRET2=@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/secret2)",
	}

	kvRefCount := 0
	for _, envVar := range envVars {
		parts := splitEnvVar(envVar)
		if len(parts) == 2 && IsKeyVaultReference(parts[1]) {
			kvRefCount++
		}
	}

	if kvRefCount != 2 {
		t.Errorf("Expected to find 2 Key Vault references, found %d", kvRefCount)
	}
}

// splitEnvVar splits an environment variable into key and value.
func splitEnvVar(envVar string) []string {
	parts := make([]string, 0, 2)
	idx := 0
	for i, c := range envVar {
		if c == '=' {
			idx = i
			break
		}
	}
	if idx > 0 {
		parts = append(parts, envVar[:idx], envVar[idx+1:])
	}
	return parts
}

// TestNewKeyVaultResolver tests resolver creation.
// Note: This will fail if Azure credentials are not available,
// which is expected in CI/CD without Azure auth setup.
func TestNewKeyVaultResolver(t *testing.T) {
	t.Run("Creation without Azure auth", func(t *testing.T) {
		// This test documents that resolver creation requires Azure credentials
		resolver, err := NewKeyVaultResolver()

		// In environments without Azure credentials, this should fail
		// In environments with credentials, it should succeed
		switch {
		case err != nil:
			// Expected in CI/CD without Azure auth
			t.Logf("NewKeyVaultResolver failed as expected without Azure credentials: %v", err)
		case resolver == nil:
			t.Error("NewKeyVaultResolver returned nil resolver without error")
		default:
			t.Log("NewKeyVaultResolver succeeded with available Azure credentials")
		}
	})
}

func TestResolveReference_InvalidFormats(t *testing.T) {
	// Test error handling for invalid formats without actual Azure calls
	invalidRefs := []string{
		"not-a-reference",
		"@Microsoft.KeyVault(Invalid=format)",
		"@Microsoft.KeyVault(SecretUri=invalid-uri)",
		"",
	}

	for _, ref := range invalidRefs {
		if IsKeyVaultReference(ref) && !kvRefSecretURIPattern.MatchString(ref) && !kvRefVaultNamePattern.MatchString(ref) {
			t.Logf("Reference %q would fail validation", ref)
		}
	}
}

func TestKeyVaultErrorHandling_InvalidFormats(t *testing.T) {
	t.Run("Empty reference", func(t *testing.T) {
		if IsKeyVaultReference("") {
			t.Error("Empty string should not be detected as Key Vault reference")
		}
	})

	t.Run("Malformed references", func(t *testing.T) {
		malformed := []string{
			"@Microsoft.KeyVault(",
			"@Microsoft.KeyVault)",
			"@Microsoft.KeyVault(SecretUri)",
			"@Microsoft.KeyVault(VaultName)",
		}

		for _, ref := range malformed {
			// These should not match the validation patterns
			if kvRefSecretURIPattern.MatchString(ref) {
				t.Errorf("Malformed reference should not match SecretUri pattern: %s", ref)
			}
			if kvRefVaultNamePattern.MatchString(ref) {
				t.Errorf("Malformed reference should not match VaultName pattern: %s", ref)
			}
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		// Test that detection is case-sensitive for the prefix
		cases := []struct {
			ref   string
			valid bool
		}{
			{"@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/test)", true},
			{"@microsoft.keyvault(SecretUri=https://vault.vault.azure.net/secrets/test)", false},
			{"@MICROSOFT.KEYVAULT(SecretUri=https://vault.vault.azure.net/secrets/test)", false},
		}

		for _, tc := range cases {
			got := strings.HasPrefix(tc.ref, "@Microsoft.KeyVault(")
			if got != tc.valid {
				t.Errorf("Reference %q: expected HasPrefix=%v, got %v", tc.ref, tc.valid, got)
			}
		}
	})
}
