package jsonpath

import (
	"github.com/pkg/errors"
)

type jsonAction interface {
	Execute(ctx *evalContext) (jsonNode, error)
}

var (
	_ jsonAction = arrayIndexAction(0)
	_ jsonAction = fieldAccessAction("")
	_ jsonAction = rootAccessAction{}
)

type arrayIndexAction int

func (a arrayIndexAction) Execute(ctx *evalContext) (jsonNode, error) {
	if isArray(ctx.data) {
		// TODO (elliotcourant) What should we do with ok here?
		item, _ := getIndex(ctx.data, int(a))
		return item, nil
	}

	return nil, errors.Errorf("item is not an array")
}

type fieldAccessAction string

func (f fieldAccessAction) Execute(ctx *evalContext) (jsonNode, error) {
	if isObject(ctx.data) {
		item, ok := getField(ctx.data, string(f))
		if !ok {
			return nil, nil
		}

		return item, nil
	}

	return nil, errors.Errorf("item is not an object")
}

type rootAccessAction struct{}

func (r rootAccessAction) Execute(ctx *evalContext) (jsonNode, error) {
	// Traverse the evaluation context upwards until we reach the top.
	current := ctx

	for {
		if current.parent == nil {
			return current.data, nil
		}

		current = current.parent
	}
}


