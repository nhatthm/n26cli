package io

import (
	"io"

	"github.com/gocarina/gocsv"
)

var _ DataWriter = (*CSVDataWriter)(nil)

// CSVDataWriter writes data in csv format.
type CSVDataWriter struct {
	writer io.Writer
}

// WriteData writes data in csv format.
func (w *CSVDataWriter) WriteData(v interface{}) error {
	return gocsv.Marshal(v, w.writer)
}

// DataWriter provides a data writer.
func (w *CSVDataWriter) DataWriter() DataWriter {
	return w
}

// CSVWriter initiates a new CSVDataWriter.
func CSVWriter(w io.Writer) *CSVDataWriter {
	return &CSVDataWriter{
		writer: w,
	}
}
