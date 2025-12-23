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
			} else if matches != nil {
				t.Errorf("Expected pattern not to match, but got: %v", matches)
			}
		})
	}
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
