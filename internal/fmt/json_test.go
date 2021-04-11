package fmt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONDataWriter_WriteData(t *testing.T) {
	t.Parallel()

	data := []testData{
		{
			Name:  "foo",
			Value: "bar",
		},
	}

	t.Run("compact", func(t *testing.T) {
		t.Parallel()

		var sb strings.Builder
		w := JSONWriter(&sb)

		err := w.WriteData(data)

		expected := "[{\"name\":\"foo\",\"value\":\"bar\"}]\n"

		assert.Equal(t, expected, sb.String())
		assert.NoError(t, err)
	})

	t.Run("pretty", func(t *testing.T) {
		t.Parallel()

		var sb strings.Builder
		w := JSONWriter(&sb)
		w.SetIndent("", "    ")

		err := w.WriteData(data)

		expected := `[
    {
        "name": "foo",
        "value": "bar"
    }
]
`

		assert.Equal(t, expected, sb.String())
		assert.NoError(t, err)
	})
}

func TestJSONDataWriter_DataWriter(t *testing.T) {
	t.Parallel()

	w := JSONWriter(nil)

	assert.Equal(t, w, w.DataWriter())
}
