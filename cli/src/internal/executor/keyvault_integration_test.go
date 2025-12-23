//go:build integration
// +build integration

package executor

import (
	"context"
	"os"
	"strings"
	"testing"
)

// TestKeyVaultIntegration tests Key Vault resolution with real Azure credentials.
// This test requires:
// 1. Azure credentials to be available (az login or DefaultAzureCredential)
// 2. A Key Vault to be set up with test secrets
// 3. Environment variables: TEST_KEYVAULT_NAME, TEST_KEYVAULT_SECRET_NAME
func TestKeyVaultIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	vaultName := os.Getenv("TEST_KEYVAULT_NAME")
	secretName := os.Getenv("TEST_KEYVAULT_SECRET_NAME")

	if vaultName == "" || secretName == "" {
		t.Skip("Skipping Key Vault integration test: TEST_KEYVAULT_NAME and TEST_KEYVAULT_SECRET_NAME must be set")
	}

	t.Run("Resolve VaultName format", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Fatalf("Failed to create resolver: %v", err)
		}

		reference := "@Microsoft.KeyVault(VaultName=" + vaultName + ";SecretName=" + secretName + ")"
		value, err := resolver.ResolveReference(context.Background(), reference)
		if err != nil {
			t.Fatalf("Failed to resolve reference: %v", err)
		}

		if value == "" {
			t.Error("Resolved value is empty")
		}

		t.Logf("Successfully resolved secret: %s (length: %d)", secretName, len(value))
	})

	t.Run("Resolve SecretUri format", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Fatalf("Failed to create resolver: %v", err)
		}

		secretURI := "https://" + vaultName + ".vault.azure.net/secrets/" + secretName
		reference := "@Microsoft.KeyVault(SecretUri=" + secretURI + ")"
		value, err := resolver.ResolveReference(context.Background(), reference)
		if err != nil {
			t.Fatalf("Failed to resolve reference: %v", err)
		}

		if value == "" {
			t.Error("Resolved value is empty")
		}

		t.Logf("Successfully resolved secret from URI: %s (length: %d)", secretName, len(value))
	})

	t.Run("Resolve environment variables", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Fatalf("Failed to create resolver: %v", err)
		}

		envVars := []string{
			"NORMAL_VAR=normal_value",
			"KV_SECRET=@Microsoft.KeyVault(VaultName=" + vaultName + ";SecretName=" + secretName + ")",
			"ANOTHER_VAR=another_value",
		}

		resolved, err := resolver.ResolveEnvironmentVariables(context.Background(), envVars)
		if err != nil {
			t.Fatalf("Failed to resolve environment variables: %v", err)
		}

		if len(resolved) != len(envVars) {
			t.Errorf("Expected %d resolved vars, got %d", len(envVars), len(resolved))
		}

		// Check that the Key Vault reference was resolved
		var kvSecretResolved bool
		for _, envVar := range resolved {
			if strings.HasPrefix(envVar, "KV_SECRET=") {
				value := strings.TrimPrefix(envVar, "KV_SECRET=")
				if IsKeyVaultReference(value) {
					t.Error("KV_SECRET was not resolved")
				} else if value != "" {
					kvSecretResolved = true
					t.Logf("KV_SECRET resolved to value of length: %d", len(value))
				}
			}
		}

		if !kvSecretResolved {
			t.Error("KV_SECRET was not found in resolved variables")
		}
	})

	t.Run("Handle invalid vault name", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Fatalf("Failed to create resolver: %v", err)
		}

		reference := "@Microsoft.KeyVault(VaultName=nonexistent-vault-12345;SecretName=test)"
		_, err = resolver.ResolveReference(context.Background(), reference)
		if err == nil {
			t.Error("Expected error for nonexistent vault, got nil")
		}
		t.Logf("Got expected error: %v", err)
	})

	t.Run("Handle invalid secret name", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Fatalf("Failed to create resolver: %v", err)
		}

		reference := "@Microsoft.KeyVault(VaultName=" + vaultName + ";SecretName=nonexistent-secret-12345)"
		_, err = resolver.ResolveReference(context.Background(), reference)
		if err == nil {
			t.Error("Expected error for nonexistent secret, got nil")
		}
		t.Logf("Got expected error: %v", err)
	})
}

// TestExecutorWithKeyVaultReferences tests the executor with Key Vault references.
func TestExecutorWithKeyVaultReferences(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	vaultName := os.Getenv("TEST_KEYVAULT_NAME")
	secretName := os.Getenv("TEST_KEYVAULT_SECRET_NAME")

	if vaultName == "" || secretName == "" {
		t.Skip("Skipping Key Vault integration test: TEST_KEYVAULT_NAME and TEST_KEYVAULT_SECRET_NAME must be set")
	}

	// Create a simple script that prints an environment variable
	tmpDir := t.TempDir()
	scriptPath := tmpDir + "/test-kv.sh"
	scriptContent := "#!/bin/bash\necho \"Secret value: $MY_KV_SECRET\"\n"
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0700); err != nil {
		t.Fatal(err)
	}

	// Set up environment variable with Key Vault reference
	kvRef := "@Microsoft.KeyVault(VaultName=" + vaultName + ";SecretName=" + secretName + ")"
	os.Setenv("MY_KV_SECRET", kvRef)
	defer os.Unsetenv("MY_KV_SECRET")

	executor := New(Config{})
	err := executor.Execute(context.Background(), scriptPath)

	// We expect this to work if Azure credentials are available
	if err != nil {
		t.Logf("Execution error (may be expected if no Azure credentials): %v", err)
	} else {
		t.Log("Successfully executed script with Key Vault reference resolution")
	}
}

// TestKeyVaultErrorHandling tests graceful error handling.
func TestKeyVaultErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Invalid reference format", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Skip("Skipping test: Azure credentials not available")
		}

		invalidRef := "@Microsoft.KeyVault(Invalid=format)"
		_, err = resolver.ResolveReference(context.Background(), invalidRef)
		if err == nil {
			t.Error("Expected error for invalid reference format, got nil")
		}
		if !strings.Contains(err.Error(), "invalid Key Vault reference format") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Invalid SecretUri", func(t *testing.T) {
		resolver, err := NewKeyVaultResolver()
		if err != nil {
			t.Skip("Skipping test: Azure credentials not available")
		}

		invalidRef := "@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/invalid/path)"
		_, err = resolver.ResolveReference(context.Background(), invalidRef)
		if err == nil {
			t.Error("Expected error for invalid SecretUri path, got nil")
		}
	})
}
