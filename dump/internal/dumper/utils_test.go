package dumper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	assert.Equal(t, "''", getValue(""))
	assert.Equal(t, "1", getValue("1"))
	assert.Equal(t, "'foo'", getValue("foo"))
}

func TestEscape(t *testing.T) {
	input := string([]byte{0, '\n', '\r', '\\', '\'', '"', '\032', 'a'})
	expected := `\0\n\r\\\'\"\Za`
	result := escape(input)
	assert.Equal(t, expected, result)
}
