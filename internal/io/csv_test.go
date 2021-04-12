package io

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVDataWriter_WriteData_Error(t *testing.T) {
	t.Parallel()

	data := []testData{
		{
			Name:  "foo",
			Value: "bar",
		},
	}

	w := CSVWriter(&errorWriter{})
	err := w.WriteData(data)

	assert.EqualError(t, err, "write error")
}

func TestCSVDataWriter_WriteData_Success(t *testing.T) {
	t.Parallel()

	var sb strings.Builder

	data := []testData{
		{
			Name:  "foo",
			Value: "bar",
		},
	}

	expected := "name,value\nfoo,bar\n"

	w := CSVWriter(&sb)
	err := w.WriteData(data)

	assert.Equal(t, expected, sb.String())
	assert.NoError(t, err)
}

func TestCSVDataWriter_DataWriter(t *testing.T) {
	t.Parallel()

	w := CSVWriter(nil)

	assert.Equal(t, w, w.DataWriter())
}
