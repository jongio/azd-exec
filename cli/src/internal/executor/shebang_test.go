package executor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectShellWithShebang(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		filename string
		want     string
	}{
		{
			name:     "Bash shebang overrides .txt extension",
			content:  "#!/bin/bash\necho hello",
			filename: "script.txt",
			want:     "bash",
		},
		{
			name:     "Python shebang",
			content:  "#!/usr/bin/env python3\nprint('hello')",
			filename: "script",
			want:     "python3",
		},
		{
			name:     "Zsh shebang",
			content:  "#!/usr/bin/zsh\necho hello",
			filename: "script",
			want:     "zsh",
		},
	}

	exec := New(Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(scriptPath, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}

			got := exec.detectShell(scriptPath)
			if got != tt.want {
				t.Errorf("detectShell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadShebangFileNotFound(t *testing.T) {
	exec := New(Config{})
	got := exec.readShebang("nonexistent-file.sh")

	if got != "" {
		t.Errorf("readShebang(nonexistent) = %v, want empty string", got)
	}
}
