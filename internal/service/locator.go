package service

import (
	"github.com/bool64/ctxd"
	"github.com/nhatthm/n26aas"
	"github.com/spf13/afero"
	clock "go.nhat.io/clock/service"

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
