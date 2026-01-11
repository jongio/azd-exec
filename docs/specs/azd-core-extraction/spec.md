# azd-core Extraction Opportunities - azd-exec

## Goal
Identify components in azd-exec that could be extracted to azd-core for reuse by other azd extensions.

## Background

azd-exec currently uses azd-core for:
- ✅ **shellutil**: Shell detection from file extensions and shebangs
- ✅ **keyvault**: Azure Key Vault reference resolution

This spec evaluates remaining components in azd-exec for potential extraction to azd-core.

## Analysis Summary

| Component | Extract? | Priority | Rationale |
|-----------|----------|----------|-----------|
| testhelpers | ✅ Yes | **P1** | Generic test utilities, high reuse potential |
| version pattern | ⚠️ Maybe | **P2** | Common pattern, but simple to duplicate |
| errors | ⚠️ Maybe | **P3** | Some generic types, but many are executor-specific |
| command_builder | ❌ No | - | Executor-specific implementation |
| executor core | ❌ No | - | Extension-specific business logic |

---

## Extraction Candidates

### 1. testhelpers Package ✅ RECOMMEND EXTRACTION

**Priority**: **P1** (High Value)

**Current Location**: `azd-exec/cli/src/internal/testhelpers`

**Capabilities**:
- `CaptureOutput(t, fn)`: Capture stdout during test execution
- `GetTestProjectsDir(t)`: Locate test fixtures with smart path searching
- `Contains(s, substr)`: String containment helper

**Why Extract**:
- **High Reuse Potential**: Any extension with CLI commands needs output capture for testing
- **Generic Implementation**: Nothing azd-exec specific, pure test infrastructure
- **Quality of Life**: Simplifies writing reliable CLI tests across extensions
- **Pattern Standardization**: Consistent test helpers across azd ecosystem

**Proposed azd-core Package**: `testutil`

**Functions to Extract**:
```go
// CaptureOutput captures stdout during function execution for testing
func CaptureOutput(t *testing.T, fn func() error) string

// FindTestData locates test fixture directories relative to current working directory
// Renamed from GetTestProjectsDir to be more generic
func FindTestData(t *testing.T, subdirs ...string) string

// TempDir creates a temporary directory for testing with automatic cleanup
func TempDir(t *testing.T) string

// Contains checks if string contains substring (convenience helper)
func Contains(s, substr string) bool
```

**Benefits**:
- **azd-app**: Can use CaptureOutput for CLI command tests, FindTestData for locating test fixtures
- **azd-exec**: Standardized test infrastructure
- **Future Extensions**: Consistent testing patterns from day one

**Test Coverage Target**: ≥85%

**Implementation Notes**:
- `GetTestProjectsDir` should be generalized to `FindTestData` with configurable subdirectories
- Add `TempDir` helper with automatic cleanup via `t.Cleanup()`
- Consider adding `AssertContains(t, s, substr)` for cleaner test assertions

**Risks**: Low - Pure test utilities with no production dependencies

---

### 2. Version Pattern ⚠️ CONSIDER

**Priority**: **P2** (Medium Value)

**Current Location**: `azd-exec/cli/src/internal/version`

**Capabilities**:
- Version, BuildDate, GitCommit variables set via ldflags
- ExtensionID constant for registry identification
- Name constant for human-readable extension name

**Why Consider**:
- **Common Pattern**: All extensions need version management
- **Build-time Configuration**: Standardized ldflags approach
- **Registry Integration**: ExtensionID pattern for extension registry

**Why Not Extract**:
- **Simple to Duplicate**: Only ~30 lines, each extension can implement easily
- **Extension-Specific Constants**: ExtensionID and Name are unique per extension
- **Limited Shared Logic**: Mostly just variable declarations

**Proposed Alternative**: **Documentation Pattern** instead of code extraction

**Recommendation**: 
- **Do NOT extract code** - too simple, extension-specific
- **DO document pattern** in azd-core README with example:

```go
// Example: Extension version management pattern
package version

var Version = "0.0.0-dev"      // Set via: -ldflags "-X .../version.Version=1.0.0"
var BuildDate = "unknown"       // Set via: -ldflags "-X .../version.BuildDate=2026-01-10T12:00:00Z"
var GitCommit = "unknown"       // Set via: -ldflags "-X .../version.GitCommit=abc123"

const ExtensionID = "your.extension.id"  // Must match extension.yaml
const Name = "Your Extension Name"
```

