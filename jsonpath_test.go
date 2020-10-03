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

func MustFailOnTestJson(t *testing.T, path string) error {
	result, err := Jsonpath([]byte(TestJson), path)
	require.Error(t, err, "there should be an error")
	require.Nil(t, result, "result should be nil")
	return err
}

func TestEvaluator_Evaluate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.firstName")
		assert.Equal(t, "John", result)
	})

	t.Run("field does not exist", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$.iDontExist")
		assert.Nil(t, result)
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

	t.Run("array index fails on object", func(t *testing.T) {
		err := MustFailOnTestJson(t, "[0]")
		assert.EqualError(t, err, "item is not an array")
	})

	t.Run("array index list fails on object", func(t *testing.T) {
		err := MustFailOnTestJson(t, "[0,1]")
		assert.EqualError(t, err, "item is not an array")
	})

	t.Run("cannot access field on non-mutated array", func(t *testing.T) {
		err := MustFailOnTestJson(t, "$.phoneNumbers.type")
		assert.EqualError(t, err, "cannot extract field from array")
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

	t.Run("recursive first item", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$..type[0]")
		assert.Equal(t, "iPhone", result)
	})

	t.Run("recursive phone type", func(t *testing.T) {
		result := EvaluateOnTestJson(t, "$..streetAddress")
		assert.Equal(t, []interface{}{
			"naist street",
		}, result)
	})
}

func TestJsonpath(t *testing.T) {
	t.Run("bad path", func(t *testing.T) {
		result, err := Jsonpath(nil, `"thing`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("bad json", func(t *testing.T) {
		result, err := Jsonpath([]byte(`{"test:true}`), `test`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
