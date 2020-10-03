package jsonpath

type (
	Evaluator struct {
		path    string
		actions []jsonAction
	}

	evalContext struct {
		parent *evalContext
		data   jsonNode
	}
)

func Jsonpath(data []byte, path string) (interface{}, error) {
	eval, err := NewEvaluator(path)
	if err != nil {
		return nil, err
	}

	return eval.Evaluate(data)
}

func NewEvaluator(path string) (*Evaluator, error) {
	actions, err := ParsePath(path)
	if err != nil {
		return nil, err
	}

	eval := &Evaluator{
		path:    path,
		actions: actions.actions,
	}

	return eval, nil
}

func (e *Evaluator) Evaluate(data []byte) (interface{}, error) {
	node, err := parseJson(data)
	if err != nil {
		return nil, err
	}

	return e.run(node)
}

func (e *Evaluator) run(root jsonNode) (interface{}, error) {
	ctx := &evalContext{
		parent: nil,
		data:   root,
	}

	for _, action := range e.actions {
		result, err := action.Execute(ctx)
		if err != nil {
			return nil, err
		}

		ctx = &evalContext{
			parent: ctx,
			data:   result,
		}
	}

	if ctx.data == nil {
		return nil, nil
	}

	return []interface{}{
		ctx.data,
	}, nil
}
