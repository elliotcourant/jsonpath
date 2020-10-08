package jsonpath

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func parseJsonString(input string) (jsonNode, error) {
	return parseJson([]byte(input))
}

func parseJson(input []byte) (jsonNode, error) {
	var data jsonNode
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal input")
	}

	return data, nil
}

type (
	jsonNode         interface{}
	jsonArray        = []interface{}
	jsonObject       = map[string]interface{}
	jsonMutatedArray []jsonNode
)

func isArray(data jsonNode) bool {
	// We can type check this against an array of interfaces instead of using
	// reflect. This is probably faster.
	_, ok := data.(jsonArray)
	return ok
}

func isObject(data jsonNode) bool {
	// We can type check this against an array of interfaces instead of using
	// reflect. This is probably faster.
	// TODO (elliotcourant) Can integers be keys in json? This assertion would
	//  fail if they are.
	_, ok := data.(jsonObject)
	return ok
}

func getIndex(data jsonNode, index int) (jsonNode, bool) {
	array, ok := data.(jsonArray)
	if !ok {
		return nil, false
	}

	return jsonNode(array[index]), true
}
