package configurator

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/nhatthm/n26cli/internal/io"
)

// WithStdio configures stdio for prompt.
func WithStdio(stdio terminal.Stdio) Option {
	return func(c *PromptConfigurator) {
		c.defaultOptions = append(c.defaultOptions, survey.WithStdio(
			stdio.In, stdio.Out, stdio.Err,
		))
	}
}

// WithStdioProvider configures stdio for prompt.
func WithStdioProvider(p io.StdioProvider) Option {
	in, ok := p.InOrStdin().(terminal.FileReader)
	if !ok {
		return configureNothing
	}

	out, ok := p.OutOrStdout().(terminal.FileWriter)
	if !ok {
		return configureNothing
	}

	return func(c *PromptConfigurator) {
		WithStdio(terminal.Stdio{
			In:  in,
			Out: out,
			Err: p.ErrOrStderr(),
		})(c)
	}
}

func configureNothing(*PromptConfigurator) {}
