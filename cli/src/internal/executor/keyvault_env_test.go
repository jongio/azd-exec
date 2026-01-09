package executor

import (
	"strings"
	"testing"

	"github.com/jongio/azd-core/keyvault"
)

func TestIsKeyVaultReference_DetectsValidFormats(t *testing.T) {
	valid := []string{
		"@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret)",
		"@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)",
		"@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret;SecretVersion=abc123)",
		"akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD",
		"\"akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD\"",
		"  'akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD/abc123'  ",
	}

	for _, v := range valid {
		if !keyvault.IsKeyVaultReference(v) {
			t.Errorf("expected value to be detected as reference: %s", v)
		}
	}
}

func TestIsKeyVaultReference_RejectsInvalidFormats(t *testing.T) {
	invalid := []string{
		"",
		"just-a-regular-value",
		"@Microsoft.KeyVault(",
		"@Microsoft.KeyVault(SecretUri)",
		"@Microsoft.KeyVault(VaultName)",
		"@microsoft.keyvault(SecretUri=https://vault.vault.azure.net/secrets/test)",
	}

	for _, v := range invalid {
		if keyvault.IsKeyVaultReference(v) {
			t.Errorf("expected value NOT to be detected as reference: %s", v)
		}
	}
}

func TestResolveEnvironmentVariables_DetectionWithQuotes(t *testing.T) {
	envVars := []string{
		"NORMAL_VAR=value1",
		"AKVS_UNQUOTED=akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD",
		"AKVS_DQUOTED=\"akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD\"",
		"AKVS_SQUOTED='akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD/abc123'",
	}

	count := 0
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 && keyvault.IsKeyVaultReference(parts[1]) {
			count++
		}
	}

	if count != 3 {
		t.Fatalf("expected 3 references detected, got %d", count)
	}
}
