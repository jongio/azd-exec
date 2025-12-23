## [0.2.6] - 2025-12-23

- feat: add test scripts for azd x release flow and update release workflow to build all binaries (cf24401)
- feat: add verification steps for binaries and packages in release workflow (2b42cf5)

## [0.2.5] - 2025-12-23



## [0.2.4] - 2025-12-23



## [0.2.3] - 2025-12-23



## [0.2.2] - 2025-12-23



## [0.2.1] - 2025-12-23

- feat: update GitHub token for release workflow to use RELEASE_PAT (bdbf836)
- feat: add permissions for actions in release workflow (00c6594)
- feat: update Open Graph image for azd exec (515ee8e)
- feat: add script to generate Open Graph image for azd exec (a349920)
- Initial commit (5ecc1c5)

# Changelog

All notable changes to the azd exec extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-12-19

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
- Backward compatibility alias (`azd script` to `azd exec`)

### Features
- `azd exec` - Execute script files with azd context
- `azd exec version` - Display extension version

[0.1.0]: https://github.com/jongio/azd-exec/releases/tag/v0.1.0
