package adapter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jbirtley88/gremel/data"
)

// GenericCSVParser is a blunt but effective instrument
//
// It:
//
//   - loads the CSV
//   - uses the first row as headings
//   - parses the rest of the data as rows
type GenericCSVParser struct {
	BaseAdapter
}

func NewGenericCSVParser(ctx data.GremelContext) data.Parser {
	p := &GenericCSVParser{
		BaseAdapter: *NewBaseAdapter("csv", ctx),
	}
	return p
}

func (p *GenericCSVParser) Parse(input io.Reader) (*data.RowList, error) {
	// TODO(john): check the context for any hints on how to parse the CSV.
	// For now, we assume it's:
	//   - Comma-separated
	//   - First row is headings
	if p.Ctx != nil {
	}

	// Step 1: Parse the CSV
	r := csv.NewReader(input)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Parse(%s): read error: %w", p.GetName(), err)
	}

	// Step 2: convert it to rows
	var rows []data.Row
	var headings []string
	for _, value := range records[0] {
		headings = append(headings, fmt.Sprint(value))
	}
	for _, record := range records[1:] {
		row := make(map[string]any)
		for i, value := range record {
			row[headings[i]] = data.InferValue(value)
		}
		rows = append(rows, row)
	}
	return data.NewRowList(rows, p.GetHeadings(rows), nil), nil
}
