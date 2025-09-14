package adapter

import (
	"fmt"
	"io"
	"os"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
)

func CreateTableFromFile(ctx data.GremelContext, database db.GremelDB, tableName string, fileType string, datafile string) error {
	// Step 1: Parse the underlying file into rows
	var parser data.Parser
	switch fileType {
	case "json":
		// Parse JSON file
		parser = NewGenericJsonParser(ctx)

	case "csv":
		// Parse CSV file
		parser = NewGenericCSVParser(ctx)

	case "log":
		// Parse log file
		parser = NewGenericLogParser(ctx)

	case "xlsx", "xls":
		// Parse Excel file
		parser = NewGenericExcelParser(ctx)

	default:
		return fmt.Errorf("CreateDB(%s): unsupported file type: %s", datafile, fileType)
	}

	// Step 2: Parse the data into rows
	f, err := os.Open(datafile)
	if err != nil {
		return fmt.Errorf("CreateDB(%s): failed to open file: %w", datafile, err)
	}
	defer f.Close()

	err = CreateTableFromReader(ctx, database, tableName, f, parser)
	if err != nil {
		return fmt.Errorf("CreateDB(%s): failed to create DB from reader: %w", datafile, err)
	}
	return nil
}

func CreateTableFromReader(ctx data.GremelContext, database db.GremelDB, tableName string, input io.Reader, parser data.Parser) error {
	rows, err := parser.Parse(input)
	if err != nil {
		return fmt.Errorf("CreateDBFromReader(%s): failed to parse data: %w", tableName, err)
	}
	ctx.Values().SetValue(tableName+".headings", rows.Headings)

	// Create the database schema based on the headings in the first row
	if len(rows.Rows) == 0 {
		return fmt.Errorf("CreateDBFromReader(%s): no data rows found", tableName)
	}

	// Make sure that we don't already have a table of this name

	err = database.CreateSchema(tableName, rows.Rows[0])
	if err != nil {
		return fmt.Errorf("CreateDBFromReader(%s): failed to create schema: %w", tableName, err)
	}

	// Insert all the rows
	for rowNum, rowMap := range rows.Rows {
		// Ensure all rows have the same columns as the first row
		err = database.InsertRows(tableName, data.NewSingleRow(rowMap))
		if err != nil {
			return fmt.Errorf("CreateDBFromReader(%s): failed to insert row %d: %w", tableName, rowNum, err)
		}
	}
	return nil
}
