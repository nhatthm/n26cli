package service

import (
	"github.com/bool64/ctxd"
	clock "github.com/nhatthm/go-clock/service"
	"github.com/nhatthm/n26aas"

	"github.com/nhatthm/n26cli/internal/io"
)

// Locator is a service locator.
type Locator struct {
	Config

	clock.ClockProvider
	io.DataWriterProvider
	io.StdioProvider
	ctxd.LoggerProvider

	n26aas.TransactionsFinderProvider
}
