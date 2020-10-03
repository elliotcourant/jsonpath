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
		items := make(jsonArray, 0)
		for _, item := range ctx.data.(jsonArray) {
			if !isArray(item) {
				return nil, errors.Errorf("item is not an array")
			}
			// TODO (elliotcourant) What should we do with ok here?
			result, _ := getIndex(item, int(a))
			items = append(items, result)
		}

		return items, nil
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
		data := ctx.data.(jsonArray)
		items := make(jsonArray, 0, len(a)*len(data))
		for _, item := range data {
			if !isArray(item) {
				return nil, errors.Errorf("item is not an array")
			}
			// TODO (elliotcourant) What should we do with ok here?
			for _, index := range a {
				result, _ := getIndex(item, index)
				items = append(items, result)
			}
		}

		return items, nil
	}

	return nil, errors.Errorf("item is not an array")
}

type fieldAccessAction string

func (f fieldAccessAction) Execute(ctx *evalContext) (jsonNode, error) {
	if isObject(ctx.data) {
		return f.extractField(ctx.data.(jsonObject))
	}

	items := make([]interface{}, 0)

	switch grouping := ctx.data.(type) {
	case jsonArray:
		for _, group := range grouping {
			result, err := f.extractFromGroup(group)
			if err != nil {
				return nil, err
			}

			items = append(items, result...)
		}
	case jsonMutatedArray:
		result, err := f.extractFromGroup(grouping)
		if err != nil {
			return nil, err
		}

		items = append(items, result...)

	}

	return items, nil
}

func (f fieldAccessAction) extractFromGroup(node jsonNode) ([]interface{}, error) {
	items := make([]interface{}, 0)
	switch obj := node.(type) {
	case jsonObject:
		item, err := f.extractField(obj)
		if err != nil {
			return nil, err
		}

		if item == nil {
			break
		}

		items = append(items, item)
	case jsonMutatedArray:
		for _, group := range obj {
			result, err := f.extractFromGroup(group)
			if err != nil {
				return nil, err
			}

			items = append(items, result...)
		}
	}

	return items, nil
}

func (f fieldAccessAction) extractField(data jsonObject) (jsonNode, error) {
	item, ok := data[string(f)]
	if !ok {
		return nil, nil
	}

	return item, nil
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

type wildcardAccessAction struct{}

func (w wildcardAccessAction) Execute(ctx *evalContext) (jsonNode, error) {
	items := make([]interface{}, 0)
	if isArray(ctx.data) {
		for _, sub := range ctx.data.(jsonArray) {
			if isArray(sub) {
				for _, item := range sub.(jsonArray) {
					items = append(items, item)
				}
			}
		}

		return items, nil
	}

	// TODO (elliotcourant) Implement object wildcard.

	return nil, nil
}
