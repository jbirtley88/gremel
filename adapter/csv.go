package adapter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jbirtley88/gremel/data"
)

type CSVAdapter struct {
	BaseAdapter
	source io.ReadCloser
}

func NewCSVAdapter(name string, ctx data.GremelContext, source io.ReadCloser) Adapter {
	base := NewBaseAdapter(name, ctx)
	a := &CSVAdapter{
		BaseAdapter: *base,
		source:      source,
	}
	if a.Ctx == nil {
		a.Ctx = data.NewGremelContext(context.TODO())
	}
	return a
}

func (a *CSVAdapter) Load() ([]map[string]any, []string, error) {
	if a.source == nil {
		return nil, nil, fmt.Errorf("Load(): No source provided for CSVAdapter '%s'", a.GetName())
	}

	defer a.source.Close()
	r := csv.NewReader(a.source)
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("Load(%s): read error: %w", a.GetName(), err)
	}

	var rows []map[string]any
	var headings []string
	for _, value := range records[0] {
		headings = append(headings, fmt.Sprint(value))
	}
	for _, record := range records[1:] {
		row := make(map[string]any)
		for i, value := range record {
			row[headings[i]] = value
		}
		rows = append(rows, row)
	}
	return rows, headings, nil
}
