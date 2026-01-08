package executor

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteWithInteractiveMode(t *testing.T) {
	exec := New(Config{
		Interactive: true,
	})

	if !exec.config.Interactive {
		t.Error("Interactive mode not set")
	}
}

func TestExecuteWithCustomShell(t *testing.T) {
	exec := New(Config{
		Shell: "bash",
	})

	if exec.config.Shell != "bash" {
		t.Errorf("Shell = %v, want bash", exec.config.Shell)
	}
}

func TestExecuteError(t *testing.T) {
	t.Run("Nonexistent script", func(t *testing.T) {
		exec := New(Config{})
		err := exec.Execute(context.Background(), "nonexistent-script.sh")
		if err == nil {
			t.Error("Expected error for nonexistent script, got nil")
		}
	})

	t.Run("Empty script path", func(t *testing.T) {
		exec := New(Config{})
		err := exec.Execute(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty script path, got nil")
		}
		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("Expected 'cannot be empty' error, got: %v", err)
		}
	})
}

func TestExecuteWithArguments(t *testing.T) {
	args := []string{"arg1", "arg2", "arg3"}
	exec := New(Config{Args: args})

	if len(exec.config.Args) != 3 {
		t.Errorf("Args length = %v, want 3", len(exec.config.Args))
	}

	if exec.config.Args[0] != "arg1" {
		t.Errorf("Args[0] = %v, want arg1", exec.config.Args[0])
	}
}
