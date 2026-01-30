package main

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jongio/azd-core/env"
	"github.com/jongio/azd-exec/cli/src/internal/executor"
)

type fakeExecutor struct {
	executePath   string
	inlineContent string
	args          []string
}

func (f *fakeExecutor) Execute(_ context.Context, scriptPath string) error {
	f.executePath = scriptPath
	return nil
}

func (f *fakeExecutor) ExecuteInline(_ context.Context, scriptContent string) error {
	f.inlineContent = scriptContent
	return nil
}

func TestPersistentPreRunE_SetsEnvAndCwd(t *testing.T) {
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}

	newWd := t.TempDir()

	oldDebug := os.Getenv("AZD_DEBUG")
	oldNoPrompt := os.Getenv("AZD_NO_PROMPT")
	oldTraceFile := os.Getenv("AZD_TRACE_LOG_FILE")
	oldTraceURL := os.Getenv("AZD_TRACE_LOG_URL")

	cmd := newRootCmd()
	cmd.SetContext(context.Background())
	if setErr := cmd.PersistentFlags().Set("debug", "true"); setErr != nil {
		t.Fatalf("setting debug flag failed: %v", setErr)
	}
	if setErr := cmd.PersistentFlags().Set("no-prompt", "true"); setErr != nil {
		t.Fatalf("setting no-prompt flag failed: %v", setErr)
	}
	if setErr := cmd.PersistentFlags().Set("cwd", newWd); setErr != nil {
		t.Fatalf("setting cwd flag failed: %v", setErr)
	}
	// Note: -e/--environment flag is not tested here because it now calls 'azd env get-values'
	// which requires an actual azd environment to exist. See TestLoadAzdEnvironment for unit tests.
	if setErr := cmd.PersistentFlags().Set("trace-log-file", "trace.log"); setErr != nil {
		t.Fatalf("setting trace-log-file flag failed: %v", setErr)
	}
	if setErr := cmd.PersistentFlags().Set("trace-log-url", "http://example.invalid"); setErr != nil {
		t.Fatalf("setting trace-log-url flag failed: %v", setErr)
	}
	if cmd.PersistentPreRunE == nil {
		t.Fatalf("expected PersistentPreRunE to be set")
	}

	if runErr := cmd.PersistentPreRunE(cmd, []string{"echo"}); runErr != nil {
		t.Fatalf("PersistentPreRunE failed: %v", runErr)
	}

	gotWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	// Normalize paths for comparison (macOS has /var -> /private/var symlink)
	gotWdNorm, err := filepath.EvalSymlinks(gotWd)
	if err != nil {
		t.Fatalf("EvalSymlinks failed for %q: %v", gotWd, err)
	}
	newWdNorm, err := filepath.EvalSymlinks(newWd)
	if err != nil {
		t.Fatalf("EvalSymlinks failed for %q: %v", newWd, err)
	}
	if gotWdNorm != newWdNorm {
		t.Fatalf("expected cwd %q, got %q", newWdNorm, gotWdNorm)
	}

	if os.Getenv("AZD_DEBUG") != "true" {
		t.Fatalf("expected AZD_DEBUG=true")
	}
	if os.Getenv("AZD_NO_PROMPT") != "true" {
		t.Fatalf("expected AZD_NO_PROMPT=true")
	}
	if os.Getenv("AZD_TRACE_LOG_FILE") != "trace.log" {
		t.Fatalf("expected AZD_TRACE_LOG_FILE=trace.log")
	}
	if os.Getenv("AZD_TRACE_LOG_URL") != "http://example.invalid" {
		t.Fatalf("expected AZD_TRACE_LOG_URL=http://example.invalid")
	}

	// Restore process state.
	_ = os.Chdir(oldWd)
	_ = os.Setenv("AZD_DEBUG", oldDebug)
	_ = os.Setenv("AZD_NO_PROMPT", oldNoPrompt)
	_ = os.Setenv("AZD_TRACE_LOG_FILE", oldTraceFile)
	_ = os.Setenv("AZD_TRACE_LOG_URL", oldTraceURL)
}

// TestLoadAzdEnvironment_Validation verifies that the azd-core env package
// is correctly integrated. The actual implementation is tested in azd-core/env.
func TestLoadAzdEnvironment_Validation(t *testing.T) {
	ctx := context.Background()

	// Just verify the function is accessible and basic validation works
	err := env.LoadAzdEnvironment(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty environment name")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestRunE_DispatchesFileOrInline(t *testing.T) {
	oldNew := newScriptExecutor
	defer func() { newScriptExecutor = oldNew }()

	fake := &fakeExecutor{}
	newScriptExecutor = func(cfg executor.Config) scriptExecutor {
		fake.args = append([]string{}, cfg.Args...)
		return fake
	}

	// Avoid changing env/cwd during Execute.
	debugMode = false
	noPrompt = false
	cwd = ""
	environment = ""
	traceLogFile = ""
	traceLogURL = ""
	shell = ""
	interactive = false

	t.Run("file path", func(t *testing.T) {
		fake.executePath = ""
		fake.inlineContent = ""

		tmp := t.TempDir()
		filePath := filepath.Join(tmp, "script.sh")
		if err := os.WriteFile(filePath, []byte("echo hi\n"), 0o600); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}

		cmd := newRootCmd()
		cmd.SetArgs([]string{filePath})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		abs, err := filepath.Abs(filePath)
		if err != nil {
			t.Fatalf("Abs failed: %v", err)
		}
		if fake.executePath != abs {
			t.Fatalf("expected Execute to be called with %q, got %q", abs, fake.executePath)
		}
		if fake.inlineContent != "" {
			t.Fatalf("expected ExecuteInline not to be called")
		}
	})

	t.Run("inline script", func(t *testing.T) {
		fake.executePath = ""
		fake.inlineContent = ""

		cmd := newRootCmd()
		cmd.SetArgs([]string{"echo hello"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if fake.inlineContent != "echo hello" {
			t.Fatalf("expected ExecuteInline to be called with %q, got %q", "echo hello", fake.inlineContent)
		}
		if fake.executePath != "" {
			t.Fatalf("expected Execute not to be called")
		}
	})
}

func TestRunE_AllowsPassthroughArgsWithoutDoubleDash(t *testing.T) {
	oldNew := newScriptExecutor
	defer func() { newScriptExecutor = oldNew }()

	fake := &fakeExecutor{}
	newScriptExecutor = func(cfg executor.Config) scriptExecutor {
		fake.args = append([]string{}, cfg.Args...)
		return fake
	}

	// Avoid changing env/cwd during Execute.
	debugMode = false
	noPrompt = false
	cwd = ""
	environment = ""
	traceLogFile = ""
	traceLogURL = ""
	shell = ""
	interactive = false

	cmd := newRootCmd()
	cmd.SetArgs([]string{"pnpm", "sync", "--skip-sync"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	expectedArgs := []string{"sync", "--skip-sync"}
	if !reflect.DeepEqual(fake.args, expectedArgs) {
		t.Fatalf("expected passthrough args %v, got %v", expectedArgs, fake.args)
	}

	if fake.inlineContent != "pnpm" {
		t.Fatalf("expected inline execution of 'pnpm', got %q", fake.inlineContent)
	}
}
