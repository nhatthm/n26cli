package command

import (
	"runtime"
	"sort"

	"github.com/spf13/cobra"

	"github.com/nhatthm/n26cli/internal/fmt"
	"github.com/nhatthm/n26cli/internal/version"
)

// NewVersion creates a new version command.
func NewVersion() *cobra.Command {
	var showFull bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show version information",
		Run: func(cmd *cobra.Command, _ []string) {
			showVersion(cmd, showFull)
		},
	}

	cmd.Flags().BoolVarP(&showFull, "full", "f", false, "show full information")

	return cmd
}

func showVersion(fmt fmt.Fmt, showFull bool) {
	info := version.Info()

	fmt.Printf("%s (rev: %s; %s; %s/%s)\n", info.Version, info.Revision, info.GoVersion, runtime.GOOS, runtime.GOARCH)

	if !showFull {
		return
	}

	fmt.Println()
	fmt.Printf("build user: %s\n", info.BuildUser)
	fmt.Printf("build date: %s\n", info.BuildDate)
	fmt.Println()
	fmt.Println("dependencies:")

	keys := make([]string, 0, len(info.Dependencies))

	for path := range info.Dependencies {
		keys = append(keys, path)
	}

	sort.Strings(keys)

	for _, path := range keys {
		fmt.Printf("  %s: %s\n", path, info.Dependencies[path])
	}
}
