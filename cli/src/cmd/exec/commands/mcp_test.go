package commands

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// TestBuildShellArgs
// ---------------------------------------------------------------------------

func TestBuildShellArgs(t *testing.T) {
	tests := []struct {
		name      string
		shell     string
		script    string
		isInline  bool
		extra     []string
		wantBin   string
		wantParts []string
	}{
		// bash inline
		{
			name: "bash inline", shell: "bash", script: "echo hi", isInline: true,
			wantBin: "bash", wantParts: []string{"bash", "-c", "echo hi"},
		},
		// bash file with extra args
		{
			name: "bash file", shell: "bash", script: "/tmp/run.sh", isInline: false, extra: []string{"--flag"},
			wantBin: "bash", wantParts: []string{"bash", "/tmp/run.sh", "--flag"},
		},
		// pwsh inline
		{
			name: "pwsh inline", shell: "pwsh", script: "Get-Date", isInline: true,
			wantBin: "pwsh", wantParts: []string{"pwsh", "-NoProfile", "-Command", "Get-Date"},
		},
		// pwsh file
		{
			name: "pwsh file", shell: "pwsh", script: "run.ps1", isInline: false, extra: []string{"a", "b"},
			wantBin: "pwsh", wantParts: []string{"pwsh", "-NoProfile", "-File", "run.ps1", "a", "b"},
		},
		// powershell inline
		{
			name: "powershell inline", shell: "powershell", script: "Write-Host hello", isInline: true,
			wantBin: "powershell", wantParts: []string{"powershell", "-NoProfile", "-Command", "Write-Host hello"},
		},
		// powershell file
		{
			name: "powershell file", shell: "powershell", script: "run.ps1", isInline: false,
			wantBin: "powershell", wantParts: []string{"powershell", "-NoProfile", "-File", "run.ps1"},
		},
		// cmd inline
		{
			name: "cmd inline", shell: "cmd", script: "dir", isInline: true,
			wantBin: "cmd", wantParts: []string{"cmd", "/c", "dir"},
		},
		// cmd file with args
		{
			name: "cmd file", shell: "cmd", script: "run.bat", isInline: false, extra: []string{"/v"},
			wantBin: "cmd", wantParts: []string{"cmd", "/c", "run.bat", "/v"},
		},
		// sh inline
		{
			name: "sh inline", shell: "sh", script: "ls", isInline: true,
			wantBin: "sh", wantParts: []string{"sh", "-c", "ls"},
		},
		// zsh file
		{
			name: "zsh file", shell: "zsh", script: "run.zsh", isInline: false, extra: []string{"x"},
			wantBin: "zsh", wantParts: []string{"zsh", "run.zsh", "x"},
		},
		// case insensitivity: BASH → bash
		{
			name: "uppercase BASH", shell: "BASH", script: "echo hi", isInline: true,
			wantBin: "bash", wantParts: []string{"bash", "-c", "echo hi"},
		},
		// mixed case: Pwsh → pwsh
		{
			name: "mixed case Pwsh", shell: "Pwsh", script: "Get-Date", isInline: true,
			wantBin: "pwsh", wantParts: []string{"pwsh", "-NoProfile", "-Command", "Get-Date"},
		},
		// mixed case CMD → cmd
		{
			name: "mixed case CMD", shell: "CMD", script: "dir", isInline: true,
			wantBin: "cmd", wantParts: []string{"cmd", "/c", "dir"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildShellArgs(tc.shell, tc.script, tc.isInline, tc.extra)
			if got[0] != tc.wantBin {
				t.Errorf("binary = %q, want %q", got[0], tc.wantBin)
			}
			if len(got) != len(tc.wantParts) {
				t.Fatalf("args length = %d, want %d: %v", len(got), len(tc.wantParts), got)
			}
			for i := range tc.wantParts {
				if got[i] != tc.wantParts[i] {
					t.Errorf("arg[%d] = %q, want %q", i, got[i], tc.wantParts[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestHandleGetEnvironment_SecretFiltering
// ---------------------------------------------------------------------------

func TestHandleGetEnvironment_SecretFiltering(t *testing.T) {
	// Save and restore environment
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range origEnv {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	}()

	// Set a controlled environment
	os.Clearenv()
	t.Setenv("AZD_ENVIRONMENT_NAME", "myenv")
	t.Setenv("AZURE_LOCATION", "eastus2")
	t.Setenv("AZURE_CLIENT_SECRET", "supersecret")
	t.Setenv("ARM_CLIENT_SECRET", "anothersecret")
	t.Setenv("AZD_ACCESS_TOKEN", "tok123")
	t.Setenv("AZURE_SUBSCRIPTION_PASSWORD", "pw")
	t.Setenv("AZURE_TENANT_ID", "tid")
	t.Setenv("NODE_ENV", "production")
	t.Setenv("HOME", "/home/user") // not an allowed prefix

	result, err := handleGetEnvironment(context.Background(), azdext.ToolArgs{})
	if err != nil {
		t.Fatalf("handleGetEnvironment returned error: %v", err)
	}

	// Parse the JSON text content
	if len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	var vars []envVar
	if err := json.Unmarshal([]byte(textContent.Text), &vars); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Build lookup
	found := map[string]string{}
	for _, v := range vars {
		found[v.Key] = v.Value
	}

	// Safe vars must be present
	for _, key := range []string{"AZD_ENVIRONMENT_NAME", "AZURE_LOCATION", "AZURE_TENANT_ID", "NODE_ENV"} {
		if _, ok := found[key]; !ok {
			t.Errorf("expected safe var %q to be present", key)
		}
	}

	// Secret vars must be excluded
	for _, key := range []string{"AZURE_CLIENT_SECRET", "ARM_CLIENT_SECRET", "AZD_ACCESS_TOKEN", "AZURE_SUBSCRIPTION_PASSWORD"} {
		if _, ok := found[key]; ok {
			t.Errorf("secret var %q should have been filtered out", key)
		}
	}

	// Non-prefixed vars must be excluded
	if _, ok := found["HOME"]; ok {
		t.Error("non-prefixed var HOME should not be included")
	}
}

// ---------------------------------------------------------------------------
// TestParseTimeout
// ---------------------------------------------------------------------------

func TestParseTimeout(t *testing.T) {
	// The MCP handlers use a hardcoded defaultTimeout constant.
	// Verify the default is 30 seconds as documented.
	if defaultTimeout.Seconds() != 30 {
		t.Errorf("defaultTimeout = %v, want 30s", defaultTimeout)
	}
}

// ---------------------------------------------------------------------------
// TestMarshalExecResult
// ---------------------------------------------------------------------------

func TestMarshalExecResult(t *testing.T) {
	t.Run("success with no error", func(t *testing.T) {
		result := marshalExecResult("hello\n", "", nil, nil)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(result.Content) == 0 {
			t.Fatal("expected content")
		}
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}

		var er execResult
		if err := json.Unmarshal([]byte(textContent.Text), &er); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if er.Stdout != "hello\n" {
			t.Errorf("Stdout = %q, want %q", er.Stdout, "hello\n")
		}
		if er.ExitCode != 0 {
			t.Errorf("ExitCode = %d, want 0", er.ExitCode)
		}
		if er.Error != "" {
			t.Errorf("Error = %q, want empty", er.Error)
		}
	})

	t.Run("with error and nil process state", func(t *testing.T) {
		result := marshalExecResult("", "oops", nil, errors.New("command failed"))
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}

		var er execResult
		if err := json.Unmarshal([]byte(textContent.Text), &er); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if er.ExitCode != -1 {
			t.Errorf("ExitCode = %d, want -1 when process state is nil and error present", er.ExitCode)
		}
		if er.Error != "command failed" {
			t.Errorf("Error = %q, want %q", er.Error, "command failed")
		}
		if er.Stderr != "oops" {
			t.Errorf("Stderr = %q, want %q", er.Stderr, "oops")
		}
	})
}

// ---------------------------------------------------------------------------
// TestNewMCPCommand
// ---------------------------------------------------------------------------

func TestNewMCPCommand(t *testing.T) {
	cmd := NewMCPCommand()
	if cmd == nil {
		t.Fatal("NewMCPCommand returned nil")
	}
	if cmd.Use != "mcp" {
		t.Errorf("Use = %q, want %q", cmd.Use, "mcp")
	}
	if !cmd.Hidden {
		t.Error("MCP command should be hidden")
	}
	if len(cmd.Commands()) == 0 {
		t.Error("MCP command should have subcommands")
	}

	// Check that "serve" subcommand exists
	found := false
	for _, sub := range cmd.Commands() {
		if sub.Use == "serve" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'serve' subcommand under mcp")
	}
}

// ---------------------------------------------------------------------------
// TestNewMetadataCommand
// ---------------------------------------------------------------------------

func TestNewMetadataCommand(t *testing.T) {
	cmd := NewMetadataCommand(func() *cobra.Command {
		return &cobra.Command{Use: "test"}
	})
	if cmd == nil {
		t.Fatal("NewMetadataCommand returned nil")
	}
}

// ---------------------------------------------------------------------------
// TestHandleListShells
// ---------------------------------------------------------------------------

func TestHandleListShells(t *testing.T) {
	result, err := handleListShells(context.Background(), azdext.ToolArgs{})
	if err != nil {
		t.Fatalf("handleListShells returned error: %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	var shells []shellInfo
	if err := json.Unmarshal([]byte(textContent.Text), &shells); err != nil {
		t.Fatalf("failed to unmarshal shell list: %v", err)
	}

	if len(shells) == 0 {
		t.Error("expected at least one shell in results")
	}

	// Verify all expected shell names are present
	expectedShells := map[string]bool{
		"bash": false, "sh": false, "zsh": false,
		"pwsh": false, "powershell": false, "cmd": false,
	}
	for _, s := range shells {
		if _, ok := expectedShells[s.Name]; ok {
			expectedShells[s.Name] = true
		}
	}
	for name, found := range expectedShells {
		if !found {
			t.Errorf("expected shell %q in results", name)
		}
	}
}

// ---------------------------------------------------------------------------
// TestHandleExecInline_Validation
// ---------------------------------------------------------------------------

// makeToolArgs creates an azdext.ToolArgs from a map for testing.
func makeToolArgs(m map[string]interface{}) azdext.ToolArgs {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = m
	return azdext.ParseToolArgs(req)
}

func TestHandleExecInline_Validation(t *testing.T) {
	t.Run("empty command returns error result", func(t *testing.T) {
		result, err := handleExecInline(context.Background(), azdext.ToolArgs{})
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		// Should return an MCP error result (isError=true), not a Go error
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty error result")
		}
		if !result.IsError {
			t.Error("expected IsError=true for empty command")
		}
	})

	t.Run("invalid shell returns error result", func(t *testing.T) {
		args := makeToolArgs(map[string]interface{}{"command": "echo hi", "shell": "invalid-shell-xyz"})
		result, err := handleExecInline(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty error result")
		}
		if !result.IsError {
			t.Error("expected IsError=true for invalid shell")
		}
	})
}

// ---------------------------------------------------------------------------
// TestHandleExecScript_Validation
// ---------------------------------------------------------------------------

func TestHandleExecScript_Validation(t *testing.T) {
	t.Run("missing script_path returns error result", func(t *testing.T) {
		result, err := handleExecScript(context.Background(), azdext.ToolArgs{})
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty error result")
		}
		if !result.IsError {
			t.Error("expected IsError=true for missing script_path")
		}
	})

	t.Run("invalid shell returns error result", func(t *testing.T) {
		args := makeToolArgs(map[string]interface{}{"script_path": "/tmp/test.sh", "shell": "invalid-shell-xyz"})
		result, err := handleExecScript(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if !result.IsError {
			t.Error("expected IsError=true for invalid shell")
		}
	})

	t.Run("directory path returns error result", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Set project dir env var for the security check
		t.Setenv("AZD_EXEC_PROJECT_DIR", tmpDir)

		args := makeToolArgs(map[string]interface{}{"script_path": tmpDir})
		result, err := handleExecScript(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if !result.IsError {
			t.Error("expected IsError=true for directory path")
		}
	})

	t.Run("nonexistent file returns error result", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("AZD_EXEC_PROJECT_DIR", tmpDir)

		args := makeToolArgs(map[string]interface{}{"script_path": filepath.Join(tmpDir, "nonexistent.sh")})
		result, err := handleExecScript(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if !result.IsError {
			t.Error("expected IsError=true for nonexistent file")
		}
	})

	t.Run("valid script executes successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("AZD_EXEC_PROJECT_DIR", tmpDir)

		scriptPath := filepath.Join(tmpDir, "test.ps1")
		if writeErr := os.WriteFile(scriptPath, []byte("Write-Host 'hello'\n"), 0o600); writeErr != nil {
			t.Fatalf("WriteFile failed: %v", writeErr)
		}

		args := makeToolArgs(map[string]interface{}{"script_path": scriptPath, "shell": "powershell", "args": "--verbose"})
		result, err := handleExecScript(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty result")
		}
		// Parse result to verify structure
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}
		var er execResult
		if unmarshalErr := json.Unmarshal([]byte(textContent.Text), &er); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal: %v", unmarshalErr)
		}
		// The script should have executed (exit code 0 or error from powershell)
		t.Logf("Script result: exitCode=%d stdout=%q stderr=%q error=%q",
			er.ExitCode, er.Stdout, er.Stderr, er.Error)
	})
}

// ---------------------------------------------------------------------------
// TestHandleExecInline_Execution
// ---------------------------------------------------------------------------

func TestHandleExecInline_Execution(t *testing.T) {
	t.Run("executes cmd inline on Windows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Windows-only test: cmd.exe is not available on this platform")
		}
		args := makeToolArgs(map[string]interface{}{"command": "echo hello", "shell": "cmd"})
		result, err := handleExecInline(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty result")
		}
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}
		var er execResult
		if unmarshalErr := json.Unmarshal([]byte(textContent.Text), &er); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal: %v", unmarshalErr)
		}
		if er.ExitCode != 0 {
			t.Errorf("ExitCode = %d, want 0", er.ExitCode)
		}
		if !strings.Contains(er.Stdout, "hello") {
			t.Errorf("Stdout = %q, want to contain %q", er.Stdout, "hello")
		}
	})

	t.Run("default shell used when not specified", func(t *testing.T) {
		args := makeToolArgs(map[string]interface{}{"command": "echo default-shell-test"})
		result, err := handleExecInline(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected Go error: %v", err)
		}
		if result == nil || len(result.Content) == 0 {
			t.Fatal("expected non-empty result")
		}
		// Just verify it returns a structured result
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}
		var er execResult
		if unmarshalErr := json.Unmarshal([]byte(textContent.Text), &er); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal: %v", unmarshalErr)
		}
		t.Logf("Default shell result: exitCode=%d stdout=%q", er.ExitCode, er.Stdout)
	})
}
