package command

import (
	"github.com/spf13/cobra"

	"github.com/nhatthm/n26cli/internal/service/configurator"
)

// NewConfig creates a new `config` command.
func NewConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "configure",
		Long:  "configure",

		RunE: func(cmd *cobra.Command, args []string) error {
			return configure(cmd)
		},
	}
}

func configure(cmd *cobra.Command) error {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	return configurator.New(configFile, configurator.WithStdioProvider(cmd)).Configure()
}