**Benefits of Documentation Approach**:
- Extensions maintain flexibility
- No dependency on azd-core for simple version management
- Consistent pattern without coupling

---

### 3. Error Types ⚠️ PARTIAL EXTRACTION

**Priority**: **P3** (Low-Medium Value)

**Current Location**: `azd-exec/cli/src/internal/executor/errors.go`

**Error Types**:
1. `ValidationError` - ✅ **Generic, extract**
2. `ScriptNotFoundError` - ⚠️ **Specific to script execution**
3. `InvalidShellError` - ❌ **Specific to shell execution**
4. `ExecutionError` - ❌ **Specific to script execution**

**Why Partial Extraction**:
- `ValidationError` is generic and useful for any extension that validates input
- Other errors are executor-specific and unlikely to be reused

**Proposed azd-core Package**: `errors` (new package)

**Extract Only**:
```go
// ValidationError indicates that input validation failed
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for %s: %s", e.Field, e.Reason)
}
```

**Additional Generic Errors to Add**:
```go
// NotFoundError indicates a resource was not found
type NotFoundError struct {
    ResourceType string  // e.g., "file", "environment", "configuration"
    Name         string
}

// ExecutionError indicates a command/operation failed
type ExecutionError struct {
    Operation string
    ExitCode  int
    Stderr    string
}
```

**Why This Matters**:
- Standardized error types improve error handling across extensions
- Easier to write tests that check error types
- Consistent error messages for better user experience

**Test Coverage Target**: ≥90% (error types are simple but critical)

**Risks**: Low - Just error type definitions, no complex logic

---

## Components NOT to Extract

### command_builder.go ❌
**Reason**: Executor-specific logic for building shell commands with different flag patterns (-c, -File, /c). This is core business logic for azd-exec, not a reusable utility.

### executor.go ❌
**Reason**: Extension-specific business logic. Orchestrates script execution, environment resolution, Key Vault integration. This is the "app logic" not a utility library.

### constants.go ❌
**Reason**: Shell identifiers (shellBash, shellPwsh, etc.) already extracted to azd-core/shellutil. Remaining constants are executor-specific.

---

## Extraction Priority & Roadmap

### Phase 1: High-Value Extraction (v0.3.0)
**Target**: azd-core v0.3.0

1. **testutil package** (P1)
   - Extract `CaptureOutput`, `FindTestData`, add `TempDir`
   - Write comprehensive tests (≥85% coverage)
   - Update azd-exec to use azd-core/testutil
   - Update azd-app to use azd-core/testutil for CLI tests

**Impact**: 
- Code reduction: ~100 lines across azd-exec
- Enables better testing in azd-app CLI commands
- Standardizes test infrastructure for future extensions

### Phase 2: Error Standardization (v0.3.0 or v0.4.0)
**Target**: azd-core v0.3.0 or v0.4.0 (depending on demand)

2. **errors package** (P3)
   - Extract `ValidationError`
   - Add `NotFoundError`, `ExecutionError` as generic types
   - Update azd-exec and azd-app to use shared error types
   - Write tests (≥90% coverage)

**Impact**:
- Modest code reduction (50-75 lines)
- Improved error handling consistency
- Better error type checking in tests

### Phase 3: Documentation (v0.3.0)
**Non-code deliverable**

3. **Extension Patterns Guide**
   - Document version management pattern
   - Document extension structure best practices
   - Provide examples from azd-exec

**Impact**:
- Easier for new extensions to follow established patterns
- No code coupling, maximum flexibility

---

## Success Criteria

### Phase 1 Complete When:
- ✅ azd-core/testutil package published with ≥85% coverage
- ✅ azd-exec migrated to use azd-core/testutil (remove internal/testhelpers)
- ✅ azd-app uses testutil for CLI command tests
- ✅ Documentation updated with usage examples

### Phase 2 Complete When:
- ✅ azd-core/errors package published with ≥90% coverage
- ✅ azd-exec migrated to use ValidationError from azd-core
- ✅ azd-app migrated to use standardized error types
- ✅ Documentation updated with error handling examples

