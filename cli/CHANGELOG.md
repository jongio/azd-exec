## [0.2.28] - 2026-01-08

### Added
- Initial public release of azd exec
- Execute scripts and commands with full access to azd environment variables
- Azure Key Vault reference resolution (`@keyvault(vaultName, secretName)` and `@akvs(secretId)` formats)
- Automatic shell detection (bash, sh, zsh, pwsh, powershell, cmd)
- Script argument support with `--` separator
- Working directory control with `--cwd` flag
- Interactive mode for scripts requiring user input
- Multi-platform support (Windows, macOS, Linux, ARM64)
- Comprehensive test coverage with unit and integration tests
- Security scanning with CodeQL and gosec
- Automated release workflow and package distribution
- Complete documentation website with examples
- Backward compatibility alias (`azd script` continues to work)

### Features
- `azd exec <script>` - Execute script files with azd context
- `azd exec version` - Display extension version
- `--cwd <path>` - Set working directory for script execution
- `--stop-on-keyvault-error` - Fail fast on Key Vault resolution errors (default: continue)
- `-- <args>` - Pass arguments to your scripts

# Changelog

All notable changes to the azd exec extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).