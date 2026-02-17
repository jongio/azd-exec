---
name: azd-exec
description: |
  Execute scripts and commands with Azure Developer CLI context, environment variables,
  and automatic Azure Key Vault secret resolution. USE FOR: run scripts with azure context,
  execute commands with env vars, resolve key vault secrets, cross-platform scripting,
  interactive script execution, shell selection. DO NOT USE FOR: service orchestration
  (use azd-app), Azure deployments (use azd deploy), long-running services.
---

# azd-exec Extension

azd-exec is an Azure Developer CLI extension that executes commands and scripts
with full access to azd environment variables, Azure credentials, and automatic
Azure Key Vault secret resolution.

## When to Use

- Running scripts that need azd environment variables (AZURE_ENV_NAME, AZURE_SUBSCRIPTION_ID, etc.)
- Executing commands that reference Azure Key Vault secrets in environment variables
- Cross-platform scripting with automatic shell detection
- Interactive script execution requiring stdin passthrough

## Command Syntax

```
azd exec [flags] <script-or-command> [-- script-args...]
```

The first positional argument is either a **file path** (executed as a script) or an
**inline command string** (executed in the detected/specified shell). Everything after
the script argument (or after `--`) is passed as arguments to the script.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--shell` | `-s` | auto | Shell to use: `bash`, `sh`, `zsh`, `pwsh`, `powershell`, `cmd` |
| `--interactive` | `-i` | false | Connect stdin to the script for interactive input |
| `--stop-on-keyvault-error` | | false | Fail-fast when any Key Vault reference fails to resolve |
| `--cwd` | `-C` | . | Set working directory before execution |
| `--environment` | `-e` | | Load a specific azd environment by name |
| `--debug` | | false | Enable debug output to stderr |
| `--output` | `-o` | default | Output format: `default` or `json` |

## Shell Support

### Supported Shells

`bash`, `sh`, `zsh`, `pwsh` (PowerShell Core 6+), `powershell` (Windows PowerShell 5.1), `cmd`

### Auto-Detection Logic

When `--shell` is not specified:

1. **File scripts**: shell is detected from the file extension (`.sh` → bash, `.ps1` → pwsh,
   `.cmd`/`.bat` → cmd) or shebang line (`#!/bin/bash`, etc.)
2. **Inline commands**: defaults to `bash` on Linux/macOS, `powershell` on Windows

## Key Vault Secret Resolution

azd-exec automatically resolves Azure Key Vault references found in environment variables
before executing the script. Three reference formats are supported:

### Format 1: SecretUri

```
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret)
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret/version-id)
```

### Format 2: VaultName + SecretName

```
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret)
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret;SecretVersion=version-id)
```

### Format 3: akvs:// URI

```
akvs://subscription-id/myvault/my-secret
akvs://subscription-id/myvault/my-secret/version-id
```

### Resolution Behavior

- By default, unresolvable references log a warning and the original value is kept.
- With `--stop-on-keyvault-error`, any resolution failure aborts execution immediately.
- Resolution uses the current Azure credential (e.g., `az login` session or managed identity).

## Argument Passing

Everything after the script path is forwarded as script arguments. The `--` separator
is optional but useful to prevent flag conflicts:

```bash
azd exec ./build.sh --verbose              # --verbose passed to build.sh
azd exec ./build.sh -- --verbose            # same, explicit separator
azd exec ./deploy.sh -- --env staging       # avoids conflict with azd's -e flag
```

## Interactive Mode

Use `--interactive` / `-i` to connect stdin to the script process. This is required
for scripts that prompt for user input:

```bash
azd exec -i ./setup-wizard.sh
```

Without `-i`, stdin is not connected and the script cannot read interactive input.

## Error Handling

azd-exec uses typed errors:

| Error | Cause |
|-------|-------|
| `ValidationError` | Invalid input (empty path, path traversal, empty inline script) |
| `ScriptNotFoundError` | Script file does not exist |
| `InvalidShellError` | Unknown shell name passed to `--shell` |
| `ExecutionError` | Script exited with a non-zero exit code |

The process exit code from the script is propagated to the caller.

## Examples

```bash
# Execute a script file with azd environment loaded
azd exec ./setup.sh

# Run inline bash command with environment variables
azd exec 'echo $AZURE_ENV_NAME'

# Use PowerShell explicitly
azd exec --shell pwsh "Write-Host 'Hello from PowerShell'"

# Execute a PowerShell script file
azd exec --shell pwsh ./deploy.ps1

# Pass arguments to a script
azd exec ./build.sh -- --target release --verbose

# Interactive mode for scripts that prompt for input
azd exec -i ./interactive-setup.sh

# Use a specific azd environment
azd exec -e staging ./migrate-db.sh

# Fail-fast on Key Vault errors
azd exec --stop-on-keyvault-error ./prod-deploy.sh

# Change working directory before execution
azd exec -C ./scripts ./run.sh

# Debug mode to see execution details
azd exec --debug ./troubleshoot.sh
```
