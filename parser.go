package jsonpath

type pathParser struct {
	path   string
	buffer *tokenBuffer
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

func (p *pathParser) Parse() error {
	return nil
}
