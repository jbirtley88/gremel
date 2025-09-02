package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jbirtley88/gremel/data"
)

type JSONAdapter struct {
	BaseAdapter
	source io.ReadCloser
}

func NewJSONAdapter(name string, ctx data.GremelContext, source io.ReadCloser) Adapter {
	base := NewBaseAdapter(name, ctx)
	a := &JSONAdapter{
		BaseAdapter: *base,
		source:      source,
	}
	if a.Ctx == nil {
		a.Ctx = data.NewGremelContext(context.TODO())
	}
	return a
}

func (a *JSONAdapter) Load() ([]map[string]any, []string, error) {
	if a.source == nil {
		return nil, nil, fmt.Errorf("Load(): No source provided for JSONAdapter '%s'", a.GetName())
	}

	defer a.source.Close()
	var rows []map[string]any
	var headings []string
	// Decode the JSON data from the source
	// Assuming the JSON is an array of objects
	// If the JSON structure is different, this part may need to be adjusted
	if err := json.NewDecoder(a.source).Decode(&rows); err != nil {
		return nil, nil, fmt.Errorf("Load(%s): read error: %w", a.GetName(), err)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("Load(%s): no data found in input", a.GetName())
	}
	// Extract headings from the first row.
	// TODO(jb): this order is random, it should be sorted.  Might leave that to post-processing
	// and allow this to concentrate purely on the Load()
	for key := range rows[0] {
		headings = append(headings, key)
	}
	return rows, headings, nil
}
