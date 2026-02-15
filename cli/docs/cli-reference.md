---
title: CLI Reference
description: Complete command reference for azd exec extension
lastUpdated: 2026-01-09
tags: [cli, reference, documentation, commands]
---

# CLI Reference

Complete reference for the `azd exec` extension commands and flags.

## Overview

The `azd exec` extension allows you to execute scripts and commands with full access to your Azure Developer CLI environment variables and Azure credentials.

## Installation

```bash
# Add the extension registry
azd extension source add -n azd-exec -t url -l https://raw.githubusercontent.com/jongio/azd-exec/main/registry.json

# Install the extension
azd extension install jongio.azd.exec

# Verify installation
azd exec version
```

## Commands Overview

| Command | Description |
|---------|-------------|
| `exec` | Execute a script file or inline command with Azure context |
| `version` | Display the extension version |

---

## `azd exec`

Execute a script file or inline command with full access to azd environment variables and Azure credentials.

### Usage

```bash
azd exec [flags-before-script] <script-file-or-command> [script-args...]
```

Place azd exec flags (like `--cwd`, `--debug`, `--environment`) before the script. Everything after the script is passed through to the script; `--` remains optional if you prefer an explicit separator.

### Description

Executes scripts and commands with access to:
- All azd environment variables (`AZURE_ENV_NAME`, `AZURE_SUBSCRIPTION_ID`, etc.)
- Azure authentication context (same credentials as azd)
- Azure Key Vault integration (automatic secret resolution)
- Custom environment variables from your azd environment

The command automatically detects the script type based on file extension or shebang, and executes it in the appropriate shell.

### Examples

```bash
# Execute a script file
azd exec ./deploy.sh

# Execute an inline command
azd exec 'echo "Environment: $AZURE_ENV_NAME"'

# Specify shell explicitly
azd exec --shell pwsh ./setup.ps1

# Run in interactive mode (for scripts with prompts)
azd exec --interactive ./interactive-setup.sh

# Pass arguments to the script
azd exec ./build.sh --verbose --config release

# Inline PowerShell command
azd exec --shell pwsh 'Write-Host "Hello from $env:AZURE_ENV_NAME"'

# Execute with debug logging (flags before the script)
azd exec --debug ./deploy.sh

# Use specific environment (flags before the script)
azd exec --environment production ./deploy.sh
```

### Flags

#### Execution Control

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--shell` | `-s` | string | (auto-detect) | Shell to use for execution. Options: `bash`, `sh`, `zsh`, `pwsh`, `powershell`, `cmd`. Auto-detected from file extension or shebang if not specified. |
| `--interactive` | `-i` | bool | false | Run script in interactive mode, enabling user input and prompts. |
| `--stop-on-keyvault-error` |  | bool | false | Fail-fast: stop execution when any Key Vault reference fails to resolve. |

#### Global Flags (inherited from azd)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | default | Output format: `default` or `json` |
| `--debug` | | bool | false | Enable debug mode with verbose logging |
| `--no-prompt` | | bool | false | Disable prompts and use default values |
| `--cwd` | `-C` | string | (current) | Sets the current working directory before execution |
| `--environment` | `-e` | string | (default env) | The name of the azd environment to use |

### Script Arguments

Arguments after the script are passed directly to it—no `--` separator required. You can still use `--` if you prefer to explicitly stop flag parsing.

```bash
# The script receives: --verbose --config release
azd exec ./build.sh --verbose --config release

# Inline script with arguments
azd exec './process.sh "$@"' file1.txt file2.txt
```

### File vs Inline Execution

**File Execution**
```bash
# Script must exist on filesystem
azd exec ./deploy.sh

# Absolute path
azd exec /home/user/scripts/setup.sh

# Relative path
azd exec ../scripts/build.sh
```

**Inline Execution**
```bash
# Single-line command
azd exec 'echo $AZURE_ENV_NAME'

# Multi-line with semicolons (bash)
azd exec 'echo "Starting"; echo $AZURE_ENV_NAME; echo "Done"'

# PowerShell inline
azd exec --shell pwsh 'Write-Host "Hello"; Get-Date'

