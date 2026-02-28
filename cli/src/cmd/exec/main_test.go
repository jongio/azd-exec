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

	cmd := newRootCmd()
	cmd.SetContext(context.Background())
	if setErr := cmd.PersistentFlags().Set("cwd", newWd); setErr != nil {
		t.Fatalf("setting cwd flag failed: %v", setErr)
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

	// Restore process state.
	_ = os.Chdir(oldWd)
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

func TestPersistentPreRunE_PropagatesDebugEnvVar(t *testing.T) {
	// Save and restore AZD_DEBUG
	origDebug := os.Getenv("AZD_DEBUG")
	defer func() {
		if origDebug != "" {
			_ = os.Setenv("AZD_DEBUG", origDebug)
		} else {
			_ = os.Unsetenv("AZD_DEBUG")
		}
	}()
	_ = os.Unsetenv("AZD_DEBUG")

	cmd := newRootCmd()
	cmd.SetContext(context.Background())
	if setErr := cmd.PersistentFlags().Set("debug", "true"); setErr != nil {
		t.Fatalf("setting debug flag failed: %v", setErr)
	}

	if runErr := cmd.PersistentPreRunE(cmd, []string{"echo"}); runErr != nil {
		t.Fatalf("PersistentPreRunE failed: %v", runErr)
	}

	if got := os.Getenv("AZD_DEBUG"); got != "true" {
		t.Errorf("AZD_DEBUG = %q, want %q", got, "true")
	}
}

func TestPersistentPreRunE_PropagatesNoPromptEnvVar(t *testing.T) {
	// Save and restore AZD_NO_PROMPT
	origNoPrompt := os.Getenv("AZD_NO_PROMPT")
	defer func() {
		if origNoPrompt != "" {
			_ = os.Setenv("AZD_NO_PROMPT", origNoPrompt)
		} else {
			_ = os.Unsetenv("AZD_NO_PROMPT")
		}
	}()
	_ = os.Unsetenv("AZD_NO_PROMPT")

	cmd := newRootCmd()
	cmd.SetContext(context.Background())
	if setErr := cmd.PersistentFlags().Set("no-prompt", "true"); setErr != nil {
		t.Fatalf("setting no-prompt flag failed: %v", setErr)
	}

	if runErr := cmd.PersistentPreRunE(cmd, []string{"echo"}); runErr != nil {
		t.Fatalf("PersistentPreRunE failed: %v", runErr)
	}

	if got := os.Getenv("AZD_NO_PROMPT"); got != "true" {
		t.Errorf("AZD_NO_PROMPT = %q, want %q", got, "true")
	}
}
