package configurator

import (
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/nhatthm/surveyexpect/options"
	"github.com/nhatthm/surveyexpect/options/cobra"
)

// WithStdio configures stdio for prompt.
func WithStdio(stdio terminal.Stdio) Option {
	return func(c *PromptConfigurator) {
		c.stdout = stdio.Out
		c.defaultOptions = append(c.defaultOptions, options.WithStdio(stdio))
	}
}

// WithStdioProvider configures stdio for prompt.
func WithStdioProvider(p cobra.StdioProvider) Option {
	return func(c *PromptConfigurator) {
		c.stdout = p.OutOrStdout()
		c.defaultOptions = append(c.defaultOptions, cobra.WithStdioProvider(p))
	}
}