# Complex inline with arguments
azd exec 'echo "Args: $@"' arg1 arg2
```

### Shell Detection

When `--shell` is not specified, the shell is detected automatically:

**By File Extension:**
- `.sh` → bash
- `.bash` → bash
- `.zsh` → zsh
- `.ps1` → pwsh
- `.cmd`, `.bat` → cmd

**By Shebang (first line of file):**
```bash
#!/bin/bash       → bash
#!/bin/sh         → sh
#!/usr/bin/env zsh → zsh
#!/usr/bin/env pwsh → pwsh
```

**For Inline Commands:**
- Defaults to the system shell (`bash` on Unix, `cmd` on Windows)
- Override with `--shell` flag

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Script execution failed or script not found |
| 2 | Invalid arguments or configuration |
| Non-zero | Exit code from the executed script |

### Security Considerations

**⚠️ IMPORTANT**: `azd exec` runs scripts with full access to:
- Azure authentication credentials
- All environment variables (including secrets)
- Filesystem access
- Network access

**Best Practices:**
- ✅ Review scripts before execution: `cat ./script.sh`
- ✅ Only run trusted scripts
- ✅ Use Azure Key Vault for secrets (not environment variables)
- ✅ Use file-based scripts for complex operations
- ❌ Never pipe untrusted scripts: ~~`curl https://site.com/script.sh | azd exec -`~~
- ❌ Don't run scripts from unknown sources

### Environment Variables

All azd environment variables are available to your scripts:

**Standard Azure Variables:**
- `AZURE_ENV_NAME` - Current azd environment name
- `AZURE_SUBSCRIPTION_ID` - Azure subscription ID
- `AZURE_LOCATION` - Azure region/location
- `AZURE_RESOURCE_GROUP` - Resource group name
- `AZURE_TENANT_ID` - Azure tenant ID
- `AZURE_PRINCIPAL_ID` - Service principal ID (when using service principal auth)

**Custom Variables:**
- All variables set with `azd env set KEY VALUE`
- Variables from `.azure/<env>/.env` file

**azd-exec Specific:**
- `AZD_DEBUG` - Set to "true" when `--debug` flag is used
- `AZD_NO_PROMPT` - Set to "true" when `--no-prompt` flag is used

### Azure Key Vault Integration

`azd exec` automatically resolves Azure Key Vault references in environment variables before running your script.

