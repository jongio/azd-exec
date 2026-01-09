---
title: azd-core Integration Tasks
description: Task tracking for azd-core library integration
lastUpdated: 2026-01-09
tags: [tasks, integration, tracking]
---

<!-- NEXT: -->
# azd-core Integration Tasks

## TODO

## IN PROGRESS

## DONE

### 6. Update documentation ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Updated README.md and cli/docs/cli-reference.md to note azd-core library usage in Key Vault Integration sections. Confirmed all behaviors unchanged (continue-on-error default, --stop-on-keyvault-error flag, supported formats, authentication methods). No breaking changes to user-facing functionality.

### 5. Update tests ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Migrated unit and integration tests to work with azd-core APIs. Updated tests to rely on public API (IsKeyVaultReference), removed reliance on internal patterns. Added test helper for stop-on-error semantics validation. Unit tests: 100% pass, coverage 83.9% (exceeds 80% target). Integration tests: 100% pass.

### 4. Remove deprecated code ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Replaced entire custom keyvault resolver implementation with thin adapters to azd-core. Removed duplicated regex patterns, normalization logic, client caching, and resolution functions. All code now delegates to github.com/jongio/azd-core/keyvault package.

### 3. Integrate azd-core into executor ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Replaced custom resolver in keyvault.go with type aliases and function adapters to azd-core keyvault package. CLI flags wire correctly to azd-core options. Env normalization, warning format, and error handling preserved through azd-core's implementation.

### 2. Add azd-core dependency ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Added github.com/jongio/azd-core@main to cli/go.mod. Ran `go mod tidy`. No version conflicts detected. Azure SDK dependencies moved to indirect (provided by azd-core).

### 1. Analyze azd-core packages ✓
**Assigned**: Developer
**Completed**: 2026-01-09
**Description**: Analyzed azd-core keyvault and env packages via GitHub API. Confirmed API surface matches azd-exec needs: IsKeyVaultReference, NewKeyVaultResolver, ResolveEnvironmentVariables, ResolveEnvironmentOptions with StopOnError flag. Supports all three reference formats (SecretUri, VaultName+SecretName, akvs://). Thread-safe client caching implemented. Compatible dependencies (same Azure SDK versions).
