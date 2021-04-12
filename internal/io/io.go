package io

import "io"

// DataWriter is the interface that wraps the basic WriteData method.
type DataWriter interface {
	WriteData(v interface{}) error
}

// StdioProvider is a wrapper around *cobra.Command to provide stdin, stdout and stderr to survey.
type StdioProvider interface {
	OutOrStdout() io.Writer
	ErrOrStderr() io.Writer
	InOrStdin() io.Reader
}

// DataWriterProvider provides DataWriter.
type DataWriterProvider interface {
	DataWriter() DataWriter
}
