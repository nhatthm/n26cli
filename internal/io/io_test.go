package io

import "errors"

type errorWriter struct{}

func (w *errorWriter) Write([]byte) (_ int, err error) {
	return 0, errors.New("write error")
}

type testData struct {
	Name  string `csv:"name"  json:"name"`
	Value string `csv:"value" json:"value"`
}
