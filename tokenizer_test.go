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
		`$.store.book[*].title`,
		`$.store..price`,
		`$..book[2]`,
		`$..book[(@.length-1)]`,
		`$..book[-1:]`,
		`$..book[0,1]`,
		`$..book[:2]`,
		`$..book[?(@.isbn)]`,
		`$..book[?(@.price<10)]`,
		`$..*`,
		`$.store.book[?(@.price < $.expensive)]`,
	}

	for _, path := range paths {
		tokenizer := newPathTokenizer(path)

		tokens, err := tokenizer.Tokenize()
		assert.NoError(t, err)
		assert.NotEmpty(t, tokens)
	}
}
