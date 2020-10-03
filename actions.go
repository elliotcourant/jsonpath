package jsonpath

import (
	"github.com/pkg/errors"
)

type jsonAction interface {
	Execute(ctx *evalContext) (jsonNode, error)
}

var (
	_ jsonAction = arrayIndexAction(0)
	_ jsonAction = arrayIndexListAction{}
	_ jsonAction = fieldAccessAction("")
	_ jsonAction = rootAccessAction{}
	_ jsonAction = recursiveAction{}
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

type arrayIndexListAction []int

func newArrayIndexListAction(indexes []integerToken) arrayIndexListAction {
	a := make(arrayIndexListAction, len(indexes), len(indexes))
	for i, index := range indexes {
		a[i] = int(index)
	}

	return a
}

func (a arrayIndexListAction) Execute(ctx *evalContext) (jsonNode, error) {
	if isArray(ctx.data) {
		items := make(jsonMutatedArray, len(a), len(a))
		for i, index := range a {
			items[i], _ = getIndex(ctx.data, index)
		}

		return items, nil
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
	} else {
		// We can only try to extract fields from mutated arrays.
		array, ok := ctx.data.(jsonMutatedArray)
		if !ok {
			return nil, errors.Errorf("cannot extract field from array")
		}

		items := make(jsonArray, 0)
		for _, item := range array {
			item, ok := getField(item, string(f))
			if !ok {
				continue
			}
			// TODO (elliotcourant) If the item in the array is not an object
			//  should this fail? Or return an error?

			items = append(items, item)
		}

		return items, nil
	}
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

type recursiveAction struct{}

func (r recursiveAction) Execute(ctx *evalContext) (jsonNode, error) {
	return r.getAllObjects(ctx.data)
}

func (r recursiveAction) getAllObjects(data jsonNode) (jsonMutatedArray, error) {
	items := make(jsonMutatedArray, 0)
	if isArray(data) {
		for _, item := range data.(jsonArray) {
			if isObject(item) || isArray(item) {
				subItems, err := r.getAllObjects(item)
				if err != nil {
					return nil, err
				}

				items = append(items, subItems...)
				items = append(items, item)
			}
		}
	} else {
		for _, item := range data.(jsonObject) {
			if isObject(item) || isArray(item) {
				subItems, err := r.getAllObjects(item)
				if err != nil {
					return nil, err
				}

				items = append(items, subItems...)
				items = append(items, item)
			}
		}
	}

	return items, nil
}
