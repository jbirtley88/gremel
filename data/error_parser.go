package data

import (
	"context"
	"io"
)

type ParserError struct {
	*BaseParser
	Err error
}

func NewParserError(err error) *ParserError {
	a := &ParserError{
		BaseParser: NewBaseParser(context.TODO(), "error"),
		Err:        err,
	}
	return a
}

func (p *ParserError) GetHeadings(rows []map[string]any) []string {
	return nil // No headings for error parser
}

func (p *ParserError) Parse(input io.Reader) ([]map[string]any, []string, error) {
	return nil, nil, p.Err
}
