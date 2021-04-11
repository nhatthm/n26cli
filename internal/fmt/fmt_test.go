package fmt

import "errors"

type errorWriter struct{}

func (w *errorWriter) Write([]byte) (_ int, err error) {
	return 0, errors.New("write error")
}

type testData struct {
	Name  string `json:"name" csv:"name"`
	Value string `json:"value" csv:"value"`
}
