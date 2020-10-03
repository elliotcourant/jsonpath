package jsonpath

type jsonAction interface {
	Action()
}

var (
	_ jsonAction = arrayIndexAction(0)
)

type arrayIndexAction int

func (a arrayIndexAction) Action() {}



