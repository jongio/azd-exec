## [0.2.25] - 2026-01-08

- Merge branch 'main' of https://github.com/jongio/azd-exec (540fd09)
- chore: update release workflow for Go setup and versioning process; enhance documentation for CLI reference (4ffb145)
- chore: update registry for v0.2.24 (7807f9a)

## [0.2.24] - 2026-01-08

- refactor: remove dashboard build steps from build scripts (4ca2e2f)
- fix: handle missing dashboard directory in build scripts (07101f7)
- chore: bump version to 0.2.23 (960cbdb)
- feat: add dynamic title to Layout component and improve theme handling logic (a7cef66)
- fix: disable auto-open of Playwright HTML report (9e0c01d)
- refactor: rename src/cmd/script to src/cmd/exec across entire codebase (de4bb82)
- chore: bump version to 0.2.22 (766510d)
- feat: add release-test workflow and use azd-app build scripts (5782bb6)
- fix: rewrite build.sh to handle azd x build environment variables like build.ps1 (182178c)

## [0.2.23] - 2026-01-08

- feat: add dynamic title to Layout component and improve theme handling logic (a7cef66)
- fix: disable auto-open of Playwright HTML report (9e0c01d)
- refactor: rename src/cmd/script to src/cmd/exec across entire codebase (de4bb82)
- chore: bump version to 0.2.22 (766510d)
- feat: add release-test workflow and use azd-app build scripts (5782bb6)
- fix: rewrite build.sh to handle azd x build environment variables like build.ps1 (182178c)

## [0.2.22] - 2026-01-08

- feat: add release-test workflow and use azd-app build scripts (5782bb6)
- fix: rewrite build.sh to handle azd x build environment variables like build.ps1 (182178c)

## [0.2.21] - 2026-01-08

- debug: add comprehensive artifact logging to release workflow (27e9913)

## [0.2.20] - 2026-01-08

- fix: build binaries with platform-specific names directly instead of copying (7f99e9d)

## [0.2.19] - 2026-01-08

- refactor: enhance multi-platform build process by consolidating platform-specific logic and improving output naming (3764603)

## [0.2.18] - 2026-01-08

- fix: normalize paths in test for macOS symlink compatibility (1049f48)
- refactor: streamline CI workflows by removing Node.js setup and pnpm installation; enhance release process with multi-platform builds and artifact uploads (de3f47f)
- fix: add missing words to cspell configuration for improved spell checking (707fe28)
- fix: reorder import statements for consistency in version integration tests (4c9263c)
- feat: enhance Key Vault reference resolution with new akvs format and improve error handling (2a39d11)

## [0.2.17] - 2025-12-24

- Update perms (30f6b3e)
- chore: update CI workflows to improve Go setup and add Node.js and pnpm installation (f767da2)

## [0.2.16] - 2025-12-24

- feat: add local test script to simulate the release workflow (9e97f20)
- fix: use azd x build instead of manual go build commands to match azd-app pattern (5e47771)

## [0.2.15] - 2025-12-23

- fix: remove --artifacts flag from azd x release - it auto-discovers packaged files (e125fe3)

## [0.2.14] - 2025-12-23

- fix: add artifacts flag to azd x release (818543a)

## [0.2.13] - 2025-12-23

- fix: place binaries directly in bin/ folder, not subdirectories (77f75c3)

## [0.2.12] - 2025-12-23

- fix: use full extension ID pattern for binary names (jongio-azd-exec-{os}-{arch}) (59c1a36)

## [0.2.11] - 2025-12-23



## [0.2.10] - 2025-12-23

- fix: use correct binary naming for azd x pack (namespace-based) (7744157)

## [0.2.9] - 2025-12-23



## [0.2.8] - 2025-12-23

- feat: update binary naming for multi-platform builds to match expected format (d9392c7)
- refactor: remove Node.js setup and website build steps from release workflow (a580dc4)

## [0.2.7] - 2025-12-23

- feat: enhance build process to support multi-platform binaries and add verification steps (bdd71b9)

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
