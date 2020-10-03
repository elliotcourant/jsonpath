package jsonpath

import (
	"github.com/pkg/errors"
)

type CompiledJsonPath struct {
	actions []jsonAction
}

type pathParser struct {
	path   string
	buffer *tokenBuffer
}

func ParsePath(path string) (CompiledJsonPath, error) {
	parser, err := newPathParser(path)
	if err != nil {
		return CompiledJsonPath{}, err
	}

	actions, err := parser.Parse()

	return CompiledJsonPath{
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

func (p *pathParser) nextAction() (jsonAction, error) {
	token := p.buffer.Peek()
	switch token {
	case dollar:
		p.buffer.Scan()
		return rootAccessAction{}, nil
	case openBracket:
		return p.parseBrackets()
	case period:
		p.buffer.Scan()
		if nextToken := p.buffer.Peek(); nextToken == period {
			// This is a recursive decent.
			return nil, errors.Errorf("recursive decent not supported")
		}

		return p.parseFieldAccess(p.buffer.Scan())
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
	case singleQuotedStringToken, doubleQuotedStringToken:
		// We are accessing a field.
		action, err = p.parseFieldAccess(t)
		if err != nil {
			return nil, err
		}
	case stringToken:
		return nil, errors.Errorf("unexpected string in brackets")
	case integerToken:
		// We are accessing an array index.
		action = arrayIndexAction(int(t))
	case characterToken:
		// We might be performing an operation.
	default:
		return nil, errors.Errorf("somethings broken")
	}

	return action, p.expectCharacterToken(closeBracket)
}

func (p *pathParser) parseFieldAccess(token pathToken) (jsonAction, error) {
	var field string
	switch raw := token.(type) {
	case singleQuotedStringToken:
		field = string(raw)
	case doubleQuotedStringToken:
		field = string(raw)
	case stringToken:
		field = string(raw)
	default:
		return nil, errors.Errorf("unexpected %T parsing field access", token)
	}

	return fieldAccessAction(field), nil
}
