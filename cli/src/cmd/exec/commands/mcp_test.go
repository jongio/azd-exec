package commands

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/jongio/azd-core/azdextutil"
	"github.com/mark3labs/mcp-go/mcp"
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

	// Reset rate limiter so we don't hit limits
	saved := globalRateLimiter
	globalRateLimiter = azdextutil.NewRateLimiter(100, 100)
	defer func() { globalRateLimiter = saved }()

	result, err := handleGetEnvironment(context.Background(), mcp.CallToolRequest{})
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
// TestRateLimiting
// ---------------------------------------------------------------------------

func TestRateLimiting(t *testing.T) {
	// Create a limiter with only 2 burst tokens and 0 refill to test exhaustion.
	rl := azdextutil.NewRateLimiter(2, 0)

	if !rl.Allow() {
		t.Error("first Allow() should succeed")
	}
	if !rl.Allow() {
		t.Error("second Allow() should succeed")
	}
	// Third call must be rejected (burst exhausted, no refill)
	if rl.Allow() {
		t.Error("third Allow() should be rejected after burst exhaustion")
	}
}

func TestRateLimiting_HandlersReject(t *testing.T) {
	// Swap the global rate limiter with an exhausted one
	saved := globalRateLimiter
	globalRateLimiter = azdextutil.NewRateLimiter(0, 0) // zero tokens
	defer func() { globalRateLimiter = saved }()

	handlers := []struct {
		name string
		fn   func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
	}{
		{"handleGetEnvironment", handleGetEnvironment},
		{"handleListShells", handleListShells},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			result, err := h.fn(context.Background(), mcp.CallToolRequest{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			text, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("expected TextContent, got %T", result.Content[0])
			}
			if !strings.Contains(text.Text, "Rate limit exceeded") {
				t.Errorf("expected rate limit error, got: %s", text.Text)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestGetArgsMap / TestGetStringParam helpers
// ---------------------------------------------------------------------------

func TestGetArgsMap(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"shell": "bash",
		"count": 42,
	}

	m := getArgsMap(req)
	if m["shell"] != "bash" {
		t.Errorf("expected shell=bash, got %v", m["shell"])
	}

	// nil Arguments returns empty map
	req2 := mcp.CallToolRequest{}
	m2 := getArgsMap(req2)
	if len(m2) != 0 {
		t.Errorf("expected empty map for nil arguments, got %v", m2)
	}
}

func TestGetStringParam(t *testing.T) {
	args := map[string]interface{}{
		"name":  "test",
		"count": 42,
	}

	val, ok := getStringParam(args, "name")
	if !ok || val != "test" {
		t.Errorf("expected (test, true), got (%q, %v)", val, ok)
	}

	// non-string value
	_, ok = getStringParam(args, "count")
	if ok {
		t.Error("expected false for non-string value")
	}

	// missing key
	_, ok = getStringParam(args, "missing")
	if ok {
		t.Error("expected false for missing key")
	}
}
