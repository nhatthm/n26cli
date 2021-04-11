// Package version provides build information.
package version

import (
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
	dependencies map[string]string
)

// Information holds app version info.
type Information struct {
	Version      string            `json:"version,omitempty"`
	Revision     string            `json:"revision,omitempty"`
	Branch       string            `json:"branch,omitempty"`
	BuildUser    string            `json:"build_user,omitempty"`
	BuildDate    string            `json:"build_date,omitempty"`
	GoVersion    string            `json:"go_version,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
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
		Dependencies: dependencies,
	}
}

//nolint:gochecknoinits
func init() {
	if info, available := debug.ReadBuildInfo(); available {
		dependencies = make(map[string]string, len(info.Deps))

		for _, dep := range info.Deps {
			dependencies[dep.Path] = dep.Version
		}
	}
}
