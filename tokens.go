package jsonpath

type (
	pathToken interface {
		PathToken()
	}

	characterToken          byte
	whitespaceToken         byte
	comparisonToken         string
	fieldToken              string
	stringToken             string
	doubleQuotedStringToken string
	singleQuotedStringToken string
	nullToken               struct{}
	booleanToken            bool
	integerToken            int64
	decimalToken            float64
)

var (
	_ pathToken = characterToken(0)
	_ pathToken = whitespaceToken(0)
	_ pathToken = comparisonToken("")
	_ pathToken = fieldToken("")
	_ pathToken = stringToken("")
	_ pathToken = doubleQuotedStringToken("")
	_ pathToken = singleQuotedStringToken("")
	_ pathToken = nullToken{}
	_ pathToken = booleanToken(false)
	_ pathToken = integerToken(0)
	_ pathToken = decimalToken(0)
)

func (c characterToken) PathToken() {
	panic("implement me")
}

func (w whitespaceToken) PathToken() {
	panic("implement me")
}

func (c comparisonToken) PathToken() {
	panic("implement me")
}

func (f fieldToken) PathToken() {
	panic("implement me")
}

func (s stringToken) PathToken() {
	panic("implement me")
}

func (d doubleQuotedStringToken) PathToken() {
	panic("implement me")
}

func (s singleQuotedStringToken) PathToken() {
	panic("implement me")
}

func (n nullToken) PathToken() {
	panic("implement me")
}

func (b booleanToken) PathToken() {
	panic("implement me")
}

func (i integerToken) PathToken() {
	panic("implement me")
}

func (d decimalToken) PathToken() {
	panic("implement me")
}

const (
	eof          characterToken = 0
	dollar       characterToken = '$'
	at           characterToken = '@'
	asterisk     characterToken = '*'
	exclamation  characterToken = '!'
	question     characterToken = '?'
	colon        characterToken = ':'
	period       characterToken = '.'
	comma        characterToken = ','
	openBracket  characterToken = '['
	closeBracket characterToken = ']'
	openParen    characterToken = '('
	closeParen   characterToken = ')'
	equal        characterToken = '=' // A single = is not a comparison.
)

const (
	space   whitespaceToken = ' '
	newline whitespaceToken = '\n'
	rtrn    whitespaceToken = '\r' // spelled wrong intentionally, return is reserved
	tab     whitespaceToken = '\t'
)

const (
	equals               comparisonToken = "=="
	notEquals            comparisonToken = "!="
	lessThan             comparisonToken = "<"
	lessThanOrEqualTo    comparisonToken = "<="
	greaterThan          comparisonToken = ">"
	greaterThanOrEqualTo comparisonToken = ">="
)
