// Package version holds build information injected at link time.
//
// Release builds set these variables with -ldflags, for example:
//
//	go build -ldflags "-X gitlab.com/corvenc/corven/internal/version.Version=v0.1.0"
package version

import (
	"fmt"
	"runtime"
)

// Build information. Development builds keep the defaults.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String returns a one-line description of the build, for --version output.
func String(binary string) string {
	return fmt.Sprintf("%s %s (commit %s, built %s, %s %s/%s)",
		binary, Version, Commit, Date, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
