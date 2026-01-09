// Package testhelpers provides common testing utilities for the azd exec extension.
// It includes helpers for capturing output, locating test resources, and common assertions.
package testhelpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// CaptureOutput captures stdout during function execution.
// It redirects os.Stdout to a pipe, executes the function, and returns the captured output.
// The original stdout is always restored, even if the function returns an error.
// This is useful for testing commands that write to stdout.
//
// Example:
//
//	output := CaptureOutput(t, func() error {
//	    fmt.Println("test output")
//	    return nil
//	})
func CaptureOutput(t *testing.T, fn func() error) string {
	t.Helper()

	// Save original stdout
	origStdout := os.Stdout

	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Replace stdout
	os.Stdout = w

	// Channel for output (buffered to avoid goroutine leak)
	outCh := make(chan string, 1)
	go func() {
		var output strings.Builder
		buf := make([]byte, 1024)
		for {
			n, readErr := r.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if readErr != nil {
				break
			}
		}
		outCh <- output.String()
	}()

	// Execute function
	fnErr := fn()

	// Close write end and restore stdout
	if err := w.Close(); err != nil {
		t.Logf("Failed to close pipe writer: %v", err)
	}
	os.Stdout = origStdout

	// Get output
	output := <-outCh

	if fnErr != nil {
		t.Logf("Command error: %v", fnErr)
	}

	return output
}

// GetTestProjectsDir finds the test projects directory relative to the current working directory.
// It searches common locations and returns the first valid path found.
// This helper is useful for integration tests that need to locate test fixtures.
// It fails the test if the directory cannot be found.
func GetTestProjectsDir(t *testing.T) string {
	t.Helper()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Try multiple possible paths relative to common test locations
	possiblePaths := []string{
		filepath.Join(cwd, "..", "..", "..", "tests", "projects"),       // From commands dir
		filepath.Join(cwd, "..", "tests", "projects"),                   // From src dir
		filepath.Join(cwd, "tests", "projects"),                         // From cli dir
		filepath.Join(cwd, "..", "..", "..", "..", "tests", "projects"), // From nested test dir
	}

	for _, testDir := range possiblePaths {
		testDir = filepath.Clean(testDir)
		if _, err := os.Stat(testDir); err == nil {
			return testDir
		}
	}

	// If not found, try to walk up the directory tree
	testDir := filepath.Join(cwd, "tests", "projects")
	for i := 0; i < 5; i++ {
		testDir = filepath.Join("..", testDir)
		testDir = filepath.Clean(testDir)
		if _, err := os.Stat(testDir); err == nil {
			return testDir
		}
	}

	t.Fatalf("Test projects directory not found. CWD: %s", cwd)
	return ""
}

// Contains checks if a string contains a substring.
// This is a convenience helper for common test assertions.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
