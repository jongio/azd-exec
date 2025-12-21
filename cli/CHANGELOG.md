# Changelog

All notable changes to the azd exec extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.18] - 2025-12-20

- fix: split registry update and git commit into separate steps (25b3f62)

## [0.1.17] - 2025-12-20

- fix: use local mode for azd x publish to avoid GitHub API timing issues (d45ebf0)

## [0.1.16] - 2025-12-20

- fix: update registry path in publish step to match azd-app pattern (1e502b1)

## [0.1.15] - 2025-12-20

- fix: add wait and both token env vars for publish step (9c7d2b5)

## [0.1.14] - 2025-12-20

- fix: rename binaries to format expected by azd x pack (e1f2c90)
- chore: bump version to 0.1.13 (f7d5582)
- fix: use build.sh with BUILD_ALL=true for cross-platform binaries (f317875)
- chore: bump version to 0.1.12 (b2624c5)
- fix: remove --all flag from build to build cross-platform (230b82b)
- chore: bump version to 0.1.11 (0d48fc6)
- debug: add logging to see what build and pack produce (593f5fa)
- chore: bump version to 0.1.10 (0e2b398)
- fix: add EXTENSION_ID and EXTENSION_VERSION env vars to pack step (63588a6)
- chore: bump version to 0.1.9 (bd91f2e)
- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.13] - 2025-12-20

- fix: use build.sh with BUILD_ALL=true for cross-platform binaries (f317875)
- chore: bump version to 0.1.12 (b2624c5)
- fix: remove --all flag from build to build cross-platform (230b82b)
- chore: bump version to 0.1.11 (0d48fc6)
- debug: add logging to see what build and pack produce (593f5fa)
- chore: bump version to 0.1.10 (0e2b398)
- fix: add EXTENSION_ID and EXTENSION_VERSION env vars to pack step (63588a6)
- chore: bump version to 0.1.9 (bd91f2e)
- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.12] - 2025-12-20

- fix: remove --all flag from build to build cross-platform (230b82b)
- chore: bump version to 0.1.11 (0d48fc6)
- debug: add logging to see what build and pack produce (593f5fa)
- chore: bump version to 0.1.10 (0e2b398)
- fix: add EXTENSION_ID and EXTENSION_VERSION env vars to pack step (63588a6)
- chore: bump version to 0.1.9 (bd91f2e)
- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.11] - 2025-12-20

- debug: add logging to see what build and pack produce (593f5fa)
- chore: bump version to 0.1.10 (0e2b398)
- fix: add EXTENSION_ID and EXTENSION_VERSION env vars to pack step (63588a6)
- chore: bump version to 0.1.9 (bd91f2e)
- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.10] - 2025-12-20

- fix: add EXTENSION_ID and EXTENSION_VERSION env vars to pack step (63588a6)
- chore: bump version to 0.1.9 (bd91f2e)
- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.9] - 2025-12-20

- fix: use gh release create with direct registry artifacts (b8eef7f)
- chore: bump version to 0.1.8 (3356471)
- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.8] - 2025-12-20

- fix: use file pattern instead of directory for artifacts (65045c7)
- chore: bump version to 0.1.7 (ebb045a)
- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.7] - 2025-12-19

- fix: add GH_TOKEN to registry update step (d0bb54f)
- chore: bump version to 0.1.6 (38d2d4d)
- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.6] - 2025-12-19

- fix: use GitHub release download for registry publish instead of local artifacts (40adcaf)
- chore: bump version to 0.1.5 (4df5355)
- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.5] - 2025-12-19

- chore: bump version to 0.1.4 (610fe10)
- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.4] - 2025-12-19

- fix: use glob pattern for artifacts path in publish step (3e2fe54)
- chore: bump version to 0.1.3 (3bca254)
- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.3] - 2025-12-19

- fix: specify artifacts path for release and publish commands (69b560c)
- chore: bump version to 0.1.2 (58abbf4)
- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.2] - 2025-12-19

- fix: format main.go properly (a47e535)
- fix: remove executeScript tests incompatible with current main.go structure (0f2932e)
- feat: add tests for script execution and improve .gitignore entries (9030650)
- chore: bump version to 0.1.1 (bb7a5e8)
- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

## [0.1.1] - 2025-12-19

- Build azd extension for script execution with azd context (#1) (d819c6a)
- Initial commit (cee63df)

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

[0.1.0]: https://github.com/jongio/azd-exec/releases/tag/v0.1.0
