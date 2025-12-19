# Integration Tests

This document describes the integration test suite for azd-exec.

## Overview

Integration tests verify that the `azd exec` command works correctly with real shells and scripts across different platforms.

## Structure

```
cli/
├── src/
│   ├── internal/executor/
│   │   ├── executor_test.go         # Integration tests (//go:build integration)
│   │   └── executor_unit_test.go    # Unit tests
│   └── cmd/script/commands/
│       ├── run_integration_test.go
│       └── version_integration_test.go
├── tests/projects/
│   ├── bash/                        # Bash test scripts
│   ├── powershell/                  # PowerShell test scripts
│   └── python/                      # Python test scripts
```

## Running Integration Tests

### All Integration Tests

```bash
cd cli

# Run all integration tests
mage testIntegration

# Or with go directly
go test -tags=integration -v ./src/...
```

### Specific Package

```bash
# Test only executor package
TEST_PACKAGE=executor mage testIntegration

# Test only commands package
TEST_PACKAGE=commands mage testIntegration
```

### Specific Test

```bash
# Run specific test by name
TEST_NAME=TestRunCommandIntegration mage testIntegration

# Run specific subtest
TEST_NAME=TestRunCommandIntegration/Bash_simple_script mage testIntegration
```

### With Timeout

```bash
# Override default 10m timeout
TEST_TIMEOUT=15m mage testIntegration
```

## Test Categories

### Executor Integration Tests

**File:** `src/internal/executor/executor_test.go`

Tests the core script execution engine:
- Shell detection from file extension
- Shell detection from shebang
- Real script execution
- Environment variable passing
- Working directory handling

**Example:**
```go
//go:build integration
// +build integration

func TestExecuteIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Create and execute real script
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

### Command Integration Tests

**Files:**
- `src/cmd/script/commands/run_integration_test.go`
- `src/cmd/script/commands/version_integration_test.go`

Tests the CLI commands end-to-end:
- Command parsing and execution
- Argument passing to scripts
- Output format handling (default, json, quiet)
- Error handling for invalid scripts

**Example:**
```go
func TestRunCommandIntegration(t *testing.T) {
    testProjectsDir := getTestProjectsDir(t)
    scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")
    
    outputFormat := "default"
    cmd := NewRunCommand(&outputFormat)
    cmd.SetArgs([]string{scriptPath})
    
    err := cmd.Execute()
    if err != nil {
        t.Errorf("Command failed: %v", err)
    }
}
```

## Test Projects

### Purpose

The `tests/projects/` directory contains real scripts used by integration tests. These scripts verify that:
- Shell/interpreter detection works correctly
- Arguments are passed properly
- Environment variables are inherited
- Scripts execute on different platforms

### Available Test Scripts

#### Bash Scripts (`tests/projects/bash/`)

- **simple.sh** - Basic output test
  ```bash
  #!/bin/bash
  echo "Hello from bash script"
  exit 0
  ```

- **with-args.sh** - Argument passing test
  ```bash
  #!/bin/bash
  echo "Arg 1: $1"
  echo "Arg 2: $2"
  ```

- **env-test.sh** - Environment variable test
  ```bash
  #!/bin/bash
  echo "PATH exists: ${PATH:+yes}"
  ```

#### PowerShell Scripts (`tests/projects/powershell/`)

- **simple.ps1** - Basic output test
  ```powershell
  Write-Host "Hello from PowerShell script"
  exit 0
  ```

- **with-params.ps1** - Parameter passing test
  ```powershell
  param([string]$Name = "World")
  Write-Host "Hello, $Name!"
  ```

- **env-test.ps1** - Environment variable test
  ```powershell
  Write-Host "PATH exists: $($null -ne $env:PATH)"
  ```

#### Python Scripts (`tests/projects/python/`)

- **simple.py** - Basic output test
  ```python
  #!/usr/bin/env python3
  print("Hello from Python script")
  sys.exit(0)
  ```

- **with-args.py** - Argument passing test
  ```python
  #!/usr/bin/env python3
  print(f"Arg 1: {sys.argv[1]}")
  ```

### Using Test Projects in Tests

```go
func getTestProjectsDir(t *testing.T) string {
    t.Helper()
    
    cwd, _ := os.Getwd()
    testDir := filepath.Join(cwd, "..", "..", "..", "tests", "projects")
    testDir = filepath.Clean(testDir)
    
    if _, err := os.Stat(testDir); os.IsNotExist(err) {
        t.Fatalf("Test projects directory not found: %s", testDir)
    }
    
    return testDir
}

