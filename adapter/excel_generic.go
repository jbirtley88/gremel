package adapter

import (
	"fmt"
	"io"

	"github.com/jbirtley88/gremel/data"
	"github.com/xuri/excelize/v2"
)

// GenericExcelParser is a blunt but effective instrument
//
// It:
//
//   - loads the Excel
//   - uses the first row as headings
//   - parses the rest of the data as rows
type GenericExcelParser struct {
	BaseAdapter
}

func NewGenericExcelParser(ctx data.GremelContext) data.Parser {
	p := &GenericExcelParser{
		BaseAdapter: *NewBaseAdapter("excel", ctx),
	}
	return p
}

func (p *GenericExcelParser) Parse(input io.Reader) (*data.RowList, error) {
	// TODO(john): check the context for any hints on how to parse the Excel.
	// For now, we assume it's:
	//   - Comma-separated
	//   - First row is headings
	spreadsheet, err := excelize.OpenReader(input)
	sheetName := spreadsheet.GetSheetList()[0]
	if p.Ctx != nil {
		if sn, ok := p.Ctx.Values().GetValue("excel.sheetname").(string); ok {
			sheetName = sn
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Parse(%s): open error: %w", p.GetName(), err)
	}

	// TODO(john): deal with multiple worksheets in the Excel file
	// For now, we assume there's only one worksheet called "Sheet1"
	// Step 1: Parse the Excel
	spreadsheetRows, err := spreadsheet.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("Parse(%s): GetRows(): %w", p.GetName(), err)
	}

	// Step 2: convert it to rows
	var rows []data.Row
	var headings []string
	headings = append(headings, spreadsheetRows[0][0:]...)

	for _, ssRow := range spreadsheetRows[1:] {
		row := make(map[string]any)
		for i, value := range ssRow {
			row[headings[i]] = data.InferValue(value)
		}
		rows = append(rows, row)
	}
	return data.NewRowList(rows, p.GetHeadings(rows), nil), nil
}
