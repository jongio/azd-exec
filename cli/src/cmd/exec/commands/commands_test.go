package commands

import (
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)

	if cmd == nil {
		t.Fatal("NewVersionCommand returned nil")
	}

	if cmd.Use != "version" {
		t.Errorf("Command Use = %v, want version", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}

func TestVersionCommandDefault(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)

	// Version command uses fmt.Printf which goes to stdout
	// We test by executing and checking error only
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully
}

func TestVersionCommandQuiet(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)
	cmd.SetArgs([]string{"--quiet"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully with quiet flag
}

func TestVersionCommandJSON(t *testing.T) {
	outputFormat := "json"
	cmd := NewVersionCommand(&outputFormat)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully with JSON format
	// Note: Output validation would require capturing stdout which is complex in tests
}

func TestNewListenCommand(t *testing.T) {
	cmd := NewListenCommand()

	if cmd == nil {
		t.Fatal("NewListenCommand returned nil")
	}

	if cmd.Use != "listen" {
		t.Errorf("Command Use = %v, want listen", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}

	if !cmd.Hidden {
		t.Error("Listen command should be hidden")
	}
}

func TestListenCommandExecution(t *testing.T) {
	cmd := NewListenCommand()

	// Listen command requires a running azd server for gRPC communication.
	// In unit tests without azd, it should fail with a connection error.
	err := cmd.Execute()
	if err == nil {
		t.Error("Listen command should error without azd server running")
	}
}
