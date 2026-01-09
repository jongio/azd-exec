---
title: azd-core Integration
description: Migrate Key Vault resolution to shared azd-core library
lastUpdated: 2026-01-09
tags: [spec, integration, azd-core, keyvault]
status: implemented
---

# azd-core integration — Migrate Key Vault resolution to shared library

## Context

`azd exec` currently implements Key Vault reference detection and resolution in-repo under [cli/src/internal/executor](../../../cli/src/internal/executor) with custom parsing, client caching, and env map/slice handling.

The new `azd-core` utility repo provides reusable Go packages for common Azure Developer CLI operations:
- **Local path**: `c:\code\azd-core`
- **GitHub**: https://github.com/jongio/azd-core
- **Relevant packages**:
  - `keyvault`: Key Vault reference detection, resolution, DefaultAzureCredential, thread-safe client caching
  - `env`: Environment variable resolution with configurable error handling (graceful vs fail-fast)

**Goal**: Replace bespoke Key Vault resolution in `azd-exec` with `azd-core` packages to reduce duplication, keep behavior parity, and ease future maintenance.

## Goals

- Use `azd-core/keyvault` + `azd-core/env` for all Key Vault detection/resolution in `azd exec`
- Preserve current CLI behavior and flags:
  - Default continue-on-error with warnings
  - `--stop-on-keyvault-error` restores fail-fast/all-or-nothing
- Keep supported reference formats unchanged (SecretUri, VaultName/SecretName, `akvs://`)
- Maintain or improve test coverage (unit + integration) for resolution and CLI flags
- Minimize user-visible changes; update docs only where needed

## Non-Goals

- Adding new Key Vault reference formats or authentication mechanisms
- Changing CLI surface beyond existing flags/behavior
- Introducing secret value caching beyond existing in-memory behavior
- Modifying azd-core library itself (use as-is)

## Discovery Phase

### Current azd-exec Implementation

Review existing keyvault implementation:
- [cli/src/internal/executor/keyvault.go](../../../cli/src/internal/executor/keyvault.go) - Main resolver
- [cli/src/internal/executor/keyvault_env.go](../../../cli/src/internal/executor/keyvault_env.go) - Env integration
- [cli/src/internal/executor/keyvault_patterns.go](../../../cli/src/internal/executor/keyvault_patterns.go) - Pattern detection
- Related test files

**Key behaviors to preserve**:
- Supported formats: SecretUri, VaultName+SecretName, `akvs://` protocol
- Default: continue-on-error with warnings to stderr
- Flag: `--stop-on-keyvault-error` for fail-fast behavior
- Error messages must not include secret values
- Thread-safe client caching

### azd-core Analysis

Analyze `c:\code\azd-core` packages:
- API surface of `keyvault` package
- API surface of `env` package
- Configuration options (fail-fast vs continue)
- Error handling patterns
- Client caching implementation
- Test patterns and coverage

**Questions to answer**:
1. Does azd-core support all three reference formats?
2. How does azd-core handle partial resolution failures?
3. What configuration options exist for error handling?
4. Are warning messages compatible?
5. How is client caching implemented?
6. What are the Go module dependencies?

## Integration Plan

### Phase 1: Add azd-core Dependency

1. Add `azd-core` as Go module dependency in [cli/go.mod](../../../cli/go.mod)
2. Run `go mod tidy` to resolve dependencies
3. Verify no version conflicts

### Phase 2: Replace Executor Resolution

1. Update executor to use `azd-core` packages:
   - Import `azd-core/keyvault` and `azd-core/env`
   - Replace custom resolver instantiation with `keyvault.NewKeyVaultResolver`
   - Wire CLI flags to azd-core options (fail-fast vs continue)
   - Update env resolution calls to use `env.ResolveMap`/`ResolveSlice`

2. Ensure compatibility:
   - Env normalization/parsing matches azd-exported values
   - Warning format aligns with current stderr output
   - Error categories map correctly

### Phase 3: Cleanup

1. Remove custom keyvault resolver code:
   - Mark deprecated or remove `keyvault.go` custom types
   - Remove duplicated regex/constants if covered by azd-core
   - Remove redundant helper functions

2. Update imports and remove unused dependencies

### Phase 4: Test Migration

1. Update unit tests:
   - Adjust mocks/fixtures for azd-core APIs
   - Validate mixed success/failure scenarios
   - Test `--stop-on-keyvault-error` flag behavior
   - Verify warning output format

2. Update integration tests:
   - Ensure tests still gate correctly (require Azure credentials)
   - Test actual Key Vault resolution with azd-core
   - Verify performance/caching behavior

3. Maintain or improve coverage targets (>=80%)

### Phase 5: Documentation

1. Update [README.md](../../../README.md):
   - Note azd-core usage for Key Vault resolution
   - Confirm behavior unchanged
   - Add troubleshooting notes if needed

2. Update [cli/docs/cli-reference.md](../../../cli/docs/cli-reference.md):
   - Document `--stop-on-keyvault-error` flag
   - Reference azd-core library
   - Update examples if needed

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Behavior drift** | Warning messages or error categories may differ from current implementation | Assert message shape in tests; align options; document any intentional changes |
| **Performance regression** | Client caching may behave differently | Add perf sanity check in tests; verify cache hit rates |
| **API surface mismatch** | azd-core may not support all current features (e.g., ignore-missing) | Identify gaps early; capture as follow-up rather than blocking migration |
| **Dependency conflicts** | azd-core deps may conflict with existing modules | Test early; coordinate with azd-core maintainers if issues arise |
| **Breaking changes** | azd-core API may change | Pin to specific version; track azd-core releases |

## Acceptance Criteria

- [x] Executor uses `azd-core` `keyvault`/`env` for resolution
- [x] Custom resolver code is no longer invoked (deprecated or removed)
- [x] Default behavior: partial resolution with warnings; stop-on-error flag restores fail-fast
- [x] Supported formats unchanged (SecretUri, VaultName+SecretName, `akvs://`)
- [x] Env normalization matches current behavior
- [x] Tests updated and passing:
  - [x] Unit tests cover mixed success/failure and flag paths
  - [x] Integration tests validate actual Key Vault resolution
  - [x] Coverage maintained or improved (>=80%) - **83.9% achieved**
- [x] Documentation updated:
  - [x] README reflects azd-core usage
  - [x] CLI reference updated
  - [x] Help text accurate
- [x] Build passes with new dependency
- [x] No regression in performance or user experience

## Success Metrics

- Zero user-visible breaking changes
- Test coverage >= 80% maintained
- Build time does not increase significantly
- Code reduction in azd-exec (removed custom resolver code)
- Easier maintenance through shared library

## Follow-up Items

Items discovered during migration that are out of scope:

- [ ] TBD based on discovery phase findings
