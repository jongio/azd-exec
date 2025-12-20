<div align="center">

# azd exec

### **Run Scripts with Azure Context**

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
- Use HTTPS for downloads, never HTTP
- Verify script sources from official Azure documentation or trusted repositories
- Review script contents before execution

**❌ Dangerous Practices**
- Never pipe untrusted scripts: ~~`curl https://random-site.com/script.sh | azd exec -`~~
- Don't run scripts from unknown sources
- Avoid storing secrets in environment variables (use Azure Key Vault instead)

**What Scripts Can Access:**
- 🔑 Azure authentication context (subscription, tenant, credentials)
- 🌍 All environment variables (including secrets)
- 📂 Full filesystem access
- 🌐 Network access

For detailed security information, see [Security Documentation](cli/docs/SECURITY-REVIEW.md) and [Threat Model](cli/docs/THREAT-MODEL.md).

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

### 2. Enable Extensions & Install azd-exec

```bash
# Enable azd extensions
azd config set alpha.extension.enabled on

# Add the extension registry
azd extension source add -n azd-exec -t url -l https://raw.githubusercontent.com/jongio/azd-exec/main/registry.json

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

## 📚 Usage Examples

### Basic Execution

```bash
azd exec ./my-script.sh
```

### Specify Shell

```bash
azd exec ./deploy.ps1 --shell pwsh
```

### Pass Arguments

```bash
azd exec ./build.sh -- --verbose --config release
```

### Set Working Directory

```bash
azd exec ./scripts/setup.sh --cwd /path/to/project
```

### Interactive Mode

```bash
azd exec ./interactive-setup.sh --interactive
```

---

## 💡 Script Examples

<table>
<tr>
<td width="50%">

**Bash Script**

```bash
#!/bin/bash
# deploy.sh

echo "Environment: $AZURE_ENV_NAME"
echo "Subscription: $AZURE_SUBSCRIPTION_ID"
echo "Location: $AZURE_LOCATION"

azd deploy --all
```

</td>
<td width="50%">

**PowerShell Script**

```powershell
# setup.ps1

Write-Host "Environment: $env:AZURE_ENV_NAME"
Write-Host "Resource Group: $env:AZURE_RESOURCE_GROUP"

# Your setup logic here
```

</td>
</tr>
</table>

**Run it:**
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

### Usage Example

**1. Store a secret in Azure Key Vault:**
```bash
az keyvault secret set --vault-name myvault --name database-password --value "SuperSecret123!"
```

**2. Set environment variable with Key Vault reference:**
```bash
azd env set DATABASE_PASSWORD "@Microsoft.KeyVault(VaultName=myvault;SecretName=database-password)"
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
- A warning is displayed to stderr
- Script continues with the original Key Vault reference string
- This allows scripts to handle missing secrets gracefully

### Security Benefits

- ✅ **No secrets in code**: Store references, not actual secrets
- ✅ **Centralized management**: Update secrets in Key Vault, not in code
- ✅ **Access control**: Use Azure RBAC to control who can access secrets
- ✅ **Audit trail**: Key Vault logs all secret access
- ✅ **Automatic rotation**: Update secrets without changing code

---

## 🔧 Development

### Build from Source

```bash
git clone https://github.com/jongio/azd-exec.git
cd azd-exec/cli
chmod +x build.sh
./build.sh
```

Binary created in `cli/bin/exec`.

### Prerequisites

- Go 1.23 or later
- golangci-lint
- Node.js 20+ (for cspell)

### Commands

### Commands

```bash
# Build
cd cli && ./build.sh

# Test
cd cli && go test ./...

# Lint
cd cli && golangci-lint run

# Spell check
npm install -g cspell
cspell "**/*.{go,md,yaml,yml}" --config cspell.json

# Security scan
go install github.com/securego/gosec/v2/cmd/gosec@latest
cd cli && gosec ./...
```

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