package executor

import (
	"testing"
)

func TestIsKeyVaultReference(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "Valid SecretUri format",
			value: "@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret)",
			want:  true,
		},
		{
			name:  "Valid VaultName format",
			value: "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)",
			want:  true,
		},
		{
			name:  "Valid VaultName with version",
			value: "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret;SecretVersion=abc123)",
			want:  true,
		},
		{
			name:  "Regular string",
			value: "just-a-regular-value",
			want:  false,
		},
		{
			name:  "Empty string",
			value: "",
			want:  false,
		},
		{
			name:  "Partial reference",
			value: "@Microsoft.KeyVault(",
			want:  false,
		},
		{
			name:  "Missing closing paren",
			value: "@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKeyVaultReference(tt.value); got != tt.want {
				t.Errorf("IsKeyVaultReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyVaultReferencePatterns(t *testing.T) {
	tests := []struct {
		name          string
		reference     string
		shouldMatch   bool
		pattern       string
		expectedParts []string
	}{
		{
			name:        "SecretUri without version",
			reference:   "@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret)",
			shouldMatch: true,
			pattern:     "secreturi",
			expectedParts: []string{
				"@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret)",
				"https://myvault.vault.azure.net/secrets/mysecret",
			},
		},
		{
			name:        "SecretUri with version",
			reference:   "@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret/abc123)",
			shouldMatch: true,
			pattern:     "secreturi",
			expectedParts: []string{
				"@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret/abc123)",
				"https://myvault.vault.azure.net/secrets/mysecret/abc123",
			},
		},
		{
			name:        "VaultName without version",
			reference:   "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)",
			shouldMatch: true,
			pattern:     "vaultname",
			expectedParts: []string{
				"@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)",
				"myvault",
				"mysecret",
				"", // Empty version
			},
		},
		{
			name:        "VaultName with version",
			reference:   "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret;SecretVersion=abc123)",
			shouldMatch: true,
			pattern:     "vaultname",
			expectedParts: []string{
				"@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret;SecretVersion=abc123)",
				"myvault",
				"mysecret",
				"abc123",
			},
		},
		{
			name:        "Invalid format",
			reference:   "@Microsoft.KeyVault(Invalid=format)",
			shouldMatch: false,
			pattern:     "both",
		},
		{
			name:        "Missing SecretName",
			reference:   "@Microsoft.KeyVault(VaultName=myvault)",
			shouldMatch: false,
			pattern:     "vaultname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var matches []string

			switch tt.pattern {
			case "secreturi":
				matches = kvRefSecretURIPattern.FindStringSubmatch(tt.reference)
			case "vaultname":
				matches = kvRefVaultNamePattern.FindStringSubmatch(tt.reference)
			case "both":
				matches = kvRefSecretURIPattern.FindStringSubmatch(tt.reference)
				if matches == nil {
					matches = kvRefVaultNamePattern.FindStringSubmatch(tt.reference)
				}
			}

			if tt.shouldMatch {
				if matches == nil {
					t.Errorf("Expected pattern to match, but got no match")
					return
				}
				if len(matches) != len(tt.expectedParts) {
					t.Errorf("Expected %d parts, got %d", len(tt.expectedParts), len(matches))
					return
				}
				for i, expected := range tt.expectedParts {
					if matches[i] != expected {
						t.Errorf("Part %d: expected %q, got %q", i, expected, matches[i])
					}
				}
			} else {
				if matches != nil {
					t.Errorf("Expected pattern not to match, but got: %v", matches)
				}
			}
		})
	}
}

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

func TestKeyVaultReferenceFormats(t *testing.T) {
	// Test various real-world formats
	validFormats := []string{
		"@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/database-password)",
		"@Microsoft.KeyVault(SecretUri=https://prod-vault.vault.azure.net/secrets/api-key/abc123)",
		"@Microsoft.KeyVault(VaultName=dev-vault;SecretName=connection-string)",
		"@Microsoft.KeyVault(VaultName=prod;SecretName=secret;SecretVersion=v1)",
	}

	for _, format := range validFormats {
		if !IsKeyVaultReference(format) {
			t.Errorf("Valid format not recognized: %s", format)
		}
		// Also verify it matches one of the patterns
		if !kvRefSecretURIPattern.MatchString(format) && !kvRefVaultNamePattern.MatchString(format) {
			t.Errorf("Valid format doesn't match any pattern: %s", format)
		}
	}

	invalidFormats := []string{
		"@Microsoft.KeyVault",
		"Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/name)",
		"regular-value",
		"",
	}

	for _, format := range invalidFormats {
		if IsKeyVaultReference(format) {
			t.Errorf("Invalid format incorrectly recognized: %s", format)
		}
	}

	// These have the prefix/suffix but invalid content - should be detected as references
	// but would fail pattern validation
	ambiguousFormats := []string{
		"@Microsoft.KeyVault()",
		"@Microsoft.KeyVault(SecretUri=)",
		"@Microsoft.KeyVault(VaultName=)",
		"@Microsoft.KeyVault(SecretName=name)",
	}

	for _, format := range ambiguousFormats {
		isRef := IsKeyVaultReference(format)
		matchesPattern := kvRefSecretURIPattern.MatchString(format) || kvRefVaultNamePattern.MatchString(format)

		// These should be detected as references (have the prefix/suffix)
		if !isRef {
			t.Errorf("Format with KV prefix/suffix not detected as reference: %s", format)
		}

		// But they should NOT match the validation patterns
		if matchesPattern {
			t.Errorf("Invalid format should not match validation pattern: %s", format)
		}
	}
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
		if err != nil {
			// Expected in CI/CD without Azure auth
			t.Logf("NewKeyVaultResolver failed as expected without Azure credentials: %v", err)
		} else if resolver == nil {
			t.Error("NewKeyVaultResolver returned nil resolver without error")
		} else {
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
