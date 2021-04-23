package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nhatthm/n26cli/internal/app"
	"github.com/nhatthm/n26cli/internal/fmt"
	"github.com/nhatthm/n26cli/internal/service"
	"github.com/nhatthm/n26cli/internal/service/configurator"
)

// newAPICommand creates a new API command and decorates it with some global flags.
func newAPICommand(l *service.Locator, newCommand func(l *service.Locator) *cobra.Command) *cobra.Command {
	cliCfg := service.Config{
		OutputFormat: service.OutputFormatPrettyJSON,
	}

	cmd := newCommand(l)

	cmd.SetIn(l.InOrStdin())
	cmd.SetOut(l.OutOrStdout())
	cmd.SetErr(l.ErrOrStderr())

	cmd.Flags().StringVarP(&cliCfg.N26.Username, "username", "u", "", "n26 username")
	cmd.Flags().StringVarP(&cliCfg.N26.Password, "password", "p", "", "n26 password")
	cmd.Flags().StringVar(&cliCfg.OutputFormat, "format", "", "output format")

	run := runner(cmd)

	// If there is no runner, we do not have to setup the service locator.
	if run == nil {
		return cmd
	}

	cmd.RunE = nil
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := makeLocator(l, cliCfg); err != nil {
			handleErr(cmd, err)

			return
		}

		run(cmd, args)
	}

	return cmd
}

func runner(cmd *cobra.Command) func(cmd *cobra.Command, args []string) {
	if cmd.RunE == nil {
		return cmd.Run
	}

	fn := cmd.RunE

	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			handleErr(cmd, err)
		}
	}
}

func makeLocator(l *service.Locator, cliCfg service.Config) error {
	fileCfg, err := configurator.New(rootCfg.ConfigFile).SafeRead()
	if err != nil {
		return err
	}

	cliCfg.Log.Level = logLevel()
	cliCfg.Log.Output = l.ErrOrStderr()

	if err := mergeConfig(&l.Config, fileCfg, cliCfg); err != nil {
		return err
	}

	return app.MakeServiceLocator(l)
}

func logLevel() zapcore.Level {
	if rootCfg.Debug {
		return zap.DebugLevel
	}

	if rootCfg.Verbose {
		return zap.InfoLevel
	}

	return zap.WarnLevel
}

func handleErr(fmt fmt.Fmt, err error) {
	if rootCfg.Debug {
		panic(err)
	}

	fmt.PrintErrln(err)
}
