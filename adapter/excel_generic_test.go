package adapter

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadExcelHappyPathSingleWorksheet(t *testing.T) {
	ctx := data.NewGremelContext(context.TODO())
	f, err := os.Open("../test_resources/accounts.xlsx")
	require.Nil(t, err)
	require.NotNil(t, f)

	parser := NewGenericExcelParser(ctx)
	require.NotNil(t, parser, "Expected adapter to be registered")

	rows, err := parser.Parse(f)
	require.Nil(t, err, "Expected no error when loading Excel data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, parser.GetHeadings(rows.Rows), "Expected some headings to be loaded")
	assert.Equal(t, 1000, len(rows.Rows), "Expected 1000 rows to be loaded")
	assert.Equal(t, 5, len(parser.GetHeadings(rows.Rows)), "Expected 5 headings to be loaded")
}

// Parse each worksheet individually
func TestLoadExcelHappyPathMultipleWorksheets(t *testing.T) {
	ctx := data.NewGremelContext(context.TODO())
	f, err := os.Open("../test_resources/accounts_multiple_sheets.xlsx")
	require.Nil(t, err)
	require.NotNil(t, f)

	// Sheet 1
	ctx.Values().SetValue("excel.sheetname", "Sheet1")
	parser := NewGenericExcelParser(ctx)
	require.NotNil(t, parser, "Expected adapter to be registered")

	rows, err := parser.Parse(f)
	require.Nil(t, err, "Expected no error when loading Excel data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, parser.GetHeadings(rows.Rows), "Expected some headings to be loaded")
	assert.Equal(t, 1000, len(rows.Rows), "Expected 1000 rows to be loaded")
	assert.Equal(t, 5, len(parser.GetHeadings(rows.Rows)), "Expected 5 headings to be loaded")

	// Sheet 2
	ctx.Values().SetValue("excel.sheetname", "Sheet2")
	parser = NewGenericExcelParser(ctx)
	require.NotNil(t, parser, "Expected adapter to be registered")

	f.Seek(0, io.SeekStart) // rewind the file
	rows, err = parser.Parse(f)
	require.Nil(t, err, "Expected no error when loading Excel data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, parser.GetHeadings(rows.Rows), "Expected some headings to be loaded")
	assert.Equal(t, 829, len(rows.Rows), "Expected 829 rows to be loaded")
	assert.Equal(t, 5, len(parser.GetHeadings(rows.Rows)), "Expected 5 headings to be loaded")

	// Sheet 3
	ctx.Values().SetValue("excel.sheetname", "Sheet3")
	parser = NewGenericExcelParser(ctx)
	require.NotNil(t, parser, "Expected adapter to be registered")

	f.Seek(0, io.SeekStart) // rewind the file
	rows, err = parser.Parse(f)
	require.Nil(t, err, "Expected no error when loading Excel data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, parser.GetHeadings(rows.Rows), "Expected some headings to be loaded")
	assert.Equal(t, 876, len(rows.Rows), "Expected 876 rows to be loaded")
	assert.Equal(t, 5, len(parser.GetHeadings(rows.Rows)), "Expected 5 headings to be loaded")
}
