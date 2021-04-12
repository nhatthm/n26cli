package configurator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nhatthm/n26cli/internal/service"
)

// Write writes configuration to config file and external storage.
func (c *PromptConfigurator) Write(cfg service.Config) error {
	if err := c.writeToKeychain(cfg); err != nil {
		return err
	}

	return c.writeToFile(cfg)
}

func (c *PromptConfigurator) writeToKeychain(cfg service.Config) error {
	if cfg.N26.CredentialsProvider != service.CredentialsProviderKeychain {
		return nil
	}

	return c.getCredentialsProvider(cfg.N26.Device).
		Update(cfg.N26.Username, cfg.N26.Password)
}

func (c *PromptConfigurator) writeToFile(cfg service.Config) error {
	fileInfo, err := c.fs.Stat(c.configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := c.fs.MkdirAll(filepath.Dir(c.configFile), os.ModePerm); err != nil {
			return err
		}

		if _, err := c.fs.Create(c.configFile); err != nil {
			return err
		}
	}

	if fileInfo != nil && fileInfo.IsDir() {
		return ErrConfigIsDir
	}

	c.viper.SetConfigFile(c.configFile)

	c.viper.Set("n26.device", cfg.N26.Device.String())
	c.viper.Set("n26.credentials", string(cfg.N26.CredentialsProvider))

	return c.viper.WriteConfig()
}

// Clean cleans old obsolete from storage.
func (c *PromptConfigurator) Clean(oldCfg service.Config, newCfg service.Config) error {
	if oldCfg.N26.CredentialsProvider == service.CredentialsProviderKeychain {
		return c.cleanKeychain(oldCfg, newCfg)
	}

	return nil
}

func (c *PromptConfigurator) cleanKeychain(oldCfg service.Config, newCfg service.Config) error {
	var (
		credentialsChanged bool
		tokenChanged       bool
	)

	if oldCfg.N26.Device != newCfg.N26.Device || oldCfg.N26.CredentialsProvider != newCfg.N26.CredentialsProvider {
		credentialsChanged = true
		tokenChanged = true
	}

	if oldCfg.N26.Username != newCfg.N26.Username {
		tokenChanged = true
	}

	if credentialsChanged {
		if err := c.getCredentialsProvider(oldCfg.N26.Device).Delete(); err != nil {
			return err
		}
	}

	if tokenChanged {
		key := fmt.Sprintf("%s:%s", oldCfg.N26.Username, oldCfg.N26.Device.String())

		if err := c.getTokenStorage().Delete(context.Background(), key); err != nil {
			return err
		}
	}

	return nil
}
