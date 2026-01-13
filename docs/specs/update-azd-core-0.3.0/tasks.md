# Tasks: Bump azd-core to v0.3.0

Overview
--------

Update `azd-exec` and `azd-app` to use `github.com/jongio/azd-core v0.3.0`.

Tasks
-----

1. Update module references
   - In `c:\code\azd-exec` and `c:\code\azd-app`, update `go.mod` or any `require` statements to `github.com/jongio/azd-core v0.3.0`.
   - If a root `go.work` references a replace directive or specific version, update it accordingly.

2. Run dependency maintenance
   - Run `go mod tidy` in each repository.

3. Build and test
   - Run `go build ./...` in each repo.
   - Run the test suite: `go test ./...` and fix failures.

4. Address regressions
   - If compile or test failures occur, apply minimal fixes to restore compatibility.
   - If API changes in `azd-core` require non-trivial changes, open a separate issue and document required follow-ups.

5. Update changelogs and PR
   - Add an entry to `CHANGELOG.md` in each repo describing the bump.
   - Open a PR titled `chore: bump azd-core to v0.3.0` against the `main` branch with test results and CI green.

Acceptance Criteria
-------------------

- PRs opened for `azd-exec` and `azd-app` with CI passing.
- Local builds and tests pass after merging (or green CI on PR).

Notes
-----

- Repos paths: `c:\code\azd-exec` and `c:\code\azd-app`.
- Module path: `github.com/jongio/azd-core`
