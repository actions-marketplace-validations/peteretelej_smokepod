package smokepod

import "fmt"

// Version information set via ldflags at build time.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// VersionString returns a formatted version string.
func VersionString() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
