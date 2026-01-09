---
title: Test Projects
description: Integration test scripts for azd exec command
lastUpdated: 2026-01-09
tags: [testing, integration-tests, scripts]
---

# Test Projects

This directory contains test scripts for integration testing of the `azd exec` command.

## Structure

- `bash/` - Bash shell scripts
  - `simple.sh` - Basic output test
  - `with-args.sh` - Argument passing test
  - `env-test.sh` - Environment variable test

- `powershell/` - PowerShell scripts
  - `simple.ps1` - Basic output test
  - `with-params.ps1` - Parameter passing test
  - `env-test.ps1` - Environment variable test

- `python/` - Python scripts
  - `simple.py` - Basic output test
  - `with-args.py` - Argument passing test

## Usage

These scripts are used by integration tests to verify that `azd exec` correctly:
- Detects the appropriate shell/interpreter
- Executes scripts with proper environment
- Passes arguments correctly
- Handles different platforms (Windows, Linux, macOS)

## Running Manually

```bash
# From cli directory
cd cli

# Run a test script directly
azd exec ../tests/projects/bash/simple.sh

# Run with arguments
azd exec ../tests/projects/bash/with-args.sh arg1 arg2

# Run PowerShell script
azd exec ../tests/projects/powershell/simple.ps1

# Run Python script
azd exec ../tests/projects/python/simple.py
```

## Adding New Test Scripts

When adding new test scripts:
1. Make them executable on Unix: `chmod +x script.sh`
2. Include appropriate shebang line
3. Add clear comments explaining the test purpose
4. Use `exit 0` for success, non-zero for failure
5. Keep output simple and predictable for test assertions
