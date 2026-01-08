# Example script demonstrating Azure Key Vault reference resolution
# This script shows how azd exec automatically resolves Key Vault references

Write-Host "=== Azure Key Vault Integration Demo ===" -ForegroundColor Cyan
Write-Host ""

# Normal environment variables
Write-Host "Environment: $env:AZURE_ENV_NAME"
Write-Host "Location: $env:AZURE_LOCATION"
Write-Host ""

# If you set an environment variable with a Key Vault reference like:
# azd env set-secret API_KEY
#
# Then azd exec will automatically fetch the secret value before running this script
# and API_KEY will contain the actual secret value, not the reference string

if ($env:API_KEY) {
    Write-Host "API Key length: $($env:API_KEY.Length) characters" -ForegroundColor Green
    Write-Host "(The actual key is securely retrieved from Key Vault)" -ForegroundColor Green
} else {
    Write-Host "API_KEY not set" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "To demo this feature:" -ForegroundColor Yellow
    Write-Host "  1. Create a Key Vault secret:"
    Write-Host "     az keyvault secret set --vault-name myvault --name api-key --value 'sk-abc123xyz'"
    Write-Host ""
    Write-Host "  2. Set the environment variable with a Key Vault reference:"
    Write-Host "     azd env set-secret API_KEY"
    Write-Host ""
    Write-Host "  3. Run this script:"
    Write-Host "     azd exec examples/keyvault-demo.ps1"
}

Write-Host ""
Write-Host "=== Supported Key Vault Reference Formats ===" -ForegroundColor Cyan
Write-Host "Format 1: @Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret)"
Write-Host "Format 2: @Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret)"
Write-Host "Format 3: @Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret;SecretVersion=abc123)"

Write-Host ""
Write-Host "=== Security Benefits ===" -ForegroundColor Cyan
Write-Host "[OK] No secrets in code or config files"
Write-Host "[OK] Centralized secret management in Azure Key Vault"
Write-Host "[OK] Access control via Azure RBAC"
Write-Host "[OK] Audit trail of secret access"
Write-Host "[OK] Automatic secret rotation support"
