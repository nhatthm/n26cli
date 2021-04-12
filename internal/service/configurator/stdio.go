package configurator

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	surveycobra "github.com/nhatthm/surveymock/cobra"
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
func WithStdioProvider(p surveycobra.StdioProvider) Option {
	return func(c *PromptConfigurator) {
		c.defaultOptions = append(c.defaultOptions, surveycobra.WithStdioProvider(p))
	}
}
