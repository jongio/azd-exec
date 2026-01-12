# Upgrade azd-core to v0.3.0

Summary
-------

Upgrade consumer repositories to use `github.com/jongio/azd-core` version `v0.3.0`.

Why
---

- `azd-core` `v0.3.0` has been released and contains improvements and bug fixes relied upon by downstream projects.
- Keep dependent repositories on a recent, supported release to receive fixes and API improvements.

Scope
-----

- Update `azd-exec` and `azd-app` to require `github.com/jongio/azd-core v0.3.0`.
- Ensure `go.work`, `go.mod` and module references are updated where applicable.
- Run `go mod tidy`, build and run tests in both repos and fix regressions if any.

Acceptance Criteria
-------------------

- Both `azd-exec` and `azd-app` build successfully after the upgrade.
- Tests that previously passed remain passing; any test regressions are either fixed or documented with a follow-up task.
- PRs opened in both repositories titled `chore: bump azd-core to v0.3.0` with changelog and CI passing.

Owner
-----

Developer (code owner for each repo)
