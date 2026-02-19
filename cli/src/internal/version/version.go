// Package version provides version information for the azd exec extension.
// Version information is set at build time via ldflags.
package version

import coreversion "github.com/jongio/azd-core/version"

// Version is the semantic version of the extension, set at build time via ldflags.
//
//	go build -ldflags "-X github.com/jongio/azd-exec/cli/src/internal/version.Version=1.0.0 -X github.com/jongio/azd-exec/cli/src/internal/version.BuildDate=2025-01-09T12:00:00Z -X github.com/jongio/azd-exec/cli/src/internal/version.GitCommit=abc123"
var Version = "0.0.0-dev"

// BuildDate is the build timestamp, set at build time via ldflags.
var BuildDate = "unknown"

// GitCommit is the git commit hash, set at build time via ldflags.
var GitCommit = "unknown"

// Info provides the shared version information for this extension.
var Info = coreversion.New("jongio.azd.exec", "azd exec")

func init() {
	Info.Version = Version
	Info.BuildDate = BuildDate
	Info.GitCommit = GitCommit
}
