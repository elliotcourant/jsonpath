package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsArrayObject(t *testing.T) {
	t.Run("is array", func(t *testing.T) {
		input := `["test", "test"]`
		data, err := parseJsonString(input)
		assert.NoError(t, err)

		assert.True(t, isArray(data))
		assert.False(t, isObject(data))
	})

	t.Run("is not array", func(t *testing.T) {
		input := `{"test": true}`
		data, err := parseJsonString(input)
		assert.NoError(t, err)

		assert.False(t, isArray(data))
		assert.True(t, isObject(data))
	})
}
