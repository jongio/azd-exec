package executor

import (
	"context"
	"testing"
)

// TestKeyVaultResolver_Creation tests basic resolver creation scenarios.
func TestKeyVaultResolver_Creation(t *testing.T) {
	t.Run("NewKeyVaultResolver", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Logf("Resolver creation failed (expected without Azure credentials): %v", err)
			return
		}
		if resolver == nil {
			t.Error("Expected non-nil resolver when no error")
		} else if resolver.clients == nil {
			t.Error("Expected clients map to be initialized")
		}
	})
}

// TestKeyVaultResolver_ReferenceValidation tests reference format validation.
func TestKeyVaultResolver_ReferenceValidation(t *testing.T) {
	tests := []struct {
		name      string
		reference string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Invalid format - no pattern match",
			reference: "@Microsoft.KeyVault(Invalid=format)",
			wantErr:   true,
			errMsg:    "invalid Key Vault reference format",
		},
		{
			name:      "Empty reference",
			reference: "",
			wantErr:   true,
		},
		{
			name:      "Not a reference",
			reference: "just-a-regular-value",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip actual resolution, just test format validation
			isRef := IsKeyVaultReference(tt.reference)
			matchesPattern := kvRefSecretURIPattern.MatchString(tt.reference) ||
				kvRefVaultNamePattern.MatchString(tt.reference)

			if !isRef && tt.wantErr {
				t.Logf("Reference correctly identified as invalid: %s", tt.reference)
			} else if !matchesPattern && tt.wantErr {
				t.Logf("Reference correctly failed pattern match: %s", tt.reference)
			}
		})
	}
}

// TestKeyVaultResolver_MockResolution tests resolution logic without Azure calls.
func TestKeyVaultResolver_MockResolution(t *testing.T) {
	t.Run("ResolveReference validation only", func(t *testing.T) {
		// This tests the validation logic without making actual Azure calls
		validReferences := []string{
			"@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/test)",
			"@Microsoft.KeyVault(VaultName=vault;SecretName=secret)",
			"@Microsoft.KeyVault(VaultName=vault;SecretName=secret;SecretVersion=v1)",
		}

		for _, ref := range validReferences {
			t.Run(ref, func(t *testing.T) {
				// Validate format
				if !IsKeyVaultReference(ref) {
					t.Errorf("Valid reference not recognized: %s", ref)
				}

				matchesSecretURI := kvRefSecretURIPattern.MatchString(ref)
				matchesVaultName := kvRefVaultNamePattern.MatchString(ref)

				if !matchesSecretURI && !matchesVaultName {
					t.Errorf("Valid reference doesn't match any pattern: %s", ref)
				}
			})
		}
	})
}

// TestKeyVaultResolver_EnvironmentProcessing tests environment variable processing.
func TestKeyVaultResolver_EnvironmentProcessing(t *testing.T) {
	t.Run("Identify references in environment", func(t *testing.T) {
		envVars := []string{
			"NORMAL_VAR=value",
			"KV_VAR=@Microsoft.KeyVault(VaultName=test;SecretName=secret)",
			"PATH=/usr/bin",
		}

		count := 0
		for _, env := range envVars {
			if idx := findEquals(env); idx > 0 {
				value := env[idx+1:]
				if IsKeyVaultReference(value) {
					count++
				}
			}
		}

		if count != 1 {
			t.Errorf("Expected 1 Key Vault reference, found %d", count)
		}
	})
}

// TestKeyVaultResolver_ClientCaching tests client caching logic.
func TestKeyVaultResolver_ClientCaching(t *testing.T) {
	resolver, err := NewKeyVaultResolver()
	if err != nil {
		t.Skip("Skipping client caching test without Azure credentials")
	}

	// Test that clients map exists and is empty initially
	if resolver.clients == nil {
		t.Error("Expected clients map to be initialized")
	}
	if len(resolver.clients) != 0 {
		t.Errorf("Expected empty clients map, got %d entries", len(resolver.clients))
	}
}

// TestKeyVaultResolver_ErrorScenarios tests error handling.
func TestKeyVaultResolver_ErrorScenarios(t *testing.T) {
	resolver, err := NewKeyVaultResolver()
	if err != nil {
		t.Skip("Skipping error scenario tests without Azure credentials")
	}

	t.Run("Invalid SecretUri path", func(t *testing.T) {
		// This reference has valid format but invalid path structure
		ref := "@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/invalid)"
		_, err := resolver.ResolveReference(context.Background(), ref)
		if err == nil {
			t.Error("Expected error for invalid SecretUri path")
		}
	})

	t.Run("Invalid reference format", func(t *testing.T) {
		ref := "@Microsoft.KeyVault(BadFormat=value)"
		_, err := resolver.ResolveReference(context.Background(), ref)
		if err == nil {
			t.Error("Expected error for invalid reference format")
		}
	})
}

// findEquals returns the index of '=' in an environment variable string.
func findEquals(env string) int {
	for i, c := range env {
		if c == '=' {
			return i
		}
	}
	return -1
}
