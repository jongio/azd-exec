package executor

import (
	"testing"

	"github.com/jongio/azd-core/keyvault"
)

func TestIsKeyVaultReference_Smoke(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "secret uri", value: "@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret)", want: true},
		{name: "vault + name", value: "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)", want: true},
		{name: "vault + name + version", value: "@Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret;SecretVersion=abc123)", want: true},
		{name: "akvs", value: "akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD", want: true},
		{name: "akvs version", value: "akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/fookv/REDIS-CACHE-PASSWORD/abc123", want: true},
		{name: "non reference", value: "just-a-regular-value", want: false},
		{name: "empty", value: "", want: false},
		{name: "partial", value: "@Microsoft.KeyVault(", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := keyvault.IsKeyVaultReference(tc.value); got != tc.want {
				t.Fatalf("IsKeyVaultReference(%s) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}
