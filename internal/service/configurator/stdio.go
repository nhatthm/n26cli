package configurator

import (
	"github.com/AlecAivazis/survey/v2/terminal"

	"go.nhat.io/surveyexpect/options"
	"go.nhat.io/surveyexpect/options/cobra"
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
