#!/bin/bash
# Example script demonstrating Azure Key Vault reference resolution
# This script shows how azd exec automatically resolves Key Vault references

echo "=== Azure Key Vault Integration Demo ==="
echo ""

# Normal environment variable
echo "Environment: $AZURE_ENV_NAME"
echo "Location: $AZURE_LOCATION"
echo ""

# If you set an environment variable with a Key Vault reference like:
# azd env set DATABASE_PASSWORD "@Microsoft.KeyVault(VaultName=myvault;SecretName=db-password)"
#
# Then azd exec will automatically fetch the secret value before running this script
# and DATABASE_PASSWORD will contain the actual secret value, not the reference string

if [ -n "$DATABASE_PASSWORD" ]; then
    echo "Database password length: ${#DATABASE_PASSWORD} characters"
    echo "(The actual password is securely retrieved from Key Vault)"
else
    echo "DATABASE_PASSWORD not set"
    echo "To demo this feature:"
    echo "  1. Create a Key Vault secret:"
    echo "     az keyvault secret set --vault-name myvault --name db-password --value 'SuperSecret123!'"
    echo ""
    echo "  2. Set the environment variable with a Key Vault reference:"
    echo "     azd env set DATABASE_PASSWORD '@Microsoft.KeyVault(VaultName=myvault;SecretName=db-password)'"
    echo ""
    echo "  3. Run this script:"
    echo "     azd exec run examples/keyvault-demo.sh"
fi

echo ""
echo "=== Supported Key Vault Reference Formats ==="
echo "Format 1: @Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret)"
echo "Format 2: @Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret)"
echo "Format 3: @Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret;SecretVersion=abc123)"
