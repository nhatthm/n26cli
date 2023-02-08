package configurator

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/google/uuid"

	"github.com/nhatthm/n26cli/internal/service"
)

// Configure runs the configuration prompt.
func (c *PromptConfigurator) Configure() (err error) {
	cfg, err := c.SafeRead()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			_, _ = fmt.Fprintf(c.stdout, "\nsaved\n")

			return
		}

		if errors.Is(err, terminal.InterruptErr) {
			err = nil

			_, _ = fmt.Fprintf(c.stdout, "\nno change\n")

			return
		}
	}()

	oldCfg := cfg
	oldCfg.N26.Device = cfg.N26.Device

	// We do not do this in Read() on purpose, we only load sensitive information on demand.
	if oldCfg.N26.Device != emptyUUID && oldCfg.N26.CredentialsProvider == service.CredentialsProviderKeychain {
		cred := c.getCredentialsProvider(oldCfg.N26.Device)

		// Preload to change it later.
		cfg.N26.Username = cred.Username()
		cfg.N26.Password = cred.Password()

		oldCfg.N26.Username = cfg.N26.Username
		oldCfg.N26.Password = cfg.N26.Password
	}

	if err = c.configureDeviceID(&oldCfg.N26.Device, &cfg.N26); err != nil {
		return err
	}

	if err = c.configureCredentials(&cfg.N26); err != nil {
		return err
	}

	if err = c.Clean(oldCfg, cfg); err != nil {
		return err
	}

	return c.Write(cfg)
}

func (c *PromptConfigurator) configureDeviceID(oldDeviceID *uuid.UUID, current *service.N26Config) error {
	if current.Device == emptyUUID {
		current.Device = uuid.New()
		*oldDeviceID = current.Device

		return nil
	}

	// Ask to change device id.
	return c.askGenerateDeviceID(&current.Device)
}

func (c *PromptConfigurator) askGenerateDeviceID(current *uuid.UUID) error {
	var regenerate bool

	err := survey.AskOne(&survey.Confirm{Message: "Do you want to generate a new device id?"}, &regenerate, c.defaultOptions...)

	if regenerate {
		*current = uuid.New()
	}

	return err
}

func (c *PromptConfigurator) configureCredentials(current *service.N26Config) error {
	if err := c.askCredentialsProvider(&current.CredentialsProvider); err != nil {
		return err
	}

	// Ask to use keychain.
	if current.CredentialsProvider == service.CredentialsProviderKeychain {
		// Ask to set username and password.
		if err := c.askCredentials(&current.Username, &current.Password); err != nil {
			return err
		}
	}

	return nil
}

func (c *PromptConfigurator) askCredentialsProvider(current *service.CredentialsProviderType) error {
	use := *current == service.CredentialsProviderKeychain

	q := &survey.Confirm{
		Message: "Do you want to save your credentials to system keychain?",
		Default: use,
	}

	err := survey.AskOne(q, &use, c.defaultOptions...)

	if use {
		*current = service.CredentialsProviderKeychain
	} else {
		*current = service.CredentialsProviderNone
	}

	return err
}

func (c *PromptConfigurator) askCredentials(username, password *string) error {
	answer := map[string]interface{}{}

	questions := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Password{
				Message: "Enter username (input is hidden, leave it empty if no change) >",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Enter password (input is hidden, leave it empty if no change) >",
			},
		},
	}

	if err := survey.Ask(questions, &answer, c.defaultOptions...); err != nil {
		return err
	}

	if answer["username"] != "" {
		*username = answer["username"].(string) //nolint: errcheck
	}

	if answer["password"] != "" {
		*password = answer["password"].(string) //nolint: errcheck
	}

	return nil
}
