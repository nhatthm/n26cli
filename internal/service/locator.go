package service

import (
	"github.com/bool64/ctxd"
	clock "github.com/nhatthm/go-clock/service"
	"github.com/nhatthm/n26aas"
)

// Locator is a service locator.
type Locator struct {
	Config

	clock.ClockProvider
	ctxd.LoggerProvider

	n26aas.TransactionsFinderProvider
}
