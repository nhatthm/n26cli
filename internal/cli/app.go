package cli

import (
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nhatthm/n26cli/internal/command"
	"github.com/nhatthm/n26cli/internal/command/transaction"
)

var rootCfg = defaultConfig()

// NewApp creates a new cli application using cobra.Command.
func NewApp() *cobra.Command {
	root := &cobra.Command{
		Use:   "n26",
		Short: "n26 command-line interface",
		Long:  "An awesome tool for managing your N26 account from the terminal",
	}

	root.PersistentFlags().StringVarP(&rootCfg.ConfigFile, "config", "c", rootCfg.ConfigFile, "configuration file")
	root.PersistentFlags().BoolVarP(&rootCfg.Verbose, "verbose", "v", rootCfg.Verbose, "verbose output")
	root.PersistentFlags().BoolVarP(&rootCfg.Debug, "debug", "d", rootCfg.Debug, "debug output")

	root.AddCommand(
		newAPICommand(transaction.NewTransactions),
		command.NewConfig(),
		command.NewVersion(),
	)

	return root
}

func defaultConfig() globalConfig {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return globalConfig{
		ConfigFile: filepath.Join(usr.HomeDir, ".n26/config.toml"),
	}
}
