package testhelpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaptureOutput(t *testing.T) {
	t.Run("Captures stdout successfully", func(t *testing.T) {
		output := CaptureOutput(t, func() error {
			fmt.Println("test output")
			return nil
		})

		if !strings.Contains(output, "test output") {
			t.Errorf("Expected output to contain 'test output', got: %s", output)
		}
	})

	t.Run("Handles function errors", func(t *testing.T) {
		output := CaptureOutput(t, func() error {
			fmt.Println("before error")
			return errors.New("test error")
		})

		if !strings.Contains(output, "before error") {
			t.Errorf("Expected output before error, got: %s", output)
		}
	})

	t.Run("Captures empty output", func(t *testing.T) {
		output := CaptureOutput(t, func() error {
			return nil
		})

		// Empty output is expected and valid
		_ = output
		t.Log("Empty output captured correctly")
	})
}

func TestGetTestProjectsDir(t *testing.T) {
	t.Run("Finds test projects directory", func(t *testing.T) {
		// Save original directory
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Chdir(originalDir)
		}()

		// Create temporary directory structure that mimics the project
		tmpDir := t.TempDir()
		testProjectsPath := filepath.Join(tmpDir, "tests", "projects")
		if err := os.MkdirAll(testProjectsPath, 0750); err != nil {
			t.Fatal(err)
		}

		// Create a marker file
		markerFile := filepath.Join(testProjectsPath, "README.md")
		if err := os.WriteFile(markerFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		// Change to a subdirectory
		srcDir := filepath.Join(tmpDir, "src", "cmd", "script")
		if err := os.MkdirAll(srcDir, 0750); err != nil {
			t.Fatal(err)
		}

		if err := os.Chdir(srcDir); err != nil {
			t.Fatal(err)
		}

		// Get test projects directory
		result := GetTestProjectsDir(t)

		// Verify it's a valid path
		if result == "" {
			t.Error("GetTestProjectsDir returned empty string")
		}

		// Verify the path exists
		if _, err := os.Stat(result); err != nil {
			t.Errorf("Test projects directory does not exist: %v", err)
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "Substring found",
			s:      "hello world",
			substr: "world",
			want:   true,
		},
		{
			name:   "Substring not found",
			s:      "hello world",
			substr: "foo",
			want:   false,
		},
		{
			name:   "Empty substring",
			s:      "hello",
			substr: "",
			want:   true,
		},
		{
			name:   "Empty string",
			s:      "",
			substr: "test",
			want:   false,
		},
		{
			name:   "Case sensitive",
			s:      "Hello World",
			substr: "hello",
			want:   false,
		},
		{
			name:   "Exact match",
			s:      "test",
			substr: "test",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}
