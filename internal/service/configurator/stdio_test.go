package configurator

import (
	"bytes"
	"io"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
)

type stdio struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (s *stdio) OutOrStdout() io.Writer {
	return s.out
}

func (s *stdio) ErrOrStderr() io.Writer {
	return s.err
}

func (s *stdio) InOrStdin() io.Reader {
	return s.in
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Fd() uintptr {
	return 0
}

func TestWithStdioProvider(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		in       io.Reader
		out      io.Writer
		err      io.Writer
		expected bool
	}{
		{
			scenario: "in is not a terminal.FileReader",
		},
		{
			scenario: "out is not a terminal.FileWriter",
			in:       &buffer{},
		},
		{
			scenario: "success",
			in:       &buffer{},
			out:      &buffer{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := &PromptConfigurator{}

			WithStdioProvider(&stdio{
				in:  tc.in,
				out: tc.out,
				err: tc.err,
			})(p)

			result := &survey.AskOptions{}

			for _, o := range p.defaultOptions {
				_ = o(result)
			}

			if tc.expected {
				assert.Equal(t, tc.in, result.Stdio.In)
				assert.Equal(t, tc.out, result.Stdio.Out)
				assert.Equal(t, tc.err, result.Stdio.Err)
			} else {
				assert.Nil(t, result.Stdio.In)
				assert.Nil(t, result.Stdio.Out)
				assert.Nil(t, result.Stdio.Err)
			}
		})
	}
}
