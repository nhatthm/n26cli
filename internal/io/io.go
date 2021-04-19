package io

import (
	"io"
	"os"
)

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

type stdioProvider struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (p *stdioProvider) InOrStdin() io.Reader {
	return p.in
}

func (p *stdioProvider) OutOrStdout() io.Writer {
	return p.out
}

func (p *stdioProvider) ErrOrStderr() io.Writer {
	return p.err
}

// DefaultStdio returns default os.Stdin, os.Stdout and os.Stderr.
func DefaultStdio() StdioProvider {
	return &stdioProvider{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

// Stdio creates a new provider with then given stdio.
func Stdio(stdin io.Reader, stdout io.Writer, stderr io.Writer) StdioProvider {
	return &stdioProvider{
		in:  stdin,
		out: stdout,
		err: stderr,
	}
}
