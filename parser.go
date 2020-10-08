package jsonpath

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	compiledJsonPath struct {
		actions []jsonAction
	}

	pathParser struct {
		path   string
		buffer *tokenBuffer
	}

	sliceAccessType uint8
)

const (
	sliceAccessPrecise sliceAccessType = iota
	sliceAccessRangeSimple
	sliceAccessRangeComplex
	sliceAccessList
)

func parsePath(path string) (compiledJsonPath, error) {
	parser, err := newPathParser(path)
	if err != nil {
		return compiledJsonPath{}, err
	}

	actions, err := parser.Parse()

	return compiledJsonPath{
		actions: actions,
	}, err
}

func newPathParser(path string) (*pathParser, error) {
	buffer, err := newPathTokenBuffer(path)
	if err != nil {
		return nil, err
	}

	return &pathParser{
		path:   path,
		buffer: buffer,
	}, nil
}

func (p *pathParser) Parse() ([]jsonAction, error) {
	actions := make([]jsonAction, 0)
	for {
		if p.buffer.Peek() == eof {
			break
		}

		action, err := p.nextAction()
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (p *pathParser) expectCharacterToken(char characterToken) error {
	token := p.buffer.Peek()
	switch token {
	case char:
		p.buffer.Scan() // If we found the token we were looking for, move forward.
		return nil
	default:
		return errors.Errorf("unexpected '%s'", token)
	}
}

func (p *pathParser) consumeMaybe(char characterToken) bool {
	if token := p.buffer.Peek(); token == char {
		p.buffer.Scan()
		return true
	}

	return false
}

func (p *pathParser) consumeInteger() (integerToken, bool) {
	nextToken := p.buffer.Peek()
	integer, ok := nextToken.(integerToken)
	if ok {
		p.buffer.Scan()
	}
	return integer, ok
}

func (p *pathParser) nextAction() (jsonAction, error) {
	token := p.buffer.Peek()
	switch t := token.(type) {
	case stringToken:
		return p.parseFieldAccess(p.buffer.Scan())
	case characterToken:
		switch t {
		case dollar:
			p.buffer.Scan()
			return rootAccessAction{}, nil
		case openBracket:
			return p.parseBrackets()
		case asterisk:
			return p.parseFieldAccess(p.buffer.Scan())
		case period:
			p.buffer.Scan()
			if nextToken := p.buffer.Peek(); nextToken == period {
				// This is a recursive decent.
				p.buffer.Scan()
				return recursiveAction{}, nil
			}

			return p.parseFieldAccess(p.buffer.Scan())
		default:
			return nil, errors.Errorf("unexpected %T: %s", token, token)
		}
	default:
		return nil, errors.Errorf("unexpected %T: %s", token, token)
	}
}

func (p *pathParser) parseBrackets() (jsonAction, error) {
	// Check for array or inner expression.
	p.buffer.Scan()

	var action jsonAction
	var err error

	token := p.buffer.Scan()
	switch t := token.(type) {
	case singleQuotedStringToken, doubleQuotedStringToken, stringToken:
		return p.parseFieldSliceAccess(t)
	case integerToken:
		return p.parseSliceAccess(t)
	case characterToken:
		// We might be performing an operation.
		switch t {
		case colon:
			return p.parseSliceAccess(t)
		case question:
			return nil, errors.Errorf("filtering not implemented")
		case asterisk:
			action, err = p.parseFieldAccess(t)
		default:
			return nil, errors.Errorf("unexpected token '%s'", string(t))
		}
	default:
		return nil, errors.Errorf("somethings broken")
	}

	if err != nil {
		return nil, err
	}

	return action, p.expectCharacterToken(closeBracket)
}

func (p *pathParser) parseFieldSliceAccess(firstToken pathToken) (jsonAction, error) {
	action := make(arrayFieldAccessAction, 0)
	action = append(action, fmt.Sprint(firstToken))
	for {
		nextToken := p.buffer.Scan()
		switch token := nextToken.(type) {
		case stringToken, doubleQuotedStringToken, singleQuotedStringToken:
			field, err := p.parseString(token)
			if err != nil {
				return nil, err
			}

			action = append(action, field)
		case characterToken:
			switch token {
			case comma: // Consume a comma token. We expect this to be comma delimited.
			case closeBracket:
				// Once we finally see a close bracket return the accessor.
				return action, nil
			default:
				return nil, errors.Errorf("unexpected '%s' in slice field access", string(token))
			}
		default:
			return nil, errors.Errorf("unexpected %T in slice field access", token)
		}
	}
}

func (p *pathParser) parseSliceAccess(firstToken pathToken) (jsonAction, error) {
	indexes := make([]integerToken, 0)

	sliceAccessType := sliceAccessPrecise

	currentToken := firstToken
ScanLoop:
	for {
		switch token := currentToken.(type) {
		case integerToken:
			indexes = append(indexes, token)
		case characterToken:
			switch token {
			case closeBracket:
				// If we find a closing bracket then we have reached the end of
				// the index access area. Now parse the data we gathered.
				break ScanLoop
			case colon:
				sliceAccessType = sliceAccessRangeSimple
			case comma:
				sliceAccessType = sliceAccessList

			}
		}

		currentToken = p.buffer.Scan()
	}

	switch sliceAccessType {
	case sliceAccessPrecise:
		return arrayIndexAction(int(indexes[0])), nil
	case sliceAccessList:
		return newArrayIndexListAction(indexes), nil
	}

	return nil, errors.Errorf("bad index operation")
}

func (p *pathParser) parseString(token pathToken) (string, error) {
	var field string
	switch raw := token.(type) {
	case singleQuotedStringToken:
		field = string(raw)
	case doubleQuotedStringToken:
		field = string(raw)
	case stringToken:
		field = string(raw)
	default:
		return "", errors.Errorf("token %T is not a string", token)
	}

	return field, nil
}

func (p *pathParser) parseFieldAccess(token pathToken) (jsonAction, error) {
	var field string
	var err error
	switch raw := token.(type) {
	case singleQuotedStringToken, doubleQuotedStringToken, stringToken:
		field, err = p.parseString(raw)
		if err != nil {
			return nil, err
		}
	case characterToken:
		if raw == asterisk {
			return wildcardAccessAction{}, nil
		}

		return nil, errors.Errorf("unexpected '%s' parsing field access", string(raw))
	default:
		return nil, errors.Errorf("unexpected %T parsing field access", token)
	}

	return fieldAccessAction(field), nil
}
