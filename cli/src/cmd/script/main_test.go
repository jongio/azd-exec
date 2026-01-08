package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

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

	oldDebug := os.Getenv("AZD_SCRIPT_DEBUG")
	oldNoPrompt := os.Getenv("AZD_NO_PROMPT")
	oldEnvName := os.Getenv("AZURE_ENV_NAME")
	oldTraceFile := os.Getenv("AZD_TRACE_LOG_FILE")
	oldTraceURL := os.Getenv("AZD_TRACE_LOG_URL")

	cmd := newRootCmd()
	if err := cmd.PersistentFlags().Set("debug", "true"); err != nil {
		t.Fatalf("setting debug flag failed: %v", err)
	}
	if err := cmd.PersistentFlags().Set("no-prompt", "true"); err != nil {
		t.Fatalf("setting no-prompt flag failed: %v", err)
	}
	if err := cmd.PersistentFlags().Set("cwd", newWd); err != nil {
		t.Fatalf("setting cwd flag failed: %v", err)
	}
	if err := cmd.PersistentFlags().Set("environment", "test-env"); err != nil {
		t.Fatalf("setting environment flag failed: %v", err)
	}
	if err := cmd.PersistentFlags().Set("trace-log-file", "trace.log"); err != nil {
		t.Fatalf("setting trace-log-file flag failed: %v", err)
	}
	if err := cmd.PersistentFlags().Set("trace-log-url", "http://example.invalid"); err != nil {
		t.Fatalf("setting trace-log-url flag failed: %v", err)
	}
	if cmd.PersistentPreRunE == nil {
		t.Fatalf("expected PersistentPreRunE to be set")
	}

	if err := cmd.PersistentPreRunE(cmd, []string{"echo"}); err != nil {
		t.Fatalf("PersistentPreRunE failed: %v", err)
	}

	gotWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if gotWd != newWd {
		t.Fatalf("expected cwd %q, got %q", newWd, gotWd)
	}

	if os.Getenv("AZD_SCRIPT_DEBUG") != "true" {
		t.Fatalf("expected AZD_SCRIPT_DEBUG=true")
	}
	if os.Getenv("AZD_NO_PROMPT") != "true" {
		t.Fatalf("expected AZD_NO_PROMPT=true")
	}
	if os.Getenv("AZURE_ENV_NAME") != "test-env" {
		t.Fatalf("expected AZURE_ENV_NAME=test-env")
	}
	if os.Getenv("AZD_TRACE_LOG_FILE") != "trace.log" {
		t.Fatalf("expected AZD_TRACE_LOG_FILE=trace.log")
	}
	if os.Getenv("AZD_TRACE_LOG_URL") != "http://example.invalid" {
		t.Fatalf("expected AZD_TRACE_LOG_URL=http://example.invalid")
	}

	// Restore process state.
	_ = os.Chdir(oldWd)
	_ = os.Setenv("AZD_SCRIPT_DEBUG", oldDebug)
	_ = os.Setenv("AZD_NO_PROMPT", oldNoPrompt)
	_ = os.Setenv("AZURE_ENV_NAME", oldEnvName)
	_ = os.Setenv("AZD_TRACE_LOG_FILE", oldTraceFile)
	_ = os.Setenv("AZD_TRACE_LOG_URL", oldTraceURL)
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
	workingDir = ""
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