### Phase 3 Complete When:
- ✅ Extension patterns guide published in azd-core/docs
- ✅ Version management pattern documented
- ✅ Extension structure examples added

---

## Integration Impact

### azd-exec Changes (Phase 1)
**Files Affected**:
- `cli/src/internal/testhelpers/` → Delete (move to azd-core)
- `cli/src/internal/testhelpers/testhelpers_test.go` → Migrate tests to azd-core
- `cli/src/cmd/exec/commands/commands_test.go` → Update imports
- `cli/src/internal/executor/executor_test.go` → Update imports
- All test files using testhelpers → Update imports

**Code Impact**: -100 lines, improved test reliability

### azd-app Benefits (Phase 1)
**New Capabilities**:
- CaptureOutput for testing CLI commands in dashboard/commands
- FindTestData for locating test fixtures
- TempDir for isolated test environments

**Files That Could Use testutil**:
- `cli/src/dashboard/commands/*.go` tests
- `cli/src/internal/workflow/*_test.go`
- Any new CLI commands added to azd-app

**Code Impact**: +0 to +50 lines (enabling better test coverage)

---

## Decision Matrix

| Component | Extract? | Package | Version | Reason |
|-----------|----------|---------|---------|--------|
| testhelpers | ✅ **YES** | testutil | v0.3.0 | High reuse, enables better testing across ecosystem |
| version pattern | ❌ **NO** (doc only) | - | v0.3.0 | Too simple, document pattern instead |
| ValidationError | ⚠️ **MAYBE** | errors | v0.3.0/v0.4.0 | Useful but lower priority, extract if errors package created |
| Other errors | ❌ **NO** | - | - | Executor-specific, limited reuse |
| command_builder | ❌ **NO** | - | - | Extension business logic |
| executor | ❌ **NO** | - | - | Extension business logic |

---

## Recommendations

### Immediate Action (Now)
1. **Approve testutil extraction** for azd-core v0.3.0
2. Create task: "Extract testhelpers to azd-core/testutil"
3. Create task: "Migrate azd-exec to use azd-core/testutil"
4. Create task: "Document extension patterns in azd-core"

### Future Consideration (v0.4.0+)
5. **Evaluate errors package** after v0.3.0 release based on:
   - Demand from azd-app for standardized error types
   - Emergence of additional extensions needing error standardization
   - Community feedback on error handling patterns

### Do Not Extract
- version.go (document pattern instead)
- executor-specific logic (command_builder, executor core)
- executor-specific error types (ScriptNotFoundError, InvalidShellError, ExecutionError)

---

## Related Work

### Already Extracted to azd-core
- ✅ **shellutil** (v0.2.0): Shell detection from extensions/shebangs
- ✅ **keyvault** (v0.1.0): Azure Key Vault reference resolution

