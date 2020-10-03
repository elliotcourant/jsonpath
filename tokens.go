package jsonpath

type (
	pathToken interface {
		PathToken()
	}

	characterToken          byte
	whitespaceToken         byte
	comparisonToken         string
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
	_ pathToken = stringToken("")
	_ pathToken = doubleQuotedStringToken("")
	_ pathToken = singleQuotedStringToken("")
	_ pathToken = nullToken{}
	_ pathToken = booleanToken(false)
	_ pathToken = integerToken(0)
	_ pathToken = decimalToken(0)
)

func (c characterToken) PathToken()          {}
func (w whitespaceToken) PathToken()         {}
func (c comparisonToken) PathToken()         {}
func (s stringToken) PathToken()             {}
func (d doubleQuotedStringToken) PathToken() {}
func (s singleQuotedStringToken) PathToken() {}
func (n nullToken) PathToken()               {}
func (b booleanToken) PathToken()            {}
func (i integerToken) PathToken()            {}
func (d decimalToken) PathToken()            {}

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
	minus        characterToken = '-'
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
