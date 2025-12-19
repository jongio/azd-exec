# Changelog

All notable changes to the azd exec extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - TBD

### Added
- Initial release of azd exec extension
- Execute scripts with full access to azd environment variables and context
- Automatic shell detection (bash, sh, zsh, pwsh, powershell, cmd)
- Script arguments support with `--` separator
- Working directory control with `--cwd` flag
- Interactive mode for scripts requiring user input
- Comprehensive unit tests and integration tests
- Security scanning with CodeQL and gosec
- Multi-platform support (Windows, macOS, Linux)
- Backward compatibility alias (`azd script` → `azd exec`)

### Features
- `azd exec run` - Execute script files with azd context
- `azd exec version` - Display extension version
- `azd exec listen` - Internal extension framework communication

[0.1.0]: https://github.com/jongio/azd-exec/releases/tag/v0.1.0
