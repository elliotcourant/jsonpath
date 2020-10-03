package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestJson = `{
  "firstName": "John",
  "lastName" : "doe",
  "age"      : 26,
  "address"  : {
    "streetAddress": "naist street",
    "city"         : "Nara",
    "postalCode"   : "630-0192"
  },
  "phoneNumbers": [
    {
      "type"  : "iPhone",
      "number": "0123-4567-8888"
    },
    {
      "type"  : "home",
      "number": "0123-4567-8910"
    },
	{
      "type": "mobile",
      "number": "0913-8532-8492"
	}
  ]
}`

func EvaluateOnTestJson(t *testing.T, path string) interface{} {
	result, err := Jsonpath([]byte(TestJson), path)
	require.NoError(t, err, "should succeed")
	return result
}

func TestEvaluator_Evaluate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.firstName")
		assert.Equal(t, "John", result)
	})

	t.Run("simple bracket", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$['firstName']")
		assert.Equal(t, "John", result)
	})

	t.Run("sub object", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.address.streetAddress")
		assert.Equal(t, "naist street", result)
	})

	t.Run("sub object brackets", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$['address']['streetAddress']")
		assert.Equal(t, "naist street", result)
	})

	t.Run("array index", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.phoneNumbers[0].type")
		assert.Equal(t, "iPhone", result)
	})

	t.Run("array indexes", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.phoneNumbers[0,1].type")
		assert.Equal(t, []interface{}{
			"iPhone",
			"home",
		}, result)
	})

	t.Run("recursive phone type", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$..type")
		assert.Equal(t, []interface{}{
			"iPhone",
			"home",
			"mobile",
		}, result)
	})

	t.Run("recursive phone type", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$..streetAddress")
		assert.Equal(t, []interface{}{
			"naist street",
		}, result)
	})
}
