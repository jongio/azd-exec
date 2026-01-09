---
title: Code Review Checklist
description: Code review standards and checklist for azd-exec project
lastUpdated: 2026-01-09
tags: [code-review, standards, quality]
---

# Code Review Checklist

This document outlines the code review standards for the azd-exec project.

## Security (CRITICAL)

- [ ] **Input Validation**: All user inputs are validated (file paths, shell names, script content)
- [ ] **Path Traversal**: File paths checked for `..` sequences and normalized
- [ ] **Command Injection**: Shell parameters validated against allowlist
- [ ] **Secret Handling**: No secrets in logs or error messages
- [ ] **Error Messages**: Sanitized to avoid leaking sensitive information (use base filenames, not full paths)

## Error Handling

- [ ] **Error Wrapping**: Errors wrapped with context using `fmt.Errorf` with `%w`
- [ ] **Error Types**: Use specific error types (`ValidationError`, `ScriptNotFoundError`, etc.)
- [ ] **Context Propagation**: `context.Context` passed through for cancellation
- [ ] **Resource Cleanup**: Deferred cleanup (file closes, etc.) with error checking
- [ ] **User-Friendly Messages**: Clear, actionable error messages

## Code Quality

- [ ] **No Magic Numbers**: Constants defined for all magic values
- [ ] **No Duplication**: DRY principle followed
- [ ] **Function Size**: Functions <=50 lines, single responsibility
- [ ] **Naming**: Clear, descriptive names following Go conventions
- [ ] **Comments**: Package-level godoc, complex logic explained

## Testing

- [ ] **Test Coverage**: >=80% coverage for new code
- [ ] **Test Quality**: Tests cover edge cases, error paths
- [ ] **Test Independence**: No test interdependencies
- [ ] **Fast Tests**: Unit tests <1s, use `-short` flag appropriately
- [ ] **Integration Tests**: Marked with `//go:build integration`

## Performance

- [ ] **No Unnecessary Allocations**: Reuse buffers where possible
- [ ] **Efficient Loops**: Avoid nested loops with high complexity
- [ ] **Resource Limits**: Timeouts and limits on external operations
- [ ] **Caching**: Expensive operations cached appropriately

## Documentation

- [ ] **README**: Updated if user-facing changes
- [ ] **CLI Reference**: Generated docs match implementation
- [ ] **Examples**: Working examples for new features
- [ ] **Godoc**: All exported items documented
