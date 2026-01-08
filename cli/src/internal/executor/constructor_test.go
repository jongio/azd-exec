package executor

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "Default config",
			config: Config{},
		},
		{
			name: "With shell specified",
			config: Config{
				Shell: "bash",
			},
		},
		{
			name: "With interactive mode",
			config: Config{
				Interactive: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(tt.config)
			if exec == nil {
				t.Error("New() returned nil")
				return
			}
			if exec.config.Shell != tt.config.Shell {
				t.Errorf("Shell = %v, want %v", exec.config.Shell, tt.config.Shell)
			}
			if exec.config.Interactive != tt.config.Interactive {
				t.Errorf("Interactive = %v, want %v", exec.config.Interactive, tt.config.Interactive)
			}
		})
	}
}
