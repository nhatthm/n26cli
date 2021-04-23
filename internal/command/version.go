package command

import (
	"github.com/spf13/cobra"

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
			version.WriteInformation(cmd.OutOrStdout(), version.Info(), showFull)
		},
	}

	cmd.Flags().BoolVarP(&showFull, "full", "f", false, "show full information")

	return cmd
}
