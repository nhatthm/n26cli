package fmt

import (
	"encoding/json"
	"io"
)

var _ DataWriter = (*JSONDataWriter)(nil)

// JSONDataWriter writes data as json.
type JSONDataWriter struct {
	*json.Encoder
}

// WriteData writes data as json.
func (w *JSONDataWriter) WriteData(v interface{}) error {
	return w.Encode(v)
}

// DataWriter provides a data writer.
func (w *JSONDataWriter) DataWriter() DataWriter {
	return w
}

// JSONWriter initiates a new JSONDataWriter.
func JSONWriter(w io.Writer) *JSONDataWriter {
	enc := json.NewEncoder(w)

	return &JSONDataWriter{
		Encoder: enc,
	}
}
