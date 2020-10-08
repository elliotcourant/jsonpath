package jsonpath

type (
	// Evaluator is a compiled form of a jsonpath. It can run it's jsonpath
	// against any provided json object and only needs to parse the provided
	// json.
	Evaluator struct {
		path    string
		actions []jsonAction
	}

	evalContext struct {
		parent *evalContext
		data   jsonNode
	}
)

// Jsonpath will evaluate the provided jsonpath on the provided json. It does
// not cache either of these parameters. If you are using the same jsonpath
// consistently then you would want to create an Evaluator for that path
// instead. An error is returned if there is a problem parsing the jsonpath or
// if the json could not be parsed.
func Jsonpath(data []byte, path string) ([]interface{}, error) {
	eval, err := NewEvaluator(path)
	if err != nil {
		return nil, err
	}

	return eval.Evaluate(data)
}

// NewEvaluator will compile the provided jsonpath and create an object that can
// run that expression on provided json objects. If the path is not valid then
// an error is returned.
func NewEvaluator(path string) (*Evaluator, error) {
	actions, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	eval := &Evaluator{
		path:    path,
		actions: actions.actions,
	}

	return eval, nil
}

// Evaluate will run the compiled jsonpath against the provided json. It will
// return an array of objects that is the result of the expression or an error
// if something failed to evaluate or if the json was invalid.
func (e *Evaluator) Evaluate(data []byte) ([]interface{}, error) {
	node, err := parseJson(data)
	if err != nil {
		return nil, err
	}

	return e.run(jsonArray{node})
}

func (e *Evaluator) run(root jsonNode) ([]interface{}, error) {
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

	return ctx.data.([]interface{}), nil
}
