package sliceutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendIfMissing(t *testing.T) {
	slice := []string{
		"foo",
	}

	slice = AppendIfMissing(slice, "foo")
	slice = AppendIfMissing(slice, "bar")

	assert.Equal(t, []string{"foo", "bar"}, slice)
}
