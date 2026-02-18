## [0.3.6] - 2026-02-18

- Update azd-core to v0.5.0 (#35) (f670ee6)

## [0.3.5] - 2026-02-17

- feat: add copilot skills support and Go 1.26.0 (#34) (2d5e740)

## [0.3.4] - 2026-02-15

- Remove extension preview/alpha requirement (#33) (0443ca1)

## [0.3.3] - 2026-02-10

- Enhance Build and Setup functions for improved local installation process (#32) (667d771)

## [0.3.2] - 2026-02-08

- chore: update usage examples for shell execution in extension.yaml and registry.json (e0d6970)
- chore: update registry for v0.3.1 (982a0ee)

## [0.3.1] - 2026-02-05

- chore: update dependencies and add new changelog entries for upcoming releases (#31) (2acc836)

## [0.3.0] - 2026-01-30



## [0.2.33] - 2026-01-14



## [0.2.32] - 2026-01-13

- chore: upgrade azd-core to v0.3.0 and update related documentation (#28) (7a6ba00)

## [0.2.31] - 2026-01-11

- Integrate azd-core v0.2.0 and Extract Testing Utilities (#26) (2f465d6)

## [0.2.30] - 2026-01-09

- Integrate azd-core library for Key Vault reference resolution (#25) (a9deb44)

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