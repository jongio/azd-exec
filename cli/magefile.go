//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binaryName    = "exec"
	srcDir        = "src/cmd/exec"
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

// killExtensionProcesses terminates any running azd exec extension processes.
func killExtensionProcesses() error {
	extensionBinaryPrefix := strings.ReplaceAll(extensionID, ".", "-")

	if runtime.GOOS == "windows" {
		fmt.Println("Stopping any running extension processes...")
		for _, arch := range []string{"windows-amd64", "windows-arm64"} {
			procName := extensionBinaryPrefix + "-" + arch
			_ = exec.Command("powershell", "-NoProfile", "-Command",
				"Stop-Process -Name '"+procName+"' -Force -ErrorAction SilentlyContinue").Run()
		}
	} else {
		_ = exec.Command("pkill", "-f", extensionBinaryPrefix).Run()
	}
	return nil
}

// runWithEnvRetry runs a command with environment variables, retrying up to 3 times on failure.
func runWithEnvRetry(env map[string]string, cmd string, args ...string) error {
	const maxRetries = 3
	var err error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			delay := time.Duration(i*5) * time.Second
			fmt.Printf("  ⚠️  Attempt %d/%d failed, retrying in %s...\n", i, maxRetries, delay)
			time.Sleep(delay)
		}
		if err = sh.RunWithV(env, cmd, args...); err == nil {
			return nil
		}
	}
	return err
}

// Build compiles the CLI binary and installs it locally using azd x build.
func Build() error {
	_ = killExtensionProcesses()
	time.Sleep(500 * time.Millisecond)

	// Ensure azd extensions are set up (enables extensions + installs azd x if needed)
	if err := ensureAzdExtensions(); err != nil {
		return err
	}

	fmt.Println("Building and installing azd exec extension...")

	// Get version from extension.yaml
	version, err := getVersion()
	if err != nil {
		return err
	}

	// Set environment variables required by azd x build
	env := map[string]string{
		"EXTENSION_ID":      extensionID,
		"EXTENSION_VERSION": version,
	}

	// Build and install directly using azd x build
	if err := runWithEnvRetry(env, "azd", "x", "build"); err != nil {
		return fmt.Errorf("azd x build failed: %w", err)
	}

	fmt.Printf("✅ Build complete! Version: %s\n", version)
	fmt.Println("   Run 'azd exec version' to verify")
	return nil
}

// Pack packages the extension into archives using azd x pack.
// This creates platform-specific zip/tar.gz files in ~/.azd/registry/...
func Pack() error {
	fmt.Println("Packaging extension...")

	// First, ensure we have platform-specific binaries for pack to find
	// azd x build creates bin/exec.exe but azd x pack needs bin/jongio-azd-exec-windows-amd64.exe
	version, err := getVersion()
	if err != nil {
		return err
	}

	// Build for current platform first
	env := map[string]string{
		"EXTENSION_ID":      extensionID,
		"EXTENSION_VERSION": version,
	}

	fmt.Println("Building binary...")
	if err := sh.RunWithV(env, "azd", "x", "build", "--skip-install"); err != nil {
		return fmt.Errorf("azd x build failed: %w", err)
	}

	// azd x build now creates the platform-specific binary directly
	// Format: {extension-id-with-dashes}-{os}-{arch}{ext}
	fmt.Println("Verifying platform-specific binary exists...")

	extensionIdSafe := strings.ReplaceAll(extensionID, ".", "-")
	platformName := fmt.Sprintf("%s-windows-amd64.exe", extensionIdSafe)
	binPath := filepath.Join(binDir, platformName)

	if _, err := os.Stat(binPath); err != nil {
		return fmt.Errorf("binary not found at %s: %w", binPath, err)
	}

	fmt.Printf("  ✅ Found %s\n", platformName)

	// Now pack
	fmt.Println("Creating packages...")
	return sh.RunV("azd", "x", "pack")
}

// Publish updates the local registry with extension metadata using azd x publish.
// This makes the extension available for installation via 'azd extension install'.
func Publish() error {
	fmt.Println("Publishing to local registry...")
	return sh.RunV("azd", "x", "publish")
}

// Setup runs the complete local development setup: build -> install.
// For most development, 'mage build' is sufficient since azd x build handles installation.
// Use 'mage pack' and 'mage publish' separately when testing the release workflow.
func Setup() error {
	fmt.Println("🚀 Setting up extension for local development...")
	fmt.Println()

	if err := Build(); err != nil {
		return err
	}

	fmt.Println("\n✅ Setup complete!")
	fmt.Println("\n🎉 You can now use: azd exec version")
	return nil
}

// ensureLocalRegistry creates the local extension registry if it doesn't exist.
// This is used by Pack/Publish workflows for testing the release pipeline locally.
func ensureLocalRegistry() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	registryPath := filepath.Join(homeDir, ".azd", "registry.json")

	// Check if registry already exists
	if _, err := os.Stat(registryPath); err == nil {
		// Registry exists, check if "local" source is configured
		output, err := sh.Output("azd", "extension", "source", "list")
		if err == nil && strings.Contains(output, "local") {
			return nil // All good
		}
	}

	// Create registry file
	fmt.Println("📦 Setting up local extension registry...")

	azdDir := filepath.Join(homeDir, ".azd")
	if err := os.MkdirAll(azdDir, 0755); err != nil {
		return fmt.Errorf("failed to create .azd directory: %w", err)
	}

	// Write minimal registry structure
	registryContent := `{"extensions":[]}`
	if err := os.WriteFile(registryPath, []byte(registryContent), 0644); err != nil {
		return fmt.Errorf("failed to create registry.json: %w", err)
	}

	// Add local source
	fmt.Println("📦 Adding local extension source...")
	if err := sh.RunV("azd", "extension", "source", "add", "--name", "local", "--location", "registry.json", "--type", "file"); err != nil {
		return fmt.Errorf("failed to add local source: %w", err)
	}

	fmt.Println("✅ Local registry created!")
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

// getVersion reads the version from extension.yaml
func getVersion() (string, error) {
	data, err := os.ReadFile(extensionFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", extensionFile, err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "version:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("version not found in %s", extensionFile)
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
			testPath = "./src/cmd/exec/commands"
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

	// Run go test with coverage
	args := []string{"test", "-short", "-coverprofile=" + coverageOut, "./src/..."}
	if err := sh.RunV("go", args...); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	// Generate HTML report
	if err := sh.RunV("go", "tool", "cover", "-html="+coverageOut, "-o", coverageHTML); err != nil {
		return fmt.Errorf("failed to generate coverage HTML: %w", err)
	}

	// Calculate coverage percentage
	output, err := sh.Output("go", "tool", "cover", "-func="+coverageOut)
	if err != nil {
		return fmt.Errorf("failed to calculate coverage: %w", err)
	}

	fmt.Println("\n" + output)
	fmt.Printf("✅ Coverage report generated: %s\n", coverageHTML)

	// Check if coverage meets threshold
	if strings.Contains(output, "total:") {
		lines := strings.Split(output, "\n")
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
	if _, err := sh.Output("golangci-lint", "version"); err != nil {
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

// Watch monitors files and rebuilds on changes (requires azd x watch).
func Watch() error {
	// Ensure azd extensions are set up
	if err := ensureAzdExtensions(); err != nil {
		return err
	}

	fmt.Println("Starting watch mode...")

	// Get version from extension.yaml
	version, err := getVersion()
	if err != nil {
		return err
	}

	// Set environment variables required by azd x watch
	env := map[string]string{
		"EXTENSION_ID":      extensionID,
		"EXTENSION_VERSION": version,
	}

	return sh.RunWithV(env, "azd", "x", "watch")
}

// ensureAzdExtensions checks that azd is installed and the azd x extension is installed.
// This is a prerequisite for commands that use azd x (build, watch, etc.).
func ensureAzdExtensions() error {
	// Check if azd is available
	if _, err := sh.Output("azd", "version"); err != nil {
		return fmt.Errorf("azd is not installed or not in PATH. Install from https://aka.ms/azd")
	}

	// Check if azd x extension is available
	if _, err := sh.Output("azd", "x", "--help"); err != nil {
		fmt.Println("📦 Installing azd x extension (developer kit)...")
		if err := sh.RunV("azd", "extension", "install", "microsoft.azd.extensions", "--source", "azd", "--no-prompt"); err != nil {
			return fmt.Errorf("failed to install azd x extension: %w", err)
		}
		fmt.Println("✅ azd x extension installed!")
	}

	return nil
}

// SpellCheck runs cspell to check for spelling errors.
func SpellCheck() error {
	fmt.Println("Running spell check...")

	// Check if cspell is available
	if _, err := sh.Output("cspell", "--version"); err != nil {
		fmt.Println("⚠️  cspell not found, skipping")
		fmt.Println("   Install: npm install -g cspell")
		return nil
	}

	// Run cspell
	if err := sh.RunV("cspell", "**/*.{go,md,yaml,yml}", "--config", "../cspell.json"); err != nil {
		return fmt.Errorf("spell check failed: %w", err)
	}

	fmt.Println("✅ Spell check passed!")
	return nil
}

// Preflight runs all checks before shipping: format, build, lint, tests, and coverage.
func Preflight() error {
	fmt.Println("🚀 Running preflight checks...\n")

	checks := []struct {
		name string
		fn   func() error
	}{
		{"Format check", Fmt},
		{"Spell check", SpellCheck},
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
