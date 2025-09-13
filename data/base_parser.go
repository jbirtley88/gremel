package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type BaseParser struct {
	Name string
	Ctx  context.Context
}

func NewBaseParser(ctx context.Context, name string) *BaseParser {
	a := &BaseParser{
		Name: name,
		Ctx:  ctx,
	}
	return a
}

func (a *BaseParser) GetName() string {
	return a.Name
}

func (a *BaseParser) Parse(input io.Reader) (*RowList, error) {
	// The BaseParser simply unmarshals the input into a data.RowList,
	// which is a complex wrapper for a tuple of ([]Row, headings, error)
	var rows []Row
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&rows); err != nil {
		return NewRowList(nil, nil, err), nil
	}
	if rows == nil {
		return NewRowList(nil, nil, fmt.Errorf("no data found in input")), fmt.Errorf("no data found in input")
	}
	return NewRowList(rows, a.GetHeadings(rows), nil), nil
}

// You may have specific needs for your headings (e.g. the order in which they appear)
// This is a default implementation that returns the keys of the first row sorted alphabetically.
// Override this method in your parser if you need a different behavior.
func (a *BaseParser) GetHeadings(rows []Row) []string {
	if len(rows) == 0 {
		return nil
	}
	headings := make([]string, 0, len(rows[0]))
	for key := range rows[0] {
		headings = append(headings, key)
	}
	sort.Strings(headings)
	return headings
}
