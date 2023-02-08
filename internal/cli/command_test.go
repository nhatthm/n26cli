package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/afero/mem"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.nhat.io/aferomock"

	"github.com/nhatthm/n26cli/internal/app"
	"github.com/nhatthm/n26cli/internal/service"
)

func TestNewAPICommand_NoRunner(t *testing.T) {
	constructor := func(l *service.Locator) *cobra.Command {
		return &cobra.Command{
			Use:  "test",
			Long: "test command",
		}
	}

	var sb strings.Builder

	cmd := newAPICommand(app.NewServiceLocator(), constructor)
	cmd.SetOut(&sb)

	err := cmd.Execute()

	expected := `test command`

	assert.NoError(t, err)
	assert.Equal(t, expected, strings.TrimRight(sb.String(), "\r\n"))
}

func TestRunner(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		cmd      *cobra.Command
	}{
		{
			scenario: "no RunE",
			cmd: &cobra.Command{
				Run: func(*cobra.Command, []string) {
					panic("this should happen")
				},
			},
		},
		{
			scenario: "RunE",
			cmd: &cobra.Command{
				RunE: func(*cobra.Command, []string) error {
					panic("this should happen")
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Panics(t, func() {
				runner(tc.cmd)(tc.cmd, nil)
			})
		})
	}
}

func TestMakeLocator_CouldNotReadConfigFile(t *testing.T) {
	t.Parallel()

	l := app.NewServiceLocator()
	l.Fs = aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", mock.Anything).
			Return(mem.GetFileInfo(nil), errors.New("stat error"))
	})(t)

	err := makeLocator(l, service.Config{})
	expectedError := `stat error`

	assert.EqualError(t, err, expectedError)
}
