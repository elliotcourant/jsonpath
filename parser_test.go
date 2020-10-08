package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		path := "[0].name"
		compiled, err := parsePath(path)
		assert.NoError(t, err)

		assert.NotEmpty(t, compiled)
	})
}
