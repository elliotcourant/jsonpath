package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenImplementation(t *testing.T) {
	tokens := []pathToken{
		characterToken(0),
		whitespaceToken(0),
		comparisonToken(""),
		stringToken(""),
		doubleQuotedStringToken(""),
		nullToken{},
		booleanToken(false),
		integerToken(0),
		decimalToken(0),
	}

	for _, token := range tokens {
		token.PathToken()
		assert.Implements(t, (*pathToken)(nil), token)
	}
}
