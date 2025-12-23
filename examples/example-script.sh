#!/bin/bash
# Example script demonstrating azd context access

echo "==== Azure Developer CLI Context ===="
echo "Environment Name: ${AZURE_ENV_NAME:-not set}"
echo "Subscription ID: ${AZURE_SUBSCRIPTION_ID:-not set}"
echo "Location: ${AZURE_LOCATION:-not set}"
echo "Resource Group: ${AZURE_RESOURCE_GROUP:-not set}"
echo ""
echo "Current directory: $(pwd)"
echo "Script arguments: $@"
echo ""
echo "✅ Script executed successfully with azd context!"
