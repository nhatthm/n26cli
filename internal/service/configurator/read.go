package configurator

import (
	"context"
	"os"
	"reflect"

	"github.com/bool64/ctxd"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/nhatthm/n26cli/internal/service"
)

// SafeRead is the same as Read but does not return error if the config file is missing.
func (c *PromptConfigurator) SafeRead() (service.Config, error) {
	cfg, err := c.Read()
	if err != nil && os.IsNotExist(err) {
		return service.Config{}, nil
	}

	return cfg, err
}

// Read reads configuration from config file.
func (c *PromptConfigurator) Read() (service.Config, error) {
	var cfg service.Config

	file, err := c.fs.Stat(c.configFile)
	if err != nil {
		return service.Config{}, err
	}

	if file.IsDir() {
		return service.Config{}, ErrConfigIsDir
	}

	c.viper.SetConfigFile(c.configFile)

	if err := c.viper.ReadInConfig(); err != nil {
		return service.Config{}, err
	}

	if err := c.viper.Unmarshal(&cfg, decodeOptions()...); err != nil {
		return service.Config{}, ctxd.WrapError(context.Background(), err, "could not read config")
	}

	return cfg, nil
}

func decodeOptions() []viper.DecoderConfigOption {
	return []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			stringToUUIDHook(),
		)),
		func(c *mapstructure.DecoderConfig) {
			c.TagName = "toml"
		},
	}
}

func stringToUUIDHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		// Parse uuid.
		id, err := uuid.Parse(data.(string))
		if err != nil {
			return nil, err
		}

		return id, nil
	}
}
