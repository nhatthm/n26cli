// Package version provides build information.
package version

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

// Build information. Populated at build-time.
// nolint:gochecknoglobals
var (
	version      = "dev"
	revision     string
	branch       string
	buildUser    string
	buildDate    string
	dependencies []*debug.Module
)

// Information holds app version info.
type Information struct {
	Version      string
	Revision     string
	Branch       string
	BuildUser    string
	BuildDate    string
	GoVersion    string
	GoOS         string
	GoArch       string
	Dependencies []*debug.Module
}

// Info returns app version info.
func Info() Information {
	return Information{
		Version:      version,
		Revision:     revision,
		Branch:       branch,
		BuildUser:    buildUser,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		GoOS:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		Dependencies: dependencies,
	}
}

//nolint:gochecknoinits
func init() {
	if info, available := debug.ReadBuildInfo(); available {
		dependencies = info.Deps
	}
}

// WriteInformation writes the formatted information.
func WriteInformation(w io.Writer, info Information, showFull bool) {
	_, _ = fmt.Fprintf(w, "%s (rev: %s; %s; %s/%s)\n", info.Version, info.Revision, info.GoVersion, info.GoOS, info.GoArch)

	if !showFull {
		return
	}

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "build user: %s\n", info.BuildUser)
	_, _ = fmt.Fprintf(w, "build date: %s\n", info.BuildDate)
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "dependencies:")

	for _, dep := range info.Dependencies {
		_, _ = fmt.Fprintf(w, "  %s: %s\n", dep.Path, dep.Version)
	}
}
