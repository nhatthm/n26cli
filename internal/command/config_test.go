package command

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigure_MissingConfigFlag(t *testing.T) {
	t.Parallel()

	err := configure(&cobra.Command{})
	expectedError := `flag accessed but not defined: config`

	assert.EqualError(t, err, expectedError)
}
