package testhelpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// CaptureOutput captures stdout during function execution.
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

	// Channel for output
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
	w.Close()
	os.Stdout = origStdout

	// Get output
	output := <-outCh

	if fnErr != nil {
		t.Logf("Command error: %v", fnErr)
	}

	return output
}

// GetTestProjectsDir finds the test projects directory.
func GetTestProjectsDir(t *testing.T) string {
	t.Helper()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Try multiple possible paths
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

	// If not found, try to find the cli directory
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
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
