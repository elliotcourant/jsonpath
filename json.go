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

type jsonNode interface{}

type jsonArray = []interface{}

type jsonObject = map[string]interface{}

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

func getField(data jsonNode, fieldName string) (jsonNode, bool) {
	object, ok := data.(jsonObject)
	if !ok {
		// Fail silently.
		// TODO (elliotcourant) Maybe log something? Test this with an actual
		//  jsonpath.
		return nil, false
	}

	result, ok := object[fieldName]

	return jsonNode(result), ok
}

func getIndex(data jsonNode, index int) (jsonNode, bool) {
	array, ok := data.(jsonArray)
	if !ok {
		return nil, false
	}

	return jsonNode(array[index]), true
}
