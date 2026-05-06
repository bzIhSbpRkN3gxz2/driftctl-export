// Package version provides build-time version information for driftctl-export.
package version

import (
	"fmt"
	"runtime"
)

// Build-time variables, injected via -ldflags.
var (
	// Version is the semantic version string (e.g. "1.2.3").
	Version = "dev"
	// Commit is the short Git commit SHA.
	Commit = "none"
	// Date is the build date in RFC3339 format.
	Date = "unknown"
)

// Info holds all version metadata.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	GoVer   string `json:"go_version"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

// Get returns the current build Info.
func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
		GoVer:   runtime.Version(),
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
}

// String returns a human-readable one-line version string.
func (i Info) String() string {
	return fmt.Sprintf("%s (commit=%s, built=%s, %s/%s, %s)",
		i.Version, i.Commit, i.Date, i.OS, i.Arch, i.GoVer)
}
