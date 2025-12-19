# azd-exec

[![CI](https://github.com/jongio/azd-exec/actions/workflows/ci.yml/badge.svg)](https://github.com/jongio/azd-exec/actions/workflows/ci.yml)
[![CodeQL](https://github.com/jongio/azd-exec/actions/workflows/codeql.yml/badge.svg)](https://github.com/jongio/azd-exec/actions/workflows/codeql.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Execute scripts with Azure Developer CLI (azd) context and environment variables.

## Overview

`azd-exec` is an Azure Developer CLI extension that allows you to execute commands and scripts with full access to the azd environment, including:

- All azd environment variables
- Azure subscription and tenant information
- Project configuration and settings
- Current working environment context

This extension is perfect for automation tasks, custom deployment scripts, environment setup, and any scenario where you need to run scripts within the azd context.

## ⚠️ Security Notice

**IMPORTANT**: `azd-exec` is a powerful developer tool that executes scripts with **full access** to your Azure credentials, environment variables, and azd context. Please follow these security best practices:

### ✅ Safe Practices

- ✅ **Only run scripts you trust** - Review all scripts before execution
- ✅ **Use HTTPS for downloads** - Never use HTTP to download scripts
- ✅ **Verify script sources** - Only use scripts from official Azure documentation or trusted repositories
- ✅ **Review script contents** - Always read scripts before running them, especially from tutorials or blog posts

### ❌ Dangerous Practices

- ❌ **Never pipe untrusted scripts** - Avoid: `curl https://random-site.com/script.sh | azd exec run -`
- ❌ **Don't run scripts from unknown sources** - Verify author identity and repository ownership
- ❌ **Avoid storing secrets in environment variables** - Use Azure Key Vault or managed identities instead
- ❌ **Don't blindly follow tutorials** - Always review and understand script behavior

### What Scripts Can Access

Scripts executed by `azd-exec` inherit:
- 🔑 **Azure authentication context** (subscription, tenant, credentials)
- 🌍 **All environment variables** (including any secrets you may have set)
- 📂 **Full filesystem access** (same permissions as your user account)
- 🌐 **Network access** (can make external connections)

**For detailed security information**, see [Security Documentation](cli/docs/SECURITY-REVIEW.md) and [Threat Model](cli/docs/THREAT-MODEL.md).

## Features

- ✨ **Automatic Shell Detection**: Detects the appropriate shell based on script file extension or shebang
- 🔧 **Multiple Shell Support**: Works with bash, sh, zsh, PowerShell, pwsh, and cmd
- 🌍 **Environment Context**: Full access to azd environment variables
- 🎯 **Script Arguments**: Pass arguments to your scripts
- 📂 **Working Directory Control**: Execute scripts in any directory
- 🔄 **Interactive Mode**: Run scripts with interactive input
- 🔒 **Security**: Comprehensive security scanning with CodeQL and gosec (0 vulnerabilities)
- ✅ **Well Tested**: Extensive unit and integration tests with 86%+ coverage

## Installation

### Prerequisites

- [Azure Developer CLI (azd)](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd) installed
- Go 1.23 or later (for building from source)

### Install via azd (Recommended)

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

### Install from Release

Download the latest release for your platform from the [releases page](https://github.com/jongio/azd-exec/releases).

### Build from Source

```bash
git clone https://github.com/jongio/azd-exec.git
cd azd-exec/cli
chmod +x build.sh
./build.sh
```

The binary will be created in `cli/bin/exec`.

**Note:** This extension supports the legacy invocation `azd script` as an alias to `azd exec` for backwards compatibility.

## Usage

### Basic Usage

Execute a script file:

```bash
azd exec run ./my-script.sh
```

### Specify Shell

Explicitly specify which shell to use:

```bash
azd exec run ./deploy.ps1 --shell pwsh
```

### Pass Arguments to Script

Pass arguments to your script after `--`:

```bash
azd exec run ./build.sh -- --verbose --config release
```

### Set Working Directory

Execute script from a specific directory:

```bash
azd exec run ./scripts/setup.sh --cwd /path/to/project
```

### Interactive Mode

Run script with interactive input:

```bash
azd exec run ./interactive-setup.sh --interactive
```

### Get Version

```bash
azd exec version
```

## Script Examples

### Bash Script with azd Context

```bash
#!/bin/bash
# deploy.sh

echo "Deploying to environment: $AZURE_ENV_NAME"
echo "Subscription: $AZURE_SUBSCRIPTION_ID"
echo "Location: $AZURE_LOCATION"

# Your deployment logic here
azd deploy --all
```

Run it:
```bash
# First, review the script contents
cat ./deploy.sh

# Then execute
azd exec run ./deploy.sh
```

### PowerShell Script

```powershell
# setup.ps1

Write-Host "Setting up environment: $env:AZURE_ENV_NAME"
Write-Host "Resource Group: $env:AZURE_RESOURCE_GROUP"

# Your setup logic here
```

## Environment Variables

When you run a script using `azd-exec`, it has access to all azd environment variables, including:

- `AZURE_ENV_NAME`: Current azd environment name
- `AZURE_SUBSCRIPTION_ID`: Azure subscription ID
- `AZURE_LOCATION`: Azure region/location
- `AZURE_RESOURCE_GROUP`: Resource group name
- `AZURE_TENANT_ID`: Azure tenant ID
- And all custom environment variables defined in your azd environment

## Development

### Prerequisites

- Go 1.23 or later
- golangci-lint
- Node.js 20+ (for cspell)

### Building

```bash
cd cli
./build.sh
```

### Running Tests

```bash
cd cli
go test ./...
```

### Running Linters

```bash
cd cli
golangci-lint run
```

### Running Spell Check

```bash
npm install -g cspell
cspell "**/*.{go,md,yaml,yml}" --config cspell.json
```

### Running Security Scan

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
cd cli
gosec ./...
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **CI Workflow**: Runs on every PR and push to main
  - Linting with golangci-lint
  - Spell checking with cspell
  - Tests on Linux, Windows, and macOS
  - Security scanning with gosec
  - Code coverage reporting

- **CodeQL Workflow**: Security analysis
  - Runs on push to main and weekly
  - Detects security vulnerabilities

- **Release Workflow**: Automated releases
  - Triggered on version tags (e.g., v0.1.0)
  - Builds for multiple platforms
  - Creates GitHub releases with binaries

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- All tests pass
- Code is properly linted
- Documentation is updated
- Security scans pass

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Release Process (for Maintainers)

This project uses an automated release workflow powered by `azd x` commands:

### Creating a Release

1. Go to **Actions** → **Release** workflow
2. Click **Run workflow**
3. Choose bump type:
   - **patch**: Bug fixes (0.1.0 → 0.1.1)
   - **minor**: New features (0.1.0 → 0.2.0)
   - **major**: Breaking changes (0.1.0 → 1.0.0)
4. Click **Run workflow**

The workflow automatically:
- Calculates next version
- Updates `cli/extension.yaml` and `cli/CHANGELOG.md`
- Commits version bump
- Builds binaries for all platforms
- Packages extension archives
- Creates GitHub release
- Updates `registry.json`

### Manual Release (for testing)

```bash
# Install azd extensions tooling
azd extension install microsoft.azd.extensions

# Build for all platforms
cd cli
export extension_id="jongio.azd.exec"
export extension_version="0.1.0"
azd x build --all

# Package extension
azd x pack

# Create release (test)
azd x release --repo "jongio/azd-exec" --version "0.1.0" --draft

# Update registry
azd x publish --registry ../registry.json --version "0.1.0"
```

## Related Projects

- [Azure Developer CLI](https://github.com/Azure/azure-dev)
- [azd-app](https://github.com/jongio/azd-app)

## Support

For issues, questions, or contributions, please use the [GitHub Issues](https://github.com/jongio/azd-exec/issues) page.