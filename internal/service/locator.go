package service

import (
	"github.com/bool64/ctxd"
	clock "github.com/nhatthm/go-clock/service"
	"github.com/nhatthm/n26aas"
	"github.com/spf13/afero"

	"github.com/nhatthm/n26cli/internal/io"
)

// Locator is a service locator.
type Locator struct {
	Config

	afero.Fs

	clock.ClockProvider
	io.DataWriterProvider
	io.StdioProvider
	ctxd.LoggerProvider

	ConfiguratorProvider
	n26aas.TransactionsFinderProvider
}
