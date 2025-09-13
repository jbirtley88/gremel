package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/jbirtley88/gremel/data"
)

// An Adapter must be cabale of doing two things:
//
// 1. Load data from a remote source, such as CSV files / JSON fiels
// 2. Provide a string representation of the data that can be displayed through the pseudo-filesystem.
type Adapter interface {
	GetName() string
	Load() ([]map[string]any, []string, error)
}

type BaseAdapter struct {
	Name string
	Ctx  data.GremelContext
}

func NewBaseAdapter(name string, ctx data.GremelContext) *BaseAdapter {
	a := &BaseAdapter{
		Name: name,
		Ctx:  ctx,
	}
	if a.Ctx == nil {
		a.Ctx = data.NewGremelContext(context.TODO())
	}
	return a
}

func (a *BaseAdapter) GetName() string {
	return a.Name
}

func (a *BaseAdapter) Load() ([]map[string]any, []string, error) {
	// The BaseAdapter does nothing by default.
	// Subclasses should override this method to implement specific loading logic.
	return nil, nil, fmt.Errorf("Load(): Default implementation does nothing")
}

func (p *BaseAdapter) GetHeadings(rows []data.Row) []string {
	headings := []string{}

	// Check the context
	if p.Ctx != nil {
		if selectValue := p.Ctx.Values().GetString("select"); selectValue != "" && selectValue != "*" {
			headings = []string{}
			for _, columnName := range strings.Split(selectValue, ",") {
				headings = append(headings, strings.TrimSpace(columnName))
			}
		}
	}

	// Nothing in the context.
	// We need to grab tem from the map.
	// Sadly, keys are not ordered so we're going to get a pseudo-random order
	if len(headings) == 0 {
		for k := range rows[0] {
			headings = append(headings, k)
		}
	}

	return headings
}
