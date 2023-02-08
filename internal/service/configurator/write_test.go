package configurator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	keychainCredentialsMock "github.com/nhatthm/n26keychain/credentials/mock"
	keychainTokenMock "github.com/nhatthm/n26keychain/token/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.nhat.io/aferomock"

	"github.com/nhatthm/n26cli/internal/service"
)

func TestPromptConfigurator_Write_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockFs        aferomock.FsMocker
		expectedError string
	}{
		{
			scenario: "stat error",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "path/config.toml").
					Return(nil, errors.New("stat error"))
			}),
			expectedError: "stat error",
		},
		{
			scenario: "mkdir error",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "path/config.toml").
					Return(nil, afero.ErrFileNotFound)

				fs.On("MkdirAll", "path", os.ModePerm).
					Return(errors.New("mkdir error"))
			}),
			expectedError: "mkdir error",
		},
		{
			scenario: "create file error",
			mockFs: aferomock.MockFs(func(fs *aferomock.Fs) {
				fs.On("Stat", "path/config.toml").
					Return(nil, afero.ErrFileNotFound)

				fs.On("MkdirAll", "path", os.ModePerm).
					Return(nil)

				fs.On("Create", "path/config.toml").
					Return(nil, errors.New("create file error"))
			}),
			expectedError: "create file error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := New("path/config.toml", WithFileSystem(tc.mockFs(t)))
			err := c.Write(service.Config{})

			assert.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestPromptConfigurator_Write_ConfigIsDir(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	err := c.Write(service.Config{})

	expectedError := "config file is a directory"
	assert.EqualError(t, err, expectedError)
}

func TestPromptConfigurator_Clean(t *testing.T) {
	t.Parallel()

	oldUsername := "username"
	oldID := uuid.New()
	oldKey := fmt.Sprintf("%s:%s", oldUsername, oldID.String())
	newID := uuid.New()

	testCases := []struct {
		scenario        string
		mockCredentials keychainCredentialsMock.Mocker
		mockToken       keychainTokenMock.Mocker
		oldConfig       service.Config
		newConfig       service.Config
		expectedError   string
	}{
		{
			scenario: "nothing to do if it is not keychain",
		},
		{
			scenario: "nothing to do if credentials is the same",
			oldConfig: service.Config{
				N26: service.N26Config{CredentialsProvider: service.CredentialsProviderKeychain},
			},
			newConfig: service.Config{
				N26: service.N26Config{CredentialsProvider: service.CredentialsProviderKeychain},
			},
		},
		{
			scenario: "device id is changed",
			mockCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Delete").Return(nil)
			}),
			mockToken: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), oldKey).Return(nil)
			}),
			oldConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					Username:            oldUsername,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			newConfig: service.Config{
				N26: service.N26Config{
					Device:              newID,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
		},
		{
			scenario: "provider is changed",
			mockCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Delete").Return(nil)
			}),
			mockToken: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), oldKey).Return(nil)
			}),
			oldConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					Username:            oldUsername,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			newConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					CredentialsProvider: service.CredentialsProviderNone,
				},
			},
		},
		{
			scenario: "provider is changed",
			mockCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Delete").Return(nil)
			}),
			mockToken: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), oldKey).Return(nil)
			}),
			oldConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					Username:            oldUsername,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			newConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					CredentialsProvider: service.CredentialsProviderNone,
				},
			},
		},
		{
			scenario: "username is changed",
			mockToken: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), oldKey).Return(nil)
			}),
			oldConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					Username:            oldUsername,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			newConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
		},
		{
			scenario: "cannot delete credentials",
			mockCredentials: keychainCredentialsMock.Mock(func(p *keychainCredentialsMock.KeychainCredentials) {
				p.On("Delete").Return(errors.New("delete error"))
			}),
			oldConfig: service.Config{
				N26: service.N26Config{CredentialsProvider: service.CredentialsProviderKeychain},
			},
			newConfig: service.Config{
				N26: service.N26Config{CredentialsProvider: service.CredentialsProviderNone},
			},
			expectedError: "delete error",
		},
		{
			scenario: "cannot delete token",
			mockToken: keychainTokenMock.Mock(func(s *keychainTokenMock.Storage) {
				s.On("Delete", context.Background(), oldKey).Return(errors.New("delete error"))
			}),
			oldConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					Username:            oldUsername,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			newConfig: service.Config{
				N26: service.N26Config{
					Device:              oldID,
					CredentialsProvider: service.CredentialsProviderKeychain,
				},
			},
			expectedError: "delete error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			if tc.mockCredentials == nil {
				tc.mockCredentials = keychainCredentialsMock.NoMock
			}

			if tc.mockToken == nil {
				tc.mockToken = keychainTokenMock.NoMock
			}

			c := New("",
				withKeychainCredentials(tc.mockCredentials(t)),
				withKeychainToken(tc.mockToken(t)),
			)

			err := c.Clean(tc.oldConfig, tc.newConfig)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
