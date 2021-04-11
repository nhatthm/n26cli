package app

import (
	"context"
	"errors"
	"time"

	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
	clock "github.com/nhatthm/go-clock"
	"github.com/nhatthm/n26aas"
	"github.com/nhatthm/n26api"
	keychain "github.com/nhatthm/n26keychain/credentials"
	"github.com/nhatthm/n26keychain/token"
	prompt "github.com/nhatthm/n26prompt/credentials"

	"github.com/nhatthm/n26cli/internal/fmt"
	"github.com/nhatthm/n26cli/internal/service"
)

// ErrUnsupportedCredentialsProvider indicates that the provided credentials provider is not supported.
var ErrUnsupportedCredentialsProvider = errors.New("unsupported credentials provider")

// MakeServiceLocator creates application service locator.
func MakeServiceLocator(l *service.Locator) error {
	initLogger(l)
	initFormatter(l)

	l.ClockProvider = clock.New()

	client, err := initN26Client(l.Config.N26, l.Clock(), l.CtxdLogger())
	if err != nil {
		return err
	}

	l.TransactionsFinderProvider = client

	return nil
}

func initLogger(l *service.Locator) {
	l.LoggerProvider = zapctxd.New(l.Config.Log)
}

func initFormatter(l *service.Locator) {
	switch l.Config.OutputFormat {
	case service.OutputFormatPrettyJSON,
		service.OutputFormatNone:
		w := fmt.JSONWriter(l.OutOrStdout())
		w.SetIndent("", "    ")

		l.DataWriterProvider = w

	case service.OutputFormatJSON:
		l.DataWriterProvider = fmt.JSONWriter(l.OutOrStdout())

	case service.OutputFormatCSV:
		l.DataWriterProvider = fmt.CSVWriter(l.OutOrStdout())

	default:
		panic("unknown output format")
	}
}

func initN26Client(cfg service.N26Config, clock clock.Clock, logger ctxd.Logger) (*n26aas.Service, error) {
	credOption, err := getCredentialsProviderOption(cfg, logger)
	if err != nil {
		return nil, err
	}

	c := n26aas.New(cfg.Device,
		n26api.WithCredentials(cfg.Username, cfg.Password),
		credOption,
		prompt.WithCredentialsAtLast(prompt.WithLogger(logger)),
		token.WithTokenStorage(),
		n26api.WithClock(clock),
		n26api.WithMFATimeout(2*time.Minute),
	)

	return c, nil
}

func getCredentialsProviderOption(cfg service.N26Config, logger ctxd.Logger) (n26api.Option, error) {
	if cfg.CredentialsProvider == "" {
		return noN26ClientOption, nil
	}

	switch cfg.CredentialsProvider {
	case service.CredentialsProviderKeychain:
		return func(c *n26api.Client) {
			keychain.WithCredentialsProvider(keychain.WithLogger(logger))(c)
		}, nil

	case service.CredentialsProviderNone:
	default:
		return nil,
			ctxd.WrapError(context.Background(), ErrUnsupportedCredentialsProvider,
				"could not build credentials provider option",
				"provider", cfg.CredentialsProvider,
			)
	}

	return noN26ClientOption, nil
}

func noN26ClientOption(_ *n26api.Client) {
	// Do nothing.
}
