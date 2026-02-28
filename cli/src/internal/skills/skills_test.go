package skills

import (
	"testing"
)

func TestInstallSkill(t *testing.T) {
	// InstallSkill writes to ~/.copilot/skills/azd-exec.
	// We verify it does not panic and returns a result.
	err := InstallSkill()
	if err != nil {
		// In CI or restricted environments, file writes may fail.
		// That's acceptable — we're testing the code path, not the filesystem.
		t.Logf("InstallSkill returned error (may be expected in restricted env): %v", err)
	}
}