func TestWithTestProject(t *testing.T) {
    testProjectsDir := getTestProjectsDir(t)
    scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")
    
    // Use script in test...
}
```

## Platform-Specific Testing

### Windows

- PowerShell scripts: Always available
- Bash scripts: Skipped (unless Git Bash or WSL available)
- Cmd scripts: Always available

```go
func TestBashScript(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Skip("Skipping bash test on Windows")
    }
    // ...
}
```

### Unix (Linux/macOS)

- Bash scripts: Always available
- PowerShell scripts: Requires pwsh (PowerShell Core)
- Python scripts: Usually available

```go
func TestPowerShellScript(t *testing.T) {
    if runtime.GOOS != "windows" {
        if _, err := exec.LookPath("pwsh"); err != nil {
            t.Skip("pwsh not available")
        }
    }
    // ...
}
```

## Writing New Integration Tests

### 1. Add Build Tag

Always include the integration build tag:

```go
//go:build integration
// +build integration

package commands
```

### 2. Skip in Short Mode

Allow tests to be skipped when running unit tests:

```go
func TestYourIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    // ...
}
```

### 3. Use Test Projects

Reference scripts from `tests/projects/`:

```go
testProjectsDir := getTestProjectsDir(t)
scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")
```

### 4. Handle Platform Differences

Skip tests on incompatible platforms:

```go
if runtime.GOOS == "windows" {
    t.Skip("Skipping bash test on Windows")
}
```

### 5. Clean Up Resources

Use `t.TempDir()` or defer cleanup:

```go
tmpDir := t.TempDir() // Auto-cleaned
// or
defer os.RemoveAll(tmpDir) // Manual cleanup
```

## CI/CD Integration

### GitHub Actions

Integration tests run on:
- Pull requests
- Main branch commits
- Release workflow

```yaml
- name: Run integration tests
  run: |
    cd cli
    go test -tags=integration -v ./src/...
```

### Test Matrix

Tests run on multiple platforms:
- **Windows:** PowerShell scripts only
- **Linux:** Bash and Python scripts
- **macOS:** Bash, Python, and pwsh (if available)

## Debugging Integration Tests

### Verbose Output

```bash
# Run with verbose output
go test -tags=integration -v ./src/... -run TestRunCommandIntegration
```

### Single Test

```bash
# Run one specific test
go test -tags=integration -v ./src/cmd/script/commands -run TestRunCommandIntegration/Bash_simple_script
```

### Log Script Output

Add logging to see script output:

```go
cmd := exec.Command("bash", scriptPath)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
err := cmd.Run()
```

### Check Script Existence

Verify test scripts are found:

```go
scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")
if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
    t.Fatalf("Script not found: %s", scriptPath)
}
```

## Common Issues

### Test Projects Not Found

**Problem:** `Test projects directory not found: ...`

**Solution:**
- Ensure you're running from `cli/` directory
- Check that `tests/projects/` exists
- Verify path calculation in `getTestProjectsDir()`

### Scripts Not Executable

**Problem:** Permission denied when executing script

**Solution:**
```bash
# Make scripts executable on Unix
chmod +x cli/tests/projects/bash/*.sh
chmod +x cli/tests/projects/python/*.py
```

### Shell Not Available

**Problem:** `exec: "bash": executable file not found`

**Solution:**
- Skip test on platforms without the shell
- Install missing shell (e.g., `apt-get install bash`)
- Use platform-appropriate skip logic

### Test Timeouts

**Problem:** Integration tests timeout

**Solution:**
```bash
# Increase timeout
go test -tags=integration -timeout=15m ./src/...

# Or with mage
TEST_TIMEOUT=15m mage testIntegration
```

## Best Practices

1. **Always use build tags** - Separate integration from unit tests
2. **Skip gracefully** - Use `t.Skip()` for unavailable shells
3. **Use test projects** - Reference scripts in `tests/projects/`
4. **Clean up resources** - Use `t.TempDir()` or defer cleanup
5. **Test cross-platform** - Handle Windows vs Unix differences
6. **Keep tests fast** - Even integration tests should complete quickly
7. **Verify assumptions** - Check that scripts exist before using them

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Build Constraints](https://pkg.go.dev/go/build#hdr-Build_Constraints)
- [Integration Testing in Go](https://www.ardanlabs.com/blog/2019/10/integration-testing-in-go.html)
