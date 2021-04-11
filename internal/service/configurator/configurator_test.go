package configurator

import (
	"testing"

	"github.com/google/uuid"
	keychainCredentials "github.com/nhatthm/n26keychain/credentials"
	keychainCredentialsMock "github.com/nhatthm/n26keychain/credentials/mock"
	keychainToken "github.com/nhatthm/n26keychain/token"
	keychainTokenMock "github.com/nhatthm/n26keychain/token/mock"
	"github.com/stretchr/testify/assert"
)

func TestPromptConfigurator_GetCredentialsProvider(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	testCases := []struct {
		scenario string
		provider keychainCredentials.KeychainCredentials
		expected keychainCredentials.KeychainCredentials
	}{
		{
			scenario: "not configured",
			expected: keychainCredentials.New(id),
		},
		{
			scenario: "configured",
			provider: keychainCredentialsMock.NoMock(t),
			expected: keychainCredentialsMock.NoMock(t),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := &PromptConfigurator{
				keychainCredentials: tc.provider,
			}

			assert.Equal(t, tc.expected, c.getCredentialsProvider(id))
		})
	}
}

func TestPromptConfigurator_GetTokenStorage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		storage  keychainToken.KeychainStorage
		expected keychainToken.KeychainStorage
	}{
		{
			scenario: "not configured",
			expected: keychainToken.NewStorage(),
		},
		{
			scenario: "configured",
			storage:  keychainTokenMock.NoMock(t),
			expected: keychainTokenMock.NoMock(t),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := &PromptConfigurator{
				keychainToken: tc.storage,
			}

			assert.Equal(t, tc.expected, c.getTokenStorage())
		})
	}
}
