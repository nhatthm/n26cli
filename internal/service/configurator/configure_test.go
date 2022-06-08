package configurator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/google/uuid"
	keychainCredentialsMock "github.com/nhatthm/n26keychain/credentials/mock"
	keychainTokenMock "github.com/nhatthm/n26keychain/token/mock"
	"github.com/nhatthm/surveyexpect"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/n26cli/internal/service"
)

// nolint: gochecknoinits
func init() {
	surveyexpect.ReactionTime = 10 * time.Millisecond
}

func config(
	device uuid.UUID,
	credentialsProvider service.CredentialsProviderType,
	username, password string,
) *service.Config {
	return &service.Config{N26: service.N26Config{
		Username:            username,
		Password:            password,
		Device:              device,
		CredentialsProvider: credentialsProvider,
	}}
}

func TestPromptConfigurator_Configure(t *testing.T) {
	t.Parallel()

	oldDevice := uuid.New()

	assertSameDevice := func(t *testing.T, device uuid.UUID) {
		t.Helper()

		assert.Equal(t, oldDevice, device)
	}

	assertDifferentDevice := func(t *testing.T, device uuid.UUID) {
		t.Helper()

		assert.NotEqual(t, oldDevice, device)
	}

	testCases := []struct {
		scenario                string
		mockKeychainCredentials keychainCredentialsMock.Mocker
		mockKeychainStorage     keychainTokenMock.Mocker
		expectSurvey            surveyexpect.Expector
		oldCfg                  *service.Config
		expectConfigFile        bool
		assertDevice            func(t *testing.T, device uuid.UUID)
		expectedCredentials     string
		expectedError           string
	}{
		{
			scenario: "not configured before, no use keychain",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").No()
			}),
			expectConfigFile: true,
		},
		{
			scenario: "not configured before, use keychain",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "", "").Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >")
			}),
			expectConfigFile:    true,
			expectedCredentials: "keychain",
		},
		{
			scenario: "no change device id, no use keychain",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").No()
			}),
			oldCfg:           config(oldDevice, "", "", ""),
			expectConfigFile: true,
			assertDevice:     assertSameDevice,
		},
		{
			scenario: "change device id, no use keychain",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").Yes()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").No()
			}),
			oldCfg:           config(oldDevice, "", "", ""),
			expectConfigFile: true,
			assertDevice:     assertDifferentDevice,
		},
		{
			scenario: "no change device id, stop using keychain",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "", "").Return(nil)
				p.On("Username").Return("")
				p.On("Password").Return("")
				p.On("Delete").Return(nil)
			}),
			mockKeychainStorage: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), fmt.Sprintf(":%s", oldDevice.String())).
					Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").No()
			}),
			oldCfg:           config(oldDevice, "keychain", "", ""),
			expectConfigFile: true,
			assertDevice:     assertSameDevice,
		},
		{
			scenario: "change device id, keep using keychain",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "", "").Return(nil)
				p.On("Username").Return("")
				p.On("Password").Return("")
				p.On("Delete").Return(nil)
			}),
			mockKeychainStorage: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), fmt.Sprintf(":%s", oldDevice.String())).
					Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.WithTimeout(time.Hour)
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").Yes()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >")
			}),
			oldCfg:              config(oldDevice, "keychain", "", ""),
			expectConfigFile:    true,
			assertDevice:        assertDifferentDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "start using keychain",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").
					Answer("username")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >").
					Answer("password")
			}),
			oldCfg:              config(oldDevice, "", "", ""),
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "change username",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				p.On("Update", "foobar", "password").Return(nil)
			}),
			mockKeychainStorage: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), fmt.Sprintf("username:%s", oldDevice.String())).
					Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").
					Answer("foobar")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >")
			}),
			oldCfg:              config(oldDevice, "keychain", "username", "password"),
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "change username and password",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				p.On("Update", "foobar", "123456").Return(nil)
			}),
			mockKeychainStorage: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), fmt.Sprintf("username:%s", oldDevice.String())).
					Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").
					Answer("foobar")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >").
					Answer("123456")
			}),
			oldCfg:              config(oldDevice, "keychain", "username", "password"),
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "empty username and password means no change",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				p.On("Update", "username", "password").Return(nil)
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >")
			}),
			oldCfg:              config(oldDevice, "keychain", "username", "password"),
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "interrupt at changing id",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").Interrupt()
			}),
			oldCfg: config(oldDevice, "", "", ""),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "",
		},
		{
			scenario: "not configured, interrupt at setting credentials provider",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Interrupt()
			}),
		},
		{
			scenario: "not configured, interrupt at setting username",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").Interrupt()
			}),
		},
		{
			scenario: "not configured, interrupt at setting password",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >").Interrupt()
			}),
		},
		{
			scenario: "configured, interrupt at setting credentials provider",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Interrupt()
			}),
			oldCfg: config(oldDevice, "", "", ""),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "",
		},
		{
			scenario: "configured with keychain, interrupt at setting credentials provider",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Interrupt()
			}),
			oldCfg: config(oldDevice, "keychain", "username", "password"),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "configured with keychain, interrupt at setting username",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").Interrupt()
			}),
			oldCfg: config(oldDevice, "keychain", "username", "password"),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "configured with keychain, interrupt at setting password",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				// No call to storage to ask to update because the configurator is interrupted.
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").Answer("foobar")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >").Interrupt()
			}),
			oldCfg: config(oldDevice, "keychain", "username", "password"),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
		},
		{
			scenario: "not configured, error while persisting configuration",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(errors.New("save keychain error"))
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (y/N)").Yes()
				s.ExpectPassword("Enter username (input is hidden, leave it empty if no change) >").Answer("username")
				s.ExpectPassword("Enter password (input is hidden, leave it empty if no change) >").Answer("password")
			}),
			expectedError: "save keychain error",
		},
		{
			scenario: "could not clean credentials",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				p.On("Delete").Return(errors.New("delete credentials error"))
				// No call to storage to ask to update because the configurator is interrupted.
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").No()
			}),
			oldCfg: config(oldDevice, "keychain", "username", "password"),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
			expectedError:       "delete credentials error",
		},
		{
			scenario: "could not clean token",
			mockKeychainCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Update", "username", "password").Return(nil)
				p.On("Username").Return("username")
				p.On("Password").Return("password")
				p.On("Delete").Return(nil)
				// No call to storage to ask to update because the configurator is interrupted.
			}),
			mockKeychainStorage: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), fmt.Sprintf("username:%s", oldDevice.String())).
					Return(errors.New("delete token error"))
			}),
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to generate a new device id? (y/N)").No()
				s.ExpectConfirm("Do you want to save your credentials to system keychain? (Y/n)").No()
			}),
			oldCfg: config(oldDevice, "keychain", "username", "password"),
			// old configuration should be the same.
			expectConfigFile:    true,
			assertDevice:        assertSameDevice,
			expectedCredentials: "keychain",
			expectedError:       "delete token error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.mockKeychainCredentials == nil {
				tc.mockKeychainCredentials = keychainCredentialsMock.NoMock
			}

			if tc.mockKeychainStorage == nil {
				tc.mockKeychainStorage = keychainTokenMock.NoMock
			}

			credentialsProvider := tc.mockKeychainCredentials(t)
			tokenStorage := tc.mockKeychainStorage(t)

			assertConfigFile := func(t *testing.T, configFile string) {
				t.Helper()

				if !tc.expectConfigFile {
					_, err := os.Stat(configFile)
					assert.True(t, os.IsNotExist(err))

					return
				}

				v := viper.New()

				v.SetConfigFile(configFile)
				err := v.ReadInConfig()
				assert.NoError(t, err)

				// Check device id.
				deviceRaw := v.GetString("n26.device")
				assert.NotEmpty(t, deviceRaw)

				device, err := uuid.Parse(deviceRaw)
				assert.NotEmpty(t, device)
				assert.NoError(t, err)

				if tc.assertDevice != nil {
					tc.assertDevice(t, device)
				}

				// Check credentials provider.
				provider := v.GetString("n26.credentials")
				assert.Equal(t, tc.expectedCredentials, provider)
			}

			tc.expectSurvey(t).Start(func(stdio terminal.Stdio) {
				configFile := filepath.Join(t.TempDir(), "config.toml")
				c := New(configFile,
					WithStdio(stdio),
					withKeychainCredentials(credentialsProvider),
					withKeychainToken(tokenStorage),
				)

				if tc.oldCfg != nil {
					err := c.Write(*tc.oldCfg)
					require.NoError(t, err)
				}

				err := c.Configure()

				if tc.expectedError == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedError)
				}

				assertConfigFile(t, configFile)
			})
		})
	}
}

func TestPromptConfigurator_ConfigureErrorReadingConfigFile(t *testing.T) {
	t.Parallel()

	configFile := filepath.Join(t.TempDir(), "config.unknown")
	_, err := os.Create(configFile)
	require.NoError(t, err)

	c := New(configFile)
	err = c.Configure()

	expectedError := `Unsupported Config Type "unknown"`
	assert.EqualError(t, err, expectedError)
}
