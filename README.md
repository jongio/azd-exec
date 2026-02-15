---
title: azd exec
description: Execute any script with full access to your Azure Developer CLI environment variables and Azure credentials
lastUpdated: 2026-01-09
tags: [azure, cli, devops, scripts, keyvault]
---

<div align="center">

# azd exec

### **Execute Scripts with azd Environment Context**

Execute any script with full access to your Azure Developer CLI environment variables and Azure credentials.

[![CI](https://github.com/jongio/azd-exec/actions/workflows/ci.yml/badge.svg)](https://github.com/jongio/azd-exec/actions/workflows/ci.yml)
[![CodeQL](https://github.com/jongio/azd-exec/actions/workflows/codeql.yml/badge.svg)](https://github.com/jongio/azd-exec/actions/workflows/codeql.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<br />

</div>

---

## ⚡ One-Command Execute

Run any script with your full Azure context—no manual environment setup.

```bash
azd exec ./deploy.sh
```

That's it. Your script has access to all azd environment variables, Azure credentials, and configuration.

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🔧 Multiple Shell Support
Automatically detects and runs bash, sh, zsh, PowerShell, pwsh, and cmd scripts based on file extension or shebang.

### 🎯 Script Arguments
Pass arguments to your scripts seamlessly with the `--` separator for clean parameter handling.

### 🌍 Full Azure Context
Access all azd environment variables including subscription, tenant, location, and custom variables.

</td>
<td width="50%">

### 📂 Working Directory Control
Execute scripts from any directory with the `--cwd` flag for flexible automation.

### 🔄 Interactive Mode
Run scripts with interactive input support for prompts and user interaction.

### 🔐 Azure Key Vault Integration
Automatically resolves Key Vault references in environment variables, securely fetching secrets at runtime.

### ✅ Battle-Tested
Comprehensive security scanning with CodeQL and gosec (0 vulnerabilities). 86%+ test coverage.

</td>
</tr>
</table>

---

## ⚠️ Security Notice

**IMPORTANT**: `azd exec` executes scripts with **full access** to your Azure credentials and environment. Follow these security best practices:

**✅ Safe Practices**
- Only run scripts you trust and have reviewed
- Review script contents before execution (`cat ./script.sh`)
- Use inline scripts only for simple, trusted operations
- Use HTTPS for downloads, never HTTP
- Verify script sources from official Azure documentation or trusted repositories

**❌ Dangerous Practices**
- Never pipe untrusted scripts: ~~`curl https://random-site.com/script.sh | azd exec -`~~
- Don't run scripts from unknown sources
- Avoid storing secrets in environment variables (use Azure Key Vault instead)
- Don't blindly copy/paste inline commands without understanding them

**What Scripts Can Access:**
- 🔑 Azure authentication context (subscription, tenant, credentials)
- 🌍 All environment variables (including secrets)
- 📂 Full filesystem access
- 🌐 Network access

**Inline vs File-based Scripts:**
- **File scripts**: Can be reviewed before execution with `cat` or editor
- **Inline scripts**: Execute immediately—ensure you understand the command first
- **Best practice**: Use file scripts for complex operations, inline for simple queries

For detailed security information, see [Security Documentation](cli/docs/security-review.md) and [Threat Model](cli/docs/threat-model.md).

---

## 🎯 Quick Start

### 1. Install Azure Developer CLI

<details>
<summary><b>Windows</b></summary>

```powershell
winget install microsoft.azd
```
</details>

<details>
<summary><b>macOS</b></summary>

```bash
brew tap azure/azd && brew install azd
```
</details>

<details>
<summary><b>Linux</b></summary>

```bash
curl -fsSL https://aka.ms/install-azd.sh | bash
```
</details>

### 2. Install azd-exec

```bash
# Add the extension registry
azd extension source add -n jongio -t url -l https://jongio.github.io/azd-extensions/registry.json

# Install the extension
azd extension install jongio.azd.exec

# Verify installation
azd exec version
```

### 3. Run Your Script

```bash
# Review the script first
cat ./deploy.sh

# Then execute
azd exec ./deploy.sh
```

---

## 📚 Usage Examples

### Basic Execution

```bash
# Execute a script file
azd exec ./my-script.sh

# Execute an inline command
azd exec 'echo "Hello, $AZURE_ENV_NAME"'
```

For complete command reference, see [CLI Reference](cli/docs/cli-reference.md).

### Specify Shell

```bash
azd exec --shell pwsh ./deploy.ps1

# Inline with specific shell
azd exec --shell pwsh 'Write-Host $env:AZURE_ENV_NAME'
```

### Pass Arguments

```bash
azd exec ./build.sh --verbose --config release
# azd exec flags go before the script; script args go after it
# example with cwd flag: azd exec --cwd /path/to/project ./build.sh --verbose
```

### Set Working Directory

```bash
azd exec --cwd /path/to/project ./scripts/setup.sh

# Inline with working directory
azd exec --cwd /tmp 'echo $(pwd)'
```

### Interactive Mode

```bash
azd exec --interactive ./interactive-setup.sh
```

---

## 💡 Script Examples

<table>
<tr>
<td width="50%">

**Bash Script File**

```bash
#!/bin/bash
# deploy.sh

echo "Environment: $AZURE_ENV_NAME"
echo "Subscription: $AZURE_SUBSCRIPTION_ID"
echo "Location: $AZURE_LOCATION"

azd deploy --all
```

**Bash Inline**

```bash
azd exec 'echo "Deploying to $AZURE_ENV_NAME"'
azd exec 'for i in {1..3}; do echo "Step $i"; done'
```

</td>
<td width="50%">

**PowerShell Script File**

```powershell
# setup.ps1

Write-Host "Environment: $env:AZURE_ENV_NAME"
Write-Host "Resource Group: $env:AZURE_RESOURCE_GROUP"

# Your setup logic here
```

**PowerShell Inline**

```bash
azd exec --shell pwsh 'Write-Host "Hello from $env:AZURE_ENV_NAME"'
azd exec --shell pwsh 'Get-ChildItem Env: | Where-Object Name -like "AZURE_*"'
```

</td>
</tr>
</table>

**Run script files:**
```bash
# First, review the script
cat ./deploy.sh

# Then execute
azd exec ./deploy.sh
```

---

## 🌍 Environment Variables

Scripts executed by `azd exec` have access to all azd environment variables:

| Variable | Description |
|----------|-------------|
| `AZURE_ENV_NAME` | Current azd environment name |
| `AZURE_SUBSCRIPTION_ID` | Azure subscription ID |
| `AZURE_LOCATION` | Azure region/location |
| `AZURE_RESOURCE_GROUP` | Resource group name |
| `AZURE_TENANT_ID` | Azure tenant ID |
| *Custom variables* | All environment variables from your azd environment |

---

## 🔐 Azure Key Vault Integration

`azd exec` automatically resolves Azure Key Vault references in environment variables, allowing you to securely store and access secrets without hardcoding them.

> **Note**: Key Vault resolution is provided by the [azd-core](https://github.com/jongio/azd-core) library, a shared utility for Azure Developer CLI tools.

### Supported Reference Formats

**Format 1: SecretUri**
```bash
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret)
@Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/my-secret/abc123)
```

**Format 2: VaultName and SecretName**
```bash
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret)
@Microsoft.KeyVault(VaultName=myvault;SecretName=my-secret;SecretVersion=abc123)
```

**Format 3: azd akvs URI**
```bash
akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/myvault/my-secret
akvs://c3b3091e-400e-43a7-8ee5-e6e8cefdbebf/myvault/my-secret/abc123
```

Note: `azd` may export environment variables with quotes; `azd exec` trims whitespace and strips a single pair of wrapper quotes before detecting/parsing Key Vault references. The akvs:// format (Format 3) is used internally by azd and includes a subscription/tenant GUID, vault name, secret name, and optional version.

### Usage Example

**1. Store a secret in Azure Key Vault:**
```bash
az keyvault secret set --vault-name myvault --name database-password --value "SuperSecret123!"
```

**2. Set environment variable with Key Vault reference:**
```bash
azd env set-secret DATABASE_PASSWORD
```

**3. Use in your script:**
```bash
#!/bin/bash
# deploy.sh

echo "Connecting to database..."
# DATABASE_PASSWORD is automatically resolved to the actual secret value
mysql -u admin -p"$DATABASE_PASSWORD" -h myserver.mysql.database.azure.com
```

**4. Run the script:**
```bash
azd exec ./deploy.sh
```

### How It Works

1. `azd exec` scans environment variables for Key Vault references
2. Uses Azure DefaultAzureCredential (same authentication as azd)
3. Fetches secret values from Key Vault before running your script
4. Passes resolved values to your script securely
5. If resolution fails, warns but continues with original values

### Authentication

Key Vault resolution uses the same Azure credentials that `azd` uses:
- Azure CLI (`az login`)
- Managed Identity (when running on Azure)
- Environment variables (`AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`)
- Visual Studio / VS Code authentication

### Error Handling

If Key Vault resolution fails (e.g., secret not found, no access, vault doesn't exist):
- A warning is displayed to stderr (secret values are never printed)
- `azd exec` continues resolving other Key Vault references
- Successfully resolved secrets are substituted with their values
- Failed references remain unchanged (still `akvs://...` or `@Microsoft.KeyVault(...)`)

To fail-fast (abort on the first Key Vault resolution error), use:

```bash
azd exec --stop-on-keyvault-error ./script.sh
```

### Security Benefits

- ✅ **No secrets in code**: Store references, not actual secrets
- ✅ **Centralized management**: Update secrets in Key Vault, not in code
- ✅ **Access control**: Use Azure RBAC to control who can access secrets
- ✅ **Audit trail**: Key Vault logs all secret access
- ✅ **Automatic rotation**: Update secrets without changing code

---

## 🔧 Development

### Documentation

- [CLI Reference](cli/docs/cli-reference.md) - Complete command and flag reference
- [Security Review](cli/docs/security-review.md) - Security analysis and best practices
- [Threat Model](cli/docs/threat-model.md) - Security threat analysis

### Build from Source

```bash
git clone https://github.com/jongio/azd-exec.git
cd azd-exec/cli
chmod +x build.sh
./build.sh
```

Binary created in `cli/bin/exec`.

### Prerequisites

- Go 1.25.5 or later
- golangci-lint
- Node.js 20+ (for cspell)

### Commands

```bash
# Build
cd cli && ./build.sh

# Test - Run all tests (unit, integration, e2e)
pnpm test

# Test - Individual test suites
pnpm test:cli:unit          # CLI unit tests only
pnpm test:cli:integration   # CLI integration tests only  
pnpm test:web              # Web e2e tests only

# Lint
cd cli && golangci-lint run

# Spell check
npm install -g cspell
cspell "**/*.{go,md,yaml,yml}" --config cspell.json

# Security scan
go install github.com/securego/gosec/v2/cmd/gosec@latest
cd cli && gosec ./...
```

For detailed testing information, see [TESTING.md](TESTING.md).

---

## 🚀 CI/CD

GitHub Actions workflows:

- **CI**: Linting, spell checking, tests (Linux/Windows/macOS), security scanning, coverage
- **CodeQL**: Security analysis on push to main and weekly
- **Release**: Automated releases with multi-platform binaries

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

Ensure: tests pass, code linted, documentation updated, security scans pass.

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🔖 Release Notes

**Latest**: [View releases](https://github.com/jongio/azd-exec/releases)

### For Maintainers

**Automated Release (Recommended)**

1. Go to **Actions** → **Release** workflow
2. Click **Run workflow**
3. Choose bump type: **patch** (bug fixes), **minor** (features), **major** (breaking changes)
4. Click **Run workflow**

Workflow automatically: calculates version, updates files, builds binaries, creates release, updates registry.

**Manual Release (Testing)**

```bash
# Install tooling
azd extension install microsoft.azd.extensions

# Build & package
cd cli
export extension_id="jongio.azd.exec"
export extension_version="0.1.0"
azd x build --all
azd x pack

# Create release
azd x release --repo "jongio/azd-exec" --version "0.1.0" --draft
azd x publish --registry ../registry.json --version "0.1.0"
```

---

## 📎 Related Projects

- [Azure Developer CLI](https://github.com/Azure/azure-dev) - Core azd tool
- [azd-app](https://github.com/jongio/azd-app) - Run Azure apps locally

---

<div align="center">

### Need help or have questions?

[**Open an issue on GitHub →**](https://github.com/jongio/azd-exec/issues)

<br />

**Note**: Legacy invocation `azd script` is supported as an alias to `azd exec` for backwards compatibility.

<br />

Built with ❤️ for Azure developers

</div>