package configurator

import (
	"errors"
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	keychainCredentials "github.com/nhatthm/n26keychain/credentials"
	keychainToken "github.com/nhatthm/n26keychain/token"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/nhatthm/n26cli/internal/service"
)

// ErrConfigIsDir indicates that the config file is a directory.
var ErrConfigIsDir = errors.New("config file is a directory")

var emptyUUID uuid.UUID

// Option configures Configurator.
type Option func(c *PromptConfigurator)

// Configurator manages application configuration.
type Configurator interface {
	Configure() error
	SafeRead() (service.Config, error)
	Read() (service.Config, error)
	Write(cfg service.Config) error
}

// PromptConfigurator manages service.Config.
type PromptConfigurator struct {
	fs                  afero.Fs
	viper               *viper.Viper
	stdout              io.Writer
	keychainCredentials keychainCredentials.KeychainCredentials
	keychainToken       keychainToken.KeychainStorage

	configFile     string
	defaultOptions []survey.AskOpt
}

func (c *PromptConfigurator) getCredentialsProvider(deviceID uuid.UUID) keychainCredentials.KeychainCredentials {
	if c.keychainCredentials == nil {
		return keychainCredentials.New(deviceID)
	}

	return c.keychainCredentials
}

func (c *PromptConfigurator) getTokenStorage() keychainToken.KeychainStorage {
	if c.keychainToken == nil {
		return keychainToken.NewStorage()
	}

	return c.keychainToken
}

// New creates a new Configurator.
func New(configFile string, options ...Option) *PromptConfigurator {
	c := &PromptConfigurator{
		fs:         afero.NewOsFs(),
		viper:      viper.New(),
		stdout:     os.Stdout,
		configFile: configFile,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

func withFileSystem(fs afero.Fs) Option {
	return func(c *PromptConfigurator) {
		c.fs = fs
	}
}

func withKeychainCredentials(credentials keychainCredentials.KeychainCredentials) Option {
	return func(c *PromptConfigurator) {
		c.keychainCredentials = credentials
	}
}

func withKeychainToken(token keychainToken.KeychainStorage) Option {
	return func(c *PromptConfigurator) {
		c.keychainToken = token
	}
}
