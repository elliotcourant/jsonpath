package jsonpath

import (
	"strconv"

	"github.com/pkg/errors"
)

type (
	tokenBuffer struct {
		tokens      []pathToken
		len, offset int
	}

	pathTokenizer struct {
		path        string
		len, offset int
	}
)

func newPathTokenBuffer(path string) (*tokenBuffer, error) {
	tokenizer := newPathTokenizer(path)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		return nil, err
	}

	return &tokenBuffer{
		tokens: tokens,
		len:    len(tokens),
		offset: 0,
	}, nil
}

func (t *tokenBuffer) Peek() pathToken {
	if t.len < t.offset+1 {
		return eof
	}

	return t.tokens[t.offset]
}

func (t *tokenBuffer) Scan() pathToken {
	if t.len < t.offset+1 {
		return eof
	}

	t.offset++

	return t.tokens[t.offset-1]
}

func newPathTokenizer(path string) *pathTokenizer {
	return &pathTokenizer{
		path:   path,
		len:    len(path),
		offset: 0,
	}
}

func (t *pathTokenizer) Tokenize() ([]pathToken, error) {
	tokens := make([]pathToken, 0)
	for {
		token, err := t.nextToken()
		if err != nil {
			return nil, err
		}

		if token == eof {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// peek will return the next byte in the buffer for tokenization. If the end of
// the buffer has been reached then it will return 0/eof.
func (t *pathTokenizer) peek() byte {
	if t.len < t.offset+1 {
		return 0 // 0 represents the end of the file.
	}

	return t.path[t.offset]
}

// scan will move the buffer cursor forward one place and return the byte from
// that position. scan can move past the end of the buffer if the next token is
// not consistently checked.
func (t *pathTokenizer) scan() byte {
	if t.len < t.offset+1 {
		return 0
	}

	t.offset++
	return t.path[t.offset-1]
}

// scanAndPeek will move the cursor forward one, but then will peek the next
// byte instead of consuming it.
func (t *pathTokenizer) scanAndPeek() byte {
	t.offset++
	return t.peek()
}

// consumeAndReturn will consume the current byte under the cursor and will
// return the provided token. It is assumed that these two things are the same.
func (t *pathTokenizer) consumeAndReturn(token pathToken) (pathToken, error) {
	t.offset++
	return token, nil
}

func (t *pathTokenizer) nextToken() (pathToken, error) {
	character := t.peek()
	switch character {
	case 0:
		// If we get 0 back then that means we have reached the end of the file.
		return eof, nil
	case ' ':
		return t.consumeAndReturn(space)
	case '\t':
		return t.consumeAndReturn(tab)
	case '\r':
		// If the next character is a newline then these two characters are
		// being used together, we want to return a newline instead.
		if nextCharacter := t.scanAndPeek(); nextCharacter == '\n' {
			return t.consumeAndReturn(newline)
		}

		// Otherwise just return "return" as it's own character.
		return rtrn, nil
	case '\n':
		return t.consumeAndReturn(newline)
	case '$':
		return t.consumeAndReturn(dollar)
	case '@':
		return t.consumeAndReturn(at)
	case '*':
		return t.consumeAndReturn(asterisk)
	case ',':
		return t.consumeAndReturn(comma)
	case '.':
		return t.consumeAndReturn(period)
	case '[':
		return t.consumeAndReturn(openBracket)
	case ']':
		return t.consumeAndReturn(closeBracket)
	case '(':
		return t.consumeAndReturn(openParen)
	case ')':
		return t.consumeAndReturn(closeParen)
	case '-':
		return t.consumeAndReturn(minus)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return t.tokenizeNumber()
	case '=':
		// If there are two = in a row then we are making a comparison.
		if nextCharacter := t.scanAndPeek(); nextCharacter == '=' {
			return t.consumeAndReturn(equals)
		}

		// Otherwise just return a single equal character.
		return equal, nil
	case '!':
		// If the next character is a = then we are making a comparison.
		if nextCharacter := t.scanAndPeek(); nextCharacter == '=' {
			return t.consumeAndReturn(notEquals)
		}

		// Otherwise just return the ! character.
		return exclamation, nil
	case '<':
		// If the next character is a = then we are making a comparison.
		if nextCharacter := t.scanAndPeek(); nextCharacter == '=' {
			return t.consumeAndReturn(lessThanOrEqualTo)
		}

		return lessThan, nil
	case '>':
		// If the next character is a = then we are making a comparison.
		if nextCharacter := t.scanAndPeek(); nextCharacter == '=' {
			return t.consumeAndReturn(greaterThanOrEqualTo)
		}

		return greaterThan, nil
	case ':':
		return t.consumeAndReturn(colon)
	case '?':
		return t.consumeAndReturn(question)
	case '\'':
		// Parse single quoted string.
		str, err := t.tokenizeQuotedString(character)
		if err != nil {
			return nil, err
		}

		return singleQuotedStringToken(str), nil
	case '"':
		// Parse single quoted string.
		str, err := t.tokenizeQuotedString(character)
		if err != nil {
			return nil, err
		}

		return doubleQuotedStringToken(str), nil
	default:
		// Parse as a normal string. Potentially a keyword.
		if t.isStringPart(character) {
			return t.tokenizeString()
		}
	}

	return nil, errors.Errorf("unexpected '%s'", string(character))
}

func (t *pathTokenizer) tokenizeNumber() (pathToken, error) {
	startingIndex := t.offset

	isDecimal := false
	for char := t.peek(); t.isNumericPart(char); char = t.scanAndPeek() {
		if char == '.' {
			isDecimal = true
		}
	}

	str := t.path[startingIndex:t.offset]

	if isDecimal {
		decimal, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse '%s' as float", str)
		}

		return decimalToken(decimal), nil
	}

	integer, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse '%s' as integer", str)
	}

	return integerToken(integer), nil
}

func (t *pathTokenizer) isNumericPart(character byte) bool {
	switch character {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		return true
	default:
		return false
	}
}

func (t *pathTokenizer) tokenizeQuotedString(quote byte) (string, error) {
	// Consume the first character, we are assuming that it is the quote char.
	t.offset++

	// Store the current index, we are just going to move the cursor forward
	// until we find another quote character that is not escaped.
	startingIndex := t.offset

ScanLoop:
	for {
		character := t.scan()
		switch character {
		case 0:
			return "", errors.Errorf("unexpected eof parsing string")
		case quote:
			// If there are two of the quotes in a row we want to consider that
			// an escape.
			if t.peek() == quote {
				continue
			} else {
				break ScanLoop
			}
		case '\\':
			// If there are two of the quotes in a row we want to consider that
			// an escape.
			if t.peek() == quote {
				continue
			}
		}
	}

	return t.path[startingIndex : t.offset-1], nil
}

func (t *pathTokenizer) tokenizeString() (pathToken, error) {
	startingIndex := t.offset

	for character := t.peek(); t.isStringPart(character); character = t.scanAndPeek() {
	}

	str := t.path[startingIndex:t.offset]

	// We want to take the opportunity to parse keywords here.
	switch str {
	case "null":
		return nullToken{}, nil
	case "true", "false":
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse '%s' as boolean", str)
		}

		return booleanToken(b), nil
	}

	// If the string is not a token then just return it as a basic string token.
	return stringToken(str), nil
}

func (t *pathTokenizer) isStringPart(character byte) bool {
	return (character >= 'a' && character <= 'z') ||
		(character >= 'A' && character <= 'Z') ||
		character == '_'
}
