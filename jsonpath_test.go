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

	t.Run("sub object", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.address.streetAddress")
		assert.Equal(t, "naist street", result)
	})

	t.Run("array index", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.phoneNumbers[0].type")
		assert.Equal(t, "iPhone", result)
	})
}
