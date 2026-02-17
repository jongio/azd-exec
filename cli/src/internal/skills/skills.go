package skills

import (
	"embed"

	"github.com/jongio/azd-core/copilotskills"
	"github.com/jongio/azd-exec/cli/src/internal/version"
)

//go:embed azd-exec/SKILL.md
var skillFS embed.FS

// InstallSkill installs the azd-exec skill to ~/.copilot/skills/azd-exec.
func InstallSkill() error {
	return copilotskills.Install("azd-exec", version.Version, skillFS, "azd-exec")
}
