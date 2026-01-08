# azd exec — Key Vault resolution error handling

## Problem
Today, Key Vault reference resolution is effectively all-or-nothing: if any secret reference fails to resolve, the resolution step aborts and execution proceeds with the *original* environment values (e.g., `akvs://...`), even if some secrets could have been resolved successfully.

This blocks common workflows where some secrets are intentionally not present yet (e.g., first-run provisioning), but other required secrets already exist.

## Goals
- Allow `azd exec` to continue resolving Key Vault references even when some fail.
- Replace successfully resolved references with their secret values.
- For failed references, preserve the original reference string (e.g., `akvs://...` or `@Microsoft.KeyVault(...)`) or optionally omit them.
- Emit warnings for failures without aborting script execution.
- Make “continue on error” the default behavior.
- Provide an opt-in flag for strict/fail-fast behavior.

## Non-Goals
- Adding new Key Vault reference formats.
- Changing authentication behavior (still uses DefaultAzureCredential).
- Persisting resolved secret values anywhere (no caching beyond existing in-memory behavior).

## UX / CLI
Default behavior changes to per-variable continue-on-error.

Add a new flag to `azd exec`:

- `--stop-on-keyvault-error`
  - When set, Key Vault resolution becomes fail-fast (abort resolution on the first failure).
  - This restores the prior all-or-nothing behavior for users who rely on strict failure.

Optional follow-on flag (requested alternative):

- `--ignore-missing-secrets`
  - When set, if a secret cannot be resolved specifically due to “not found” / missing secret, omit that environment variable from the executed process environment.
  - Other error types (auth denied, invalid reference, network errors) continue to behave according to whether `--stop-on-keyvault-error` is enabled.
  - If implemented, document interaction rules:
    - `--ignore-missing-secrets` implies default continue-on-error behavior for missing-secret cases.

## Behavior
### Default (no flags)
- Evaluate each environment variable independently.
- For each Key Vault reference:
  - On success: set env var to the secret value.
  - On failure: keep the original reference string (do not blank it).
  - Emit a warning containing:
    - env var name
    - vault/secret identifier (do not log secret value)
    - a short error summary / category
- Continue resolving subsequent variables.
- Script execution continues.

### With `--stop-on-keyvault-error`
- Abort resolution on the first failure and return an error.
- The command should behave like today’s all-or-nothing behavior (no partial substitution), unless implementation constraints require otherwise.

### With `--ignore-missing-secrets` (if implemented)
- If the failure is “secret not found”:
  - Omit the variable from the process environment.
  - Emit warning.
- If the failure is other categories:
  - Follow default behavior unless `--stop-on-keyvault-error` is also provided.

## Logging
- Warnings go to stderr (same channel as other warnings).
- Warnings must not include secret values.

## Compatibility
- Breaking change: default behavior becomes continue-on-error with warnings.
- Flag names should be stable and discoverable via `--help` and CLI reference docs.

## Acceptance Criteria
- By default, a mix of resolvable and unresolvable references results in partial resolution:
  - resolvable variables are replaced with values
  - unresolvable variables remain as their original Key Vault reference string
  - command does not abort resolution and does not fail solely due to one missing secret
- Warnings are emitted for each failed reference (or aggregated in a readable way).
- Unit tests cover:
  - mixed success/failure scenarios
  - ensuring successful values are still applied
  - ensuring failures do not revert previously-resolved values
  - optional: missing-secret omission behavior if `--ignore-missing-secrets` is implemented
- Docs updated:
  - README Key Vault error-handling section
  - cli/docs/cli-reference.md flags and Key Vault section
