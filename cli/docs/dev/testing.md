# Testing Guide

This guide covers testing strategies and practices for the azd-exec extension.

## Test Categories

### Unit Tests

**Purpose:** Test individual functions and components in isolation

**Characteristics:**
- Fast execution (<1 second per test)
- No external dependencies
- No file I/O or network calls
- Run on every build

**Location:** `*_test.go` or `*_unit_test.go`

**Running:**
```bash
# From cli directory
mage test

# Or directly with go
go test -short ./src/...
```

**Example:**
```go
func TestDetectShellUnit(t *testing.T) {
    tests := []struct {
        name       string
        scriptPath string
        want       string
    }{
        {
            name:       "PowerShell script",
            scriptPath: "test.ps1",
            want:       "powershell",
        },
    }
    
    exec := New(Config{})
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := exec.detectShell(tt.scriptPath)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

**Purpose:** Test real script execution with actual shells and interpreters

**Characteristics:**
- Slower execution (seconds to minutes)
- Uses real file system
- Executes actual scripts
- Requires shells/interpreters to be installed
- Tagged with `//go:build integration`

**Location:** `*_integration_test.go`

**Running:**
```bash
# Run all integration tests
mage testIntegration

# Run specific package
TEST_PACKAGE=executor mage testIntegration
TEST_PACKAGE=commands mage testIntegration

# Run specific test
TEST_NAME=TestRunCommandIntegration mage testIntegration

# Or directly with go
go test -tags=integration -v ./src/...
```

**Example:**
```go
//go:build integration
// +build integration

package commands

func TestExecuteIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    tmpDir := t.TempDir()
    scriptPath := filepath.Join(tmpDir, "test.sh")
    os.WriteFile(scriptPath, []byte("#!/bin/bash\necho hello"), 0700)
    
    exec := New(Config{})
    err := exec.Execute(context.Background(), scriptPath)
    if err != nil {
        t.Errorf("Execute failed: %v", err)
    }
}
```

## Test Projects

The `tests/projects/` directory contains real scripts for integration testing:

```
tests/projects/
├── bash/
│   ├── simple.sh          # Basic output
│   ├── with-args.sh       # Argument passing
│   └── env-test.sh        # Environment variables
├── powershell/
│   ├── simple.ps1         # Basic output
│   ├── with-params.ps1    # Parameters
│   └── env-test.ps1       # Environment variables
└── python/
    ├── simple.py          # Basic output
    └── with-args.py       # Arguments
```

**Usage in tests:**
```go
func TestRunCommandIntegration(t *testing.T) {
    testProjectsDir := filepath.Join("..", "..", "..", "tests", "projects")
    scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")
    
    cmd := NewRunCommand(&outputFormat)
    cmd.SetArgs([]string{scriptPath})
    
    err := cmd.Execute()
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
}
```

## Test Coverage

### Viewing Coverage

```bash
# Generate coverage report
mage testCoverage

# View in browser
# Windows
start coverage/coverage.html
# macOS
open coverage/coverage.html
# Linux
xdg-open coverage/coverage.html
```

### Coverage Goals

- **Overall target:** >=80% coverage
- **Critical paths:** 100% coverage (executor, command handling)
- **Error handling:** All error paths tested

### Checking Coverage

```bash
# Run with coverage
go test -short -coverprofile=coverage.out ./src/...

# View summary
go tool cover -func=coverage.out

# View by package
go tool cover -func=coverage.out | grep executor
```

## Writing Good Tests

### Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestDetectShell(t *testing.T) {
    tests := []struct {
        name       string
        scriptPath string
        want       string
    }{
        {"Bash script", "test.sh", "bash"},
        {"PowerShell script", "test.ps1", "powershell"},
        {"Python script", "test.py", "python3"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Naming

- Use descriptive names: `TestExecute_WithArguments_Success`
- Table test names: Use the `name` field for clarity
- Integration tests: Suffix with `Integration`

### Test Cleanup

Use `t.TempDir()` for automatic cleanup:

```go
func TestWithFile(t *testing.T) {
    tmpDir := t.TempDir() // Auto-cleaned up after test
    scriptPath := filepath.Join(tmpDir, "test.sh")
    // ...
}
```

Or defer cleanup:

```go
func TestManualCleanup(t *testing.T) {
    tmpFile, err := os.CreateTemp("", "test-*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile.Name())
    // ...
}
```

### Error Testing

Test both success and failure cases:

```go
func TestExecute(t *testing.T) {
    tests := []struct {
        name    string
        script  string
        wantErr bool
    }{
        {"Valid script", "#!/bin/bash\nexit 0", false},
        {"Invalid script", "#!/bin/bash\nexit 1", true},
    }
    // ...
}
```

### Platform-Specific Tests

Skip tests on incompatible platforms:

```go
func TestBashScript(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Skip("Skipping bash test on Windows")
    }
    // ...
}
```

## CI/CD Testing

Tests run automatically on:
- **Pull Requests:** Unit + integration tests
- **Main Branch:** Full test suite + coverage
- **Release:** All tests + cross-platform validation

### GitHub Actions

The `.github/workflows/release.yml` runs:
1. Unit tests (`go test -short`)
2. Integration tests (`go test -tags=integration`)
3. Linting (`golangci-lint run`)
4. Build verification (all platforms)

## Troubleshooting

### Tests Fail Locally But Pass in CI

- Check Go version: `go version` (should be 1.23+)
- Clean build cache: `go clean -testcache`
- Check for platform-specific issues
- Verify test dependencies installed

### Integration Tests Timeout

- Increase timeout: `TEST_TIMEOUT=15m mage testIntegration`
- Or: `go test -tags=integration -timeout=15m ./src/...`

### Coverage Too Low

1. Identify uncovered code: `go tool cover -func=coverage.out | grep -v "100.0%"`
2. Add tests for uncovered functions
3. Test error paths and edge cases
4. Consider if code is testable (refactor if needed)

### Test Flakiness

- Use `t.TempDir()` for temp files
- Avoid hardcoded delays (use channels/waitgroups)
- Make tests deterministic (no random data)
- Clean up resources properly

## Best Practices

1. **Fast by default:** Keep unit tests fast (<1s total)
2. **Isolation:** Tests should not depend on each other
3. **Readability:** Test code should be simple and clear
4. **Coverage:** Aim for >=80%, focus on critical paths
5. **Documentation:** Comment complex test setups
6. **Maintainability:** Keep test code clean and DRY

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Coverage](https://go.dev/blog/cover)
- [Advanced Testing](https://golang.org/doc/code)
