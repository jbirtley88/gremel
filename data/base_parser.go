package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/jbirtley88/gremel/helper"
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

func (a *BaseParser) Parse(input io.Reader) ([]map[string]any, []string, error) {
	// The BaseParser simply unmarshals the input into a map[string][]any
	var result []map[string]any
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&result); err != nil {
		return nil, nil, err
	}
	if result == nil {
		return nil, nil, fmt.Errorf("no data found in input")
	}
	return result, a.GetHeadings(result), nil
}

// You may have specific needs for your headings (e.g. the order in which they appear)
// This is a default implementation that returns the keys of the first row sorted alphabetically.
// Override this method in your parser if you need a different behavior.
func (a *BaseParser) GetHeadings(rows []map[string]any) []string {
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

// Select operates conceptually exactly the same as a SELECT in SQL
func (a *BaseParser) Select(rows []map[string]any) ([]map[string]any, error) {
	var selectClause string
	var ok bool
	selectClause, ok = a.Ctx.Value("select").(string)
	if !ok || selectClause == "" || selectClause == "*" {
		return rows, nil
	}

	selector := helper.NewBaseSelector(a.Ctx)
	return selector.Select(selectClause, rows)
}

// Where operates conceptually exactly the same as a WHERE in SQL
func (a *BaseParser) Where(rows []map[string]any) ([]map[string]any, error) {
	var whereClause string
	var ok bool
	whereClause, ok = a.Ctx.Value("where").(string)
	if !ok || whereClause == "" {
		return rows, nil
	}

	selector := helper.NewBaseSelector(a.Ctx)
	return selector.Where(rows, whereClause)
}

// Order operates conceptually exactly the same as a ORDER BY in SQL
func (a *BaseParser) Order(rows []map[string]any) ([]map[string]any, error) {
	var orderBy string
	var ok bool
	orderBy, ok = a.Ctx.Value("order").(string)
	if !ok || orderBy == "" {
		return rows, nil
	}

	selector := helper.NewBaseSelector(a.Ctx)
	return selector.Order(rows, orderBy)
}
