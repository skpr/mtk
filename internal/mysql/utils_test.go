package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	val, err := getValue("")
	assert.NoError(t, err)
	assert.Equal(t, "''", val)

	val, err = getValue("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", val)

	val, err = getValue("foo")
	assert.NoError(t, err)
	assert.Equal(t, "'foo'", val)
}

func TestEscape(t *testing.T) {
	input := string([]byte{0, '\n', '\r', '\\', '\'', '"', '\032', 'a'})
	expected := `\0\n\r\\\'\"\Za`
	result, err := escape(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