### azd-core v0.2.0 Packages
- fileutil, pathutil, browser, security, procutil, shellutil
- See: [c:\code\azd-core\release-notes-v0.2.0.md](file:///c:/code/azd-core/release-notes-v0.2.0.md)

### Integration History
- azd-exec shellutil integration: 349 lines removed, improved reliability
- azd-app fileutil integration: 50 lines removed, fixed critical bug

---

## Open Questions

1. **testutil Scope**: Should testutil include assertion helpers (AssertContains, AssertError) or keep minimal?
   - **Recommendation**: Start minimal (CaptureOutput, FindTestData, TempDir), add assertions in v0.4.0 if needed

2. **errors Package**: Extract now or wait for more demand signals?
   - **Recommendation**: Wait - only ValidationError has clear reuse, others are speculative

3. **Version Pattern**: Code extraction vs documentation?
   - **Recommendation**: Documentation only - too simple to warrant code dependency

---

---

## azd-app Extraction Analysis

Based on analysis of `azd-app/cli/src/internal/`, here are additional extraction candidates:

### High-Value Candidates

#### 1. output Package ✅ RECOMMEND EXTRACTION (P1)

**Priority**: **P1** (High Value)

**Current Location**: `azd-app/cli/src/internal/output`

**Capabilities**:
- **CLI Output Formatting**: Headers, sections, success/error/warning/info messages
- **ANSI Color Support**: Consistent color scheme (green success, red errors, yellow warnings, blue info)
- **Unicode/Emoji Detection**: Automatic fallback to ASCII for terminals without Unicode support
- **JSON Output Mode**: Structured JSON output for automation/scripting
- **Progress Indicators**: Progress bars, spinners, status badges
- **Tables**: Simple table rendering with column alignment
- **Interactive Prompts**: Confirmation dialogs with y/n input
- **Orchestration Mode**: Header suppression for subcommands in workflows

**Why Extract**:
- **Universal CLI Need**: Every extension with CLI commands needs consistent output formatting
- **Cross-Platform Complexity**: Unicode detection (Windows Terminal, VS Code, PowerShell, ConEmu)
- **JSON Mode**: Critical for automation and scripting scenarios
- **Brand Consistency**: All azd extensions should have consistent visual output
- **Rich Feature Set**: ~500 lines of battle-tested formatting logic

**Proposed azd-core Package**: `cliout` (or `output`)

**Key Functions to Extract**:
```go
// Format management
SetFormat(format string) error  // "default" or "json"
GetFormat() Format
IsJSON() bool
SetOrchestrated(bool)  // Skip headers for subcommands

// Output functions
Header(text string)
Section(icon, text string)
Success(format string, args ...interface{})
Error(format string, args ...interface{})
Warning(format string, args ...interface{})
Info(format string, args ...interface{})

// Formatted output
Bullet(format string, args ...interface{})
Label(label, value string)
Table(headers []string, rows []TableRow)
ProgressBar(current, total, width int) string

// Interactive
Confirm(message string) bool

// JSON/Default hybrid
Print(data interface{}, formatter func()) error
PrintJSON(data interface{}) error
```

**Benefits**:
- **azd-exec**: Standardized output for exec commands, better JSON mode support
- **azd-app**: Already has this, can import from azd-core
- **Future Extensions**: Instant professional CLI output out of the box

**Implementation Notes**:
- Move Unicode detection logic (Windows Terminal, VS Code, PowerShell detection)
- Keep all ANSI color constants and emoji/ASCII fallbacks
- Maintain orchestration mode for composed workflows
- Add tests for Windows/Unix Unicode detection

**Test Coverage Target**: ≥80%

**Risks**: Low - Pure output formatting, no business logic dependencies

---

#### 2. constants Package ⚠️ PARTIAL EXTRACTION (P2)

**Priority**: **P2** (Medium Value - Partial)

**Current Location**: `azd-app/cli/src/internal/constants`

**What to Extract**:
- ✅ **File Permissions** (DirPermission, FilePermission) - Already in azd-core/fileutil
- ✅ **HTTP/Network Timeouts** (HTTPIdleConnTimeout, HTTPDialTimeout, etc.) - Generic
- ✅ **Error Limits** (MaxStderrLength, MaxErrorMessageLength) - Generic
- ❌ **Dashboard-specific** constants - azd-app specific
- ❌ **Service/Health status** values - azd-app specific
- ❌ **WebSocket config** - azd-app specific

**Proposed azd-core Package**: `constants` (new package)

**Extract These Only**:
```go
// Network timeouts
const HTTPIdleConnTimeout = 90 * time.Second
const HTTPDialTimeout = 5 * time.Second
const HTTPKeepAliveTimeout = 30 * time.Second
const HTTPTLSHandshakeTimeout = 5 * time.Second
const HTTPExpectContinueTimeout = 1 * time.Second

// Error/output limits
const MaxStderrLength = 10000       // 10KB max stderr capture
const MaxErrorMessageLength = 500   // Max error message before truncation
```

**Why Extract (Partial)**:
- HTTP timeout values are useful for any extension making Azure API calls
- Error limits are good defaults for CLI error handling
- File permissions already extracted to fileutil

**Why Not Extract (Rest)**:
- Dashboard, service, health, WebSocket constants are azd-app specific
- Pattern limits, UI constraints are application-specific

**Benefits**:
- Standardized timeout values across azd ecosystem
- Consistent error handling limits
- Reduces "magic numbers" in extensions

**Risks**: Very low - Just constants, no logic

---

#### 3. logging Package ⚠️ CONSIDER (P2-P3)

**Priority**: **P2-P3** (Medium-Low Value)

**Current Location**: `azd-app/cli/src/internal/logging`

**Capabilities**:
- Structured logging built on `log/slog`
- Component-based loggers with context propagation
- Configurable log levels (Debug, Info, Warn, Error)
- Structured JSON output support
- Service/operation context tracking

**Why Consider**:
- Standardized logging across extensions
- Component-based filtering useful for debugging
- Built on stdlib `log/slog` (Go 1.21+)

**Why Hesitate**:
- Many extensions may prefer custom logging (zerolog, zap, etc.)
- Simple enough that duplicating isn't costly
- Not much shared configuration needed

**Recommendation**: 
- **Do NOT extract now** - logging preferences vary widely
- **Consider for v0.4.0+** if multiple extensions converge on slog-based logging
- **Alternative**: Document logging pattern recommendations in azd-core README

---

### Medium-Value Candidates

#### 4. testing/coverage ⚠️ DOMAIN-SPECIFIC (P3)

**Priority**: **P3** (Low-Medium Value)

**Current Location**: `azd-app/cli/src/internal/testing/coverage.go`

**Capabilities**:
- Coverage aggregation across multiple services
- Cobertura XML report generation
- HTML coverage reports with source highlighting
- JSON coverage reports
- Line-level coverage tracking with hit counts

**Why NOT Extract**:
- **Highly domain-specific**: Designed for azd-app's multi-service architecture
- **Complex dependencies**: Security validation, file I/O, HTML generation
- **Limited reuse**: Most extensions won't aggregate coverage across services
- **Test infrastructure**: More suited to internal testing package than shared utility

**Recommendation**: **Do NOT extract**
- Functionality is specific to azd-app's orchestration model
- Other extensions use standard Go coverage tools
- If needed, extensions can use `go tool cover` directly

---

#### 5. testing/testutil ⚠️ MINIMAL VALUE (P3)

**Priority**: **P3** (Low Value)

**Current Location**: `azd-app/cli/src/internal/testing/testutil/bindloopback.go`

**Capabilities**:
- `ListenLoopback(port)`: Create TCP listener on loopback interface with ephemeral port

**Why NOT Extract**:
- **Too Simple**: 1 function, ~10 lines
- **Standard Library**: `net.Listen("tcp", "127.0.0.1:0")` is nearly as simple
- **Limited Use Case**: Only needed for tests that bind network ports
- **Duplicate Effort**: azd-exec testhelpers didn't need this

**Recommendation**: **Do NOT extract**
- Functionality too simple to warrant extraction
- Extensions can implement directly if needed

---

#### 6. wellknown/services ❌ DOMAIN-SPECIFIC

**Priority**: **N/A** (azd-app specific)

**Current Location**: `azd-app/cli/src/internal/wellknown/services.go`

**Capabilities**:
- Registry of well-known Azure emulator services (Azurite, Cosmos, Redis, PostgreSQL)
- Service definitions with Docker images, ports, health checks
- Connection string templates

**Why NOT Extract**:
- **100% azd-app specific**: Designed for `azd app add` command
- **Domain knowledge**: Azure emulator configuration, not general utility
- **Extension-specific**: No other extension needs service registry

**Recommendation**: **Do NOT extract** - This is azd-app's feature, not shared infrastructure

---

### Summary: azd-app Extraction Priorities

| Component | Extract? | Package | Priority | Lines Saved | Rationale |
|-----------|----------|---------|----------|-------------|-----------|
| **output** | ✅ **YES** | `cliout` | **P1** | ~500 | Universal CLI output, brand consistency, complex Unicode logic |
| **constants** (partial) | ⚠️ **MAYBE** | `constants` | **P2** | ~30 | HTTP timeouts useful, but small scope |
| **logging** | ❌ **NO** (doc only) | - | P2-P3 | 0 | Preferences vary, document pattern instead |
| **testing/coverage** | ❌ **NO** | - | P3 | 0 | Domain-specific to azd-app orchestration |
| **testutil/bindloopback** | ❌ **NO** | - | P3 | 0 | Too simple, standard lib sufficient |
| **wellknown** | ❌ **NO** | - | N/A | 0 | azd-app feature, not shared utility |

---

## Combined Extraction Roadmap (azd-exec + azd-app)

### Phase 1: Essential Utilities (v0.3.0)
**Target**: azd-core v0.3.0

1. **testutil package** (from azd-exec) - **P1**
   - Extract `CaptureOutput`, `FindTestData`, add `TempDir`
   - Migrate azd-exec and azd-app test helpers
   - **Impact**: ~150 lines removed across projects, standardized testing

2. **cliout package** (from azd-app) - **P1**
   - Extract full output package with Unicode detection, JSON mode, formatting
   - Migrate azd-app, add to azd-exec
   - **Impact**: ~500 lines shared, brand consistency across extensions

**Combined Phase 1 Impact**:
- **Code Reduction**: ~650 lines across azd-exec and azd-app
- **Standardization**: Consistent CLI output and testing patterns
- **Quality**: Professional CLI UX out of the box for new extensions

### Phase 2: Error Standardization (v0.3.0 or v0.4.0)

3. **errors package** (from azd-exec) - **P3**
   - Extract `ValidationError` and add `NotFoundError`, `ExecutionError`
   - **Impact**: ~75 lines, improved error handling consistency

4. **constants package** (from azd-app) - **P2**
   - Extract HTTP timeout constants, error limits
   - **Impact**: ~30 lines, standardized network/error handling

### Phase 3: Documentation (v0.3.0)

5. **Extension Patterns Guide**
   - Version management pattern
   - Logging pattern recommendations
   - Extension structure best practices
   - **Impact**: Easier onboarding for new extensions

---

## Updated Success Criteria

### Phase 1 Complete When:
- ✅ azd-core/testutil package published (from azd-exec testhelpers) with ≥85% coverage
- ✅ azd-core/cliout package published (from azd-app output) with ≥80% coverage
- ✅ azd-exec migrated to use both packages
- ✅ azd-app migrated to use azd-core/testutil (already has output package)
- ✅ Documentation updated with usage examples
- ✅ ~650 lines of duplicate code eliminated

### Phase 2 Complete When:
- ✅ azd-core/errors package published with ≥90% coverage
- ✅ azd-core/constants package published with network/error constants
- ✅ azd-exec and azd-app migrated to use standardized error types
- ✅ Documentation updated

### Phase 3 Complete When:
- ✅ Extension patterns guide published in azd-core/docs
- ✅ Version management, logging patterns documented

---

## Final Recommendations

### Extract NOW (v0.3.0):
1. **testutil** (P1) - From azd-exec, enables better testing in azd-app
2. **cliout** (P1) - From azd-app, standardizes CLI UX across ecosystem

### Extract LATER (v0.3.0 or v0.4.0):
3. **errors** (P3) - If standardized error handling gains traction
4. **constants** (P2) - If HTTP timeout standardization proves valuable

### Document (v0.3.0):
5. **Extension Patterns** - Version management, logging recommendations, structure

### Do NOT Extract:
- version.go (azd-exec) - Too simple, document pattern
- logging.go (azd-app) - Preferences vary, document pattern
- coverage.go (azd-app) - Domain-specific
- testutil/bindloopback.go (azd-app) - Too simple
- wellknown/services.go (azd-app) - azd-app feature

---

## References

- azd-core consolidation spec: [c:\code\azd-core\docs\specs\consolidation\spec.md](file:///c:/code/azd-core/docs/specs/consolidation/spec.md)
- azd-core v0.2.0 release: [c:\code\azd-core\release-notes-v0.2.0.md](file:///c:/code/azd-core/release-notes-v0.2.0.md)
- azd-exec testhelpers: [c:\code\azd-exec\cli\src\internal\testhelpers\testhelpers.go](file:///c:/code/azd-exec/cli/src/internal/testhelpers/testhelpers.go)
- azd-exec errors: [c:\code\azd-exec\cli\src\internal\executor\errors.go](file:///c:/code/azd-exec/cli/src/internal/executor/errors.go)
- azd-exec version: [c:\code\azd-exec\cli\src\internal\version\version.go](file:///c:/code/azd-exec/cli/src/internal/version/version.go)
- azd-app output: [c:\code\azd-app\cli\src\internal\output\output.go](file:///c:/code/azd-app/cli/src/internal/output/output.go)
- azd-app constants: [c:\code\azd-app\cli\src\internal\constants\](file:///c:/code/azd-app/cli/src/internal/constants/)
- azd-app logging: [c:\code\azd-app\cli\src\internal\logging\logger.go](file:///c:/code/azd-app/cli/src/internal/logging/logger.go)
