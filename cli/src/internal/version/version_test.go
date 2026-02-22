package version

import "testing"

func TestVersionInfo(t *testing.T) {
	if Info.ExtensionID != "jongio.azd.exec" {
		t.Errorf("Info.ExtensionID = %q, want %q", Info.ExtensionID, "jongio.azd.exec")
	}

	if Info.Name != "azd exec" {
		t.Errorf("Info.Name = %q, want %q", Info.Name, "azd exec")
	}

	if Info.Version == "" {
		t.Error("Info.Version should not be empty")
	}
}
