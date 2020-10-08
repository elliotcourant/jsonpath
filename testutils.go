package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertMatchingStringArray(t *testing.T, expected []string, result []interface{}) {
	actual := make([]string, len(result), len(result))
	for i, item := range result {
		actual[i] = item.(string)
	}

	if !assert.True(t, sameStringSlice(expected, actual)) {
		assert.EqualValues(t, expected, result)
	}
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}
