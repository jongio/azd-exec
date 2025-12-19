//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binaryName    = "exec"
	srcDir        = "src/cmd/script"
	binDir        = "bin"
	coverageDir   = "coverage"
	extensionFile = "extension.yaml"
	extensionID   = "jongio.azd.exec"
	testTimeout   = "10m"
)

// Default target runs all checks and builds.
var Default = All

// All runs format, lint, test, and build in dependency order.
func All() error {
	mg.Deps(Fmt, Lint, Test, Build)
	return nil
}

// Build compiles the CLI binary and installs it locally.
func Build() error {
	fmt.Println("Building azd exec extension...")

	binaryExt := ""
	if runtime.GOOS == "windows" {
		binaryExt = ".exe"
	}

	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	outputPath := filepath.Join(binDir, binaryName+binaryExt)

	// Build the binary
	if err := sh.RunV("go", "build", "-o", outputPath, "./"+srcDir); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("✅ Built successfully: %s\n", outputPath)

	// Install using azd x build
	fmt.Println("Installing extension locally...")
	if err := sh.RunV("azd", "x", "build"); err != nil {
		return fmt.Errorf("azd x build failed: %w", err)
	}

	fmt.Println("✅ Extension installed successfully!")
	return nil
}

// Install builds and installs the extension (alias for Build).
func Install() error {
	return Build()
}

// Test runs unit tests only (with -short flag).
func Test() error {
	fmt.Println("Running unit tests...")
	return sh.RunV("go", "test", "-v", "-short", "./src/...")
}

// TestIntegration runs integration tests only.
// Set TEST_PACKAGE env var to filter by package (e.g., executor, commands)
// Set TEST_NAME env var to run a specific test
// Set TEST_TIMEOUT env var to override default timeout
func TestIntegration() error {
	fmt.Println("Running integration tests...")

	args := []string{"test", "-v", "-tags=integration"}

	// Handle timeout
	timeout := os.Getenv("TEST_TIMEOUT")
	if timeout == "" {
		timeout = testTimeout
	}
	args = append(args, "-timeout="+timeout)

	// Handle test filtering
	testName := os.Getenv("TEST_NAME")
	if testName != "" {
		args = append(args, "-run="+testName)
	}

	// Handle package filtering
	pkg := os.Getenv("TEST_PACKAGE")
	testPath := "./src/..."
	if pkg != "" {
		switch pkg {
		case "executor":
			testPath = "./src/internal/executor"
		case "commands":
			testPath = "./src/cmd/script/commands"
		default:
			return fmt.Errorf("unknown package: %s (valid: executor, commands)", pkg)
		}
	}
	args = append(args, testPath)

	return sh.RunV("go", args...)
}

// TestAll runs all tests (unit + integration).
func TestAll() error {
	fmt.Println("Running all tests...")
	return sh.RunV("go", "test", "-v", "-tags=integration", "./src/...")
}

// TestCoverage runs tests with coverage report.
func TestCoverage() error {
	fmt.Println("Running tests with coverage...")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create absolute path for coverage directory
	absCoverageDir := filepath.Join(cwd, coverageDir)

	// Remove existing coverage directory
	_ = os.RemoveAll(absCoverageDir)

	// Create coverage directory
	if err := os.MkdirAll(absCoverageDir, 0755); err != nil {
		return fmt.Errorf("failed to create coverage directory: %w", err)
	}

	coverageOut := filepath.Join(absCoverageDir, "coverage.out")
	coverageHTML := filepath.Join(absCoverageDir, "coverage.html")

	// Run tests with coverage (use -short to skip integration tests)
	cmd := exec.Command("go", "test", "-short", "-coverprofile="+coverageOut, "./src/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	// Generate HTML report
	if err := sh.RunV("go", "tool", "cover", "-html="+coverageOut, "-o", coverageHTML); err != nil {
		return fmt.Errorf("failed to generate coverage HTML: %w", err)
	}

	// Calculate coverage percentage
	cmd = exec.Command("go", "tool", "cover", "-func="+coverageOut)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to calculate coverage: %w", err)
	}

	fmt.Println("\n" + string(output))
	fmt.Printf("✅ Coverage report generated: %s\n", coverageHTML)

	// Check if coverage meets threshold
	outputStr := string(output)
	if strings.Contains(outputStr, "total:") {
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if strings.Contains(line, "total:") {
				fmt.Println("\n📊 " + strings.TrimSpace(line))
				break
			}
		}
	}

	return nil
}

// Fmt formats all Go code.
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.RunV("go", "fmt", "./...")
}

// Lint runs go vet and golangci-lint.
func Lint() error {
	fmt.Println("Running linters...")

	// Run go vet
	fmt.Println("Running go vet...")
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return fmt.Errorf("go vet failed: %w", err)
	}

	// Check if golangci-lint is available
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Println("⚠️  golangci-lint not found, skipping")
		fmt.Println("   Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	} else {
		fmt.Println("Running golangci-lint...")
		if err := sh.RunV("golangci-lint", "run"); err != nil {
			return fmt.Errorf("golangci-lint failed: %w", err)
		}
	}

	fmt.Println("✅ Linting passed!")
	return nil
}

// Clean removes build artifacts and coverage reports.
func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	dirs := []string{binDir, coverageDir}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove %s: %w", dir, err)
		}
	}

	fmt.Println("✅ Clean complete!")
	return nil
}

// Uninstall removes the locally installed extension.
func Uninstall() error {
	fmt.Println("Uninstalling extension...")

	if err := sh.RunV("azd", "extension", "uninstall", extensionID); err != nil {
		fmt.Println("⚠️  Extension may not be installed")
		return nil
	}

	fmt.Println("✅ Extension uninstalled!")
	return nil
}

// Watch monitors files and rebuilds on changes (requires azd x watch).
func Watch() error {
	fmt.Println("Starting watch mode...")
	return sh.RunV("azd", "x", "watch")
}

// Preflight runs all checks before shipping: format, build, lint, tests, and coverage.
func Preflight() error {
	fmt.Println("🚀 Running preflight checks...\n")

	checks := []struct {
		name string
		fn   func() error
	}{
		{"Format check", Fmt},
		{"Linting", Lint},
		{"Unit tests", Test},
		{"Integration tests", TestIntegration},
		{"Coverage report", TestCoverage},
	}

	for i, check := range checks {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(checks), check.name)
		fmt.Println(strings.Repeat("=", 50))
		if err := check.fn(); err != nil {
			return fmt.Errorf("❌ %s failed: %w", check.name, err)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ All preflight checks passed!")
	fmt.Println("🚀 Ready to ship!")
	return nil
}
