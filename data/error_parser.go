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

func (p *ParserError) GetHeadings(rows []Row) []string {
	return nil // No headings for error parser
}

func (p *ParserError) Parse(input io.Reader) (*RowList, error) {
	return nil, p.Err
}
