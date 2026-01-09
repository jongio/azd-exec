---
title: Changelog
description: Release history and changes for azd-exec CLI
lastUpdated: 2026-01-09
tags: [changelog, releases, version-history]
---

## [0.2.29] - 2026-01-08

- refactor: update Key Vault reference commands to use 'set-secret' for environment variables (8dd9d28)
- chore: update changelog to version 0.2.29 (fc2895c)
- refactor: remove unnecessary blank line in executor coverage tests (749509f)
- refactor: remove working directory support from azd exec and related documentation (9caec7a)
- Merge branch 'main' of https://github.com/jongio/azd-exec (e2a3120)
- feat: add Azure Key Vault integration examples and documentation across multiple sections (5e45a5d)
- chore: update registry for v0.2.28 (1cf69ad)

## [0.2.29] - 2026-01-08

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