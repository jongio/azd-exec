<!-- NEXT:  -->
# azd-exec Tasks

## DONE

### 5. Continue-on-error Key Vault resolution ✓
**Assigned**: Developer
**Completed**: 2026-01-08
**Description**: Default Key Vault reference resolution now continues after individual failures (successful secrets substituted; failures kept as original `akvs://`/`@Microsoft.KeyVault(...)` values) with warnings. Added `--stop-on-keyvault-error` for fail-fast/all-or-nothing behavior. Updated docs (README + cli/docs/cli-reference.md) and tests.

### 4. Preflight Verification ✓
**Assigned**: Developer (pf)
**Completed**: 2025-12-22
**Description**: Successfully executed preflight checks. Fixed spelling issue (added USERPROFILE to cspell.json), fixed all linting issues (error handling, godot comments, gocritic suggestions, thread-safety improvements). All 6 checks passed: format ✓, spell ✓, lint ✓, unit tests ✓ (100% pass), integration tests ✓, coverage ✓ (68.8%). Ready to ship.

### 4. Fix & Verify ✓
**Assigned**: Developer (fix)
**Completed**: 2025-12-22
**Description**: Completed comprehensive build and test verification. All systems green: CLI builds cleanly, 100% test pass rate, 81% coverage (exceeds 80% target), zero security issues (gosec), website builds successfully. Production-ready with no outstanding issues.

### 3. Refactor ✓
**Assigned**: Developer (rf)
**Completed**: 2025-12-22
**Description**: Fixed all HIGH and MEDIUM priority issues from code review. Modified 6 files: main.go (context cancellation), command_builder.go (shell validation), keyvault_env.go (thread-safety with mutex), build.ps1 (removed unused var), detect_shell.go (error handling + magic number). All 47 tests passing, builds cleanly.

### 2. Code Review ✓
**Assigned**: Developer (cr)
**Completed**: 2025-12-22
**Description**: Completed comprehensive code review. Found 0 critical, 3 high, 5 medium, 4 low priority issues. Overall assessment: 4/5 stars, production-ready. Key findings: excellent security practices (gosec/CodeQL passing), 86%+ test coverage, strong architecture. High-priority issues: context cancellation, shell validation, thread-safety. Full report generated with prioritized action plan.

### 1. Add Inline Script Execution ✓
**Assigned**: Developer
**Completed**: 2025-12-20
**Description**: Implemented inline script execution with `azd exec 'echo foo'` syntax. Added ExecuteInline method to executor, modified main.go to detect inline vs file scripts, supports all shells (bash, pwsh, cmd). Updated README with examples and security warnings. Includes comprehensive tests for both unit and integration scenarios. Build and tests verified passing.

### 2. Remove legacy run command ✓
**Assigned**: Developer
**Completed**: 2025-12-20
**Description**: Removed run.go and its associated tests since the command was simplified to direct execution. Updated README and all website pages to remove 'run' references and show the cleaner 'azd exec ./script.sh' syntax throughout. Build and tests verified passing.

### 1. Add Key Vault Reference Resolution ✓
**Assigned**: Developer
**Completed**: 2025-12-20
**Description**: Implemented automatic resolution of Key Vault references in environment variables. Supports both SecretUri and VaultName+SecretName formats. Includes Azure SDK integration, authentication via DefaultAzureCredential, comprehensive error handling, unit/integration tests, and full documentation.
