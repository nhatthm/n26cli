package command

import (
	"github.com/nhatthm/n26cli/internal/service"
	"github.com/spf13/cobra"
)

// NewConfig creates a new `config` command.
func NewConfig(l *service.Locator) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "configure",
		Long:  "configure",

		RunE: func(cmd *cobra.Command, args []string) error {
			return l.Configurator().Configure()
		},
	}
}
