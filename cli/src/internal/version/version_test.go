// Package version provides version information for the azd exec extension.
package version

import "testing"

func TestVersionConstants(t *testing.T) {
	if ExtensionID != "jongio.azd.exec" {
		t.Errorf("ExtensionID = %q, want %q", ExtensionID, "jongio.azd.exec")
	}

	if Name != "azd exec" {
		t.Errorf("Name = %q, want %q", Name, "azd exec")
	}

	// Version, BuildDate, and GitCommit are set at build time,
	// so we just check they're not empty after a proper build
	if Version == "" {
		t.Error("Version should not be empty")
	}
}
