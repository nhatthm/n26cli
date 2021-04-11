package configurator

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	keychainCredentialsMock "github.com/nhatthm/n26keychain/credentials/mock"
	keychainTokenMock "github.com/nhatthm/n26keychain/token/mock"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26cli/internal/service"
)

func TestPromptConfigurator_Write(t *testing.T) {}

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
