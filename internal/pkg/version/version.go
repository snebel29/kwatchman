package version

import (
	"fmt"
	"runtime"
)

var (
	// Following global variables are initialized using
	// LD_FLAGS during build process (Makefile)
	Version   string
	BuildDate = "UNKNOWN"
	Commit    = "UNKNOWN"
)

// Release hold information about the release
type Release struct {
	GitCommit string
	BuildDate string
	Version   string
	GoVersion string
	Compiler  string
	Platform  string
}

func (r Release) String() string {
	return fmt.Sprintf("%s (%s/%s) %s %s",
		r.Version,
		runtime.GOOS, runtime.GOARCH, r.GitCommit, r.BuildDate)
}

// GetVersion returns the kwatchman version.
func GetVersion() Release {
	return Release{
		GitCommit: Commit,
		BuildDate: BuildDate,
		Version:   Version,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
