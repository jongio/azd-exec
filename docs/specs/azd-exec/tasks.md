<!-- NEXT: 3 -->
# azd-exec Tasks

## TODO: Code Review
**Assigned**: Developer (cr)
**Number**: 2
**Description**: Comprehensive code review covering security, logic, types, errors, tests, and performance. Prioritize issues. Compare against existing code patterns in codebase.

## TODO: Refactor
**Assigned**: Developer (rf)
**Number**: 3
**Description**: Identify and fix duplication, large files, dead code, and magic numbers. Ensure tests pass after changes.

## TODO: Fix & Verify
**Assigned**: Developer (fix)
**Number**: 4
**Description**: Build, run tests, and fix all errors. Repeat until compilation and all tests pass cleanly.

## IN PROGRESS

### 2. Remove legacy run command ✓
**Assigned**: Developer
**Completed**: 2025-12-20
**Description**: Removed run.go and its associated tests since the command was simplified to direct execution. Updated README and all website pages to remove 'run' references and show the cleaner 'azd exec ./script.sh' syntax throughout. Build and tests verified passing.

## DONE

### 1. Add Key Vault Reference Resolution ✓
**Assigned**: Developer
**Completed**: 2025-12-20
**Description**: Implemented automatic resolution of Key Vault references in environment variables. Supports both SecretUri and VaultName+SecretName formats. Includes Azure SDK integration, authentication via DefaultAzureCredential, comprehensive error handling, unit/integration tests, and full documentation.