> **Implementation**: Key Vault resolution is powered by the [azd-core](https://github.com/jongio/azd-core) library, ensuring consistent behavior across Azure Developer CLI tools.

**Supported Formats:**

```bash
# Format 1: SecretUri
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret)
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret/abc123)

# Format 2: VaultName and SecretName
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret)
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret;SecretVersion=abc123)

# Format 3: azd akvs URI (used internally by azd)
akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/myvault/my-secret
akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/myvault/my-secret/abc123
```

**Example Workflow:**

```bash
# 1. Store secret in Key Vault
az keyvault secret set --vault-name myvault --name db-password --value "Secret123!"

# 2. Set environment variable with Key Vault reference
azd env set-secret DB_PASSWORD

# 3. Use in script (automatically resolved)
azd exec 'echo "Password: $DB_PASSWORD"'
```

**Authentication:**
Uses the same Azure credentials as `azd`:
- Azure CLI (`az login`)
- Managed Identity (on Azure)
- Service Principal (environment variables)
- Visual Studio / VS Code authentication

**Error Handling:**
If resolution fails (secret not found, no access, etc.):
- Warning displayed to stderr (secret values are never printed)
- `azd exec` continues resolving other Key Vault references
- Successfully resolved secrets are substituted with their values
- Failed references remain unchanged (still `akvs://...` or `@Microsoft.KeyVault(...)`)

To fail-fast (abort on the first Key Vault resolution error), use `--stop-on-keyvault-error`.

---

## `azd exec version`

Display the extension version information.

### Usage

```bash
azd exec version [flags]
```

### Examples

```bash
# Show version with label
azd exec version

# Show only version number
azd exec version --quiet

# JSON output
azd exec version --output json
```

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--quiet` | `-q` | bool | false | Display only the version number |
| `--output` | `-o` | string | default | Output format: `default` or `json` |

### Output Examples

**Default:**
```
azd exec version 0.1.0
```

**Quiet:**
```
0.1.0
```

**JSON:**
```json
{
  "version": "0.1.0"
}
```

---

## Global Flags

These flags are available for all commands and match azd's global flags for compatibility:

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format (default, json) |
| `--debug` | | Enable debug logging |
| `--no-prompt` | | Disable prompts |
| `--cwd` | `-C` | Sets the current working directory |
| `--environment` | `-e` | The name of the environment to use |

---

## Legacy Compatibility

For backwards compatibility, the legacy `azd script` invocation is supported:

```bash
# Legacy (still works)
azd script ./deploy.sh

# Modern (preferred)
azd exec ./deploy.sh
```

Both invoke the same functionality. New projects should use `azd exec`.

---

## Common Workflows

### First Time Setup

```bash
# 1. Add registry
azd extension source add -n azd-exec -t url -l https://raw.githubusercontent.com/jongio/azd-exec/main/registry.json

# 2. Install extension
azd extension install jongio.azd.exec

# 3. Verify
azd exec version
```

### Daily Development

```bash
# Review script before running (security best practice)
cat ./deploy.sh

# Execute deployment script
azd exec ./deploy.sh

# Quick environment check
azd exec 'echo "Deploying to $AZURE_ENV_NAME in $AZURE_LOCATION"'

# Run tests with arguments
azd exec ./test.sh --verbose --coverage
```

### Multi-Environment Deployment

```bash
# Deploy to development
azd exec --environment dev ./deploy.sh

# Deploy to production (with confirmation)
azd exec --environment prod --interactive ./deploy.sh
```

### Debugging Scripts

```bash
# Enable debug logging
azd exec --debug ./deploy.sh

# Check what environment variables are available
azd exec 'env | grep AZURE_' --debug

# PowerShell debugging
azd exec --shell pwsh --debug 'Get-ChildItem Env: | Where-Object Name -like "AZURE_*"'
```

### Working with Key Vault Secrets

```bash
# Set Key Vault reference
azd env set-secret API_KEY

# Script automatically gets the resolved secret value
azd exec ./deploy.sh

# Verify resolution (debug mode shows resolution process)
azd exec 'echo $API_KEY' --debug
```

---

## Troubleshooting

### Script Not Found

```bash
# Error: script file not found: ./deploy.sh

# Solution: Verify path
ls -la ./deploy.sh

# Use absolute path
azd exec /full/path/to/deploy.sh
```

### Permission Denied

```bash
# Error: permission denied

# Solution: Make script executable
chmod +x ./deploy.sh
azd exec ./deploy.sh
```

### Shell Not Detected

```bash
# Solution: Specify shell explicitly
azd exec --shell bash ./script

# Or add shebang to script
echo '#!/bin/bash' | cat - script.sh > temp && mv temp script.sh
```

### Key Vault Access Denied

```bash
# Error: Failed to resolve Key Vault reference

# Solution: Grant access to Key Vault
az keyvault set-policy --name myvault --upn user@example.com --secret-permissions get list

# Or use Managed Identity / Service Principal with proper permissions
```

### Environment Variables Not Available

```bash
# Verify azd environment
azd env list
azd env select <env-name>

# Check environment variables
azd env get-values

# Set missing variables
azd env set MY_VAR "my-value"
```

---

## Related Documentation

- [Security Review](./security-review.md) - Detailed security analysis and best practices
- [Threat Model](./threat-model.md) - Security threat analysis and mitigations
- [Main README](../../README.md) - Getting started guide and examples
- [Azure Developer CLI Docs](https://learn.microsoft.com/azure/developer/azure-developer-cli/) - azd documentation

---

## Contributing

For development and contribution guidelines, see [CONTRIBUTING.md](../../CONTRIBUTING.md).

---

## Support

- [Report Issues](https://github.com/jongio/azd-exec/issues)
- [View Source](https://github.com/jongio/azd-exec)
- [Release Notes](https://github.com/jongio/azd-exec/releases)
