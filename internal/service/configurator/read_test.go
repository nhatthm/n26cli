package configurator

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptConfigurator_SafeRead(t *testing.T) {}

func TestPromptConfigurator_Read_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	configFile := filepath.Join(t.TempDir(), "config.unknown")

	c := New(configFile)
	_, err := c.Read()

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestPromptConfigurator_Read_UnsupportedType(t *testing.T) {
	t.Parallel()

	configFile := filepath.Join(t.TempDir(), "config.unknown")
	_, err := os.Create(configFile)
	require.NoError(t, err)

	c := New(configFile)
	_, err = c.Read()

	expectedError := `Unsupported Config Type "unknown"`
	assert.EqualError(t, err, expectedError)
}

func TestPromptConfigurator_Read_Error(t *testing.T) {
	t.Parallel()

	configFile := filepath.Join(t.TempDir(), "config.toml")
	_, err := os.Create(configFile)
	require.NoError(t, err)

	v := viper.New()

	v.SetConfigFile(configFile)
	v.Set("n26.device", "wrong_uuid")
	v.Set("n26.credentials", "")

	err = v.WriteConfig()
	require.NoError(t, err)

	c := New(configFile)
	_, err = c.Read()

	expectedError := "could not read config: 1 error(s) decoding:\n\n* error decoding 'N26.Device': invalid UUID length: 10"
	assert.EqualError(t, err, expectedError)
}

func Test_stringToUUIDHook(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	uuidValue := reflect.ValueOf(uuid.UUID{})

	testCases := []struct {
		scenario       string
		source         reflect.Value
		destination    reflect.Value
		expectedResult interface{}
		expectedError  string
	}{
		{
			scenario:       "ignore if data is not string",
			source:         reflect.ValueOf(1),
			destination:    uuidValue,
			expectedResult: 1,
		},
		{
			scenario:       "ignore if destination is not a uuid",
			source:         reflect.ValueOf("foobar"),
			destination:    reflect.ValueOf(1),
			expectedResult: "foobar",
		},
		{
			scenario:      "parse error",
			source:        reflect.ValueOf("foobar"),
			destination:   uuidValue,
			expectedError: "invalid UUID length: 6",
		},
		{
			scenario:       "parse error",
			source:         reflect.ValueOf(id.String()),
			destination:    uuidValue,
			expectedResult: id,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := mapstructure.DecodeHookExec(stringToUUIDHook(), tc.source, tc.destination)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
