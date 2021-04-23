package app

import (
	"context"
	"errors"
	"time"

	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
	"github.com/nhatthm/go-clock"
	"github.com/nhatthm/n26aas"
	"github.com/nhatthm/n26api"
	keychain "github.com/nhatthm/n26keychain/credentials"
	"github.com/nhatthm/n26keychain/token"
	prompt "github.com/nhatthm/n26prompt/credentials"
	"github.com/spf13/afero"

	"github.com/nhatthm/n26cli/internal/io"
	"github.com/nhatthm/n26cli/internal/service"
)

// ErrUnsupportedCredentialsProvider indicates that the provided credentials provider is not supported.
var ErrUnsupportedCredentialsProvider = errors.New("unsupported credentials provider")

// NewServiceLocator initiates a new *service.Locator.
func NewServiceLocator() *service.Locator {
	l := &service.Locator{}

	l.Fs = afero.NewOsFs()

	l.ClockProvider = clock.New()
	l.StdioProvider = io.DefaultStdio()

	l.Config.N26.BaseURL = n26api.BaseURL
	l.Config.N26.MFAWait = 5 * time.Second
	l.Config.N26.MFATimeout = 2 * time.Minute

	return l
}

// MakeServiceLocator creates application service locator.
func MakeServiceLocator(l *service.Locator) error {
	initLogger(l)
	initFormatter(l)

	client, err := initN26Client(l, l.Config.N26)
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
		w := io.JSONWriter(l.OutOrStdout())
		w.SetIndent("", "    ")

		l.DataWriterProvider = w

	case service.OutputFormatJSON:
		l.DataWriterProvider = io.JSONWriter(l.OutOrStdout())

	case service.OutputFormatCSV:
		l.DataWriterProvider = io.CSVWriter(l.OutOrStdout())

	default:
		panic("unknown output format")
	}
}

func initN26Client(l *service.Locator, cfg service.N26Config) (*n26aas.Service, error) {
	credOption, err := getCredentialsProviderOption(cfg, l.CtxdLogger())
	if err != nil {
		return nil, err
	}

	c := n26aas.New(cfg.Device,
		n26api.WithBaseURL(cfg.BaseURL),
		n26api.WithCredentials(cfg.Username, cfg.Password),
		credOption,
		prompt.WithCredentialsAtLast(
			prompt.WithStdioProvider(l.StdioProvider),
			prompt.WithLogger(l.CtxdLogger()),
		),
		token.WithTokenStorage(),
		n26api.WithClock(l.Clock()),
		n26api.WithMFAWait(cfg.MFAWait),
		n26api.WithMFATimeout(cfg.MFATimeout),
	)

	return c, nil
}

func getCredentialsProviderOption(cfg service.N26Config, logger ctxd.Logger) (n26api.Option, error) {
	switch cfg.CredentialsProvider {
	case service.CredentialsProviderKeychain:
		return func(c *n26api.Client) {
			keychain.WithCredentialsProvider(keychain.WithLogger(logger))(c)
		}, nil

	case service.CredentialsProviderNone:
		return noN26ClientOption, nil

	default:
		return nil,
			ctxd.WrapError(context.Background(), ErrUnsupportedCredentialsProvider,
				"could not build credentials provider option",
				"provider", cfg.CredentialsProvider,
			)
	}
}

func noN26ClientOption(_ *n26api.Client) {
	// Do nothing.
}
