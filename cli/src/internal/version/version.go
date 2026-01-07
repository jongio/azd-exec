package version

// Version is the current version of the azd exec extension.
// It is intended to be set at build time via ldflags.
var Version = "0.0.0-dev"

// BuildDate is the UTC timestamp of the build.
// It is intended to be set at build time via ldflags.
var BuildDate = "unknown"

// GitCommit is the git SHA used for the build.
// It is intended to be set at build time via ldflags.
var GitCommit = "unknown"

// ExtensionID is the unique identifier for this extension.
const ExtensionID = "jongio.azd.exec"

// Name is the human-readable name of the extension.
const Name = "azd exec"
