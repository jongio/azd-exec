package commands

import (
	"strings"
	"testing"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewRunCommand(t *testing.T) {
	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)

	if cmd == nil {
		t.Fatal("NewRunCommand returned nil")
	}

	if !contains(cmd.Use, "run") {
		t.Errorf("Command Use = %v, should contain 'run'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}

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
}
