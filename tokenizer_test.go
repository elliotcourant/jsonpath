package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	paths := []string{
		`$.phoneNumbers[:1].type`,
		`$.store.book[*].author`,
		`$..author`,
	}

	for _, path := range paths {
		tokenizer := newPathTokenizer(path)

		tokens, err := tokenizer.Tokenize()
		assert.NoError(t, err)
		assert.NotEmpty(t, tokens)
	}
}
