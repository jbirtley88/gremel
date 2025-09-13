package adapter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/stretchr/testify/require"
)

func TestCreateDBFromJSON(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()
	err := CreateDBFromFile(ctx, database, "accounts_json", "../test_resources/accounts_nested.json")
	if err != nil {
		t.Fatalf("failed to create 'accounts_json' table in DB: %v", err)
	}
	err = CreateDBFromFile(ctx, database, "people_json", "../test_resources/people.json")
	if err != nil {
		t.Fatalf("failed to create 'people_json' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts_json",
		"people_json",
	}
	for _, src := range sourceNames {
		filename := strings.Split(src, "_")[0] + ".json"
		f, err := os.Open(fmt.Sprintf("../test_resources/%s", filename))
		require.Nil(t, err)
		require.NotNil(t, f)
		defer f.Close()

		p := NewGenericJsonParser(ctx)
		rows, err := p.Parse(f)
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.NotNil(t, rows.Rows)
		require.NotNil(t, rows.Headings)

		// Step 3: insert the data into the database
		err = database.InsertRows(src, rows.Rows)
		require.Nil(t, err)
	}
}

func TestCreateDBFromCSV(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()
	err := CreateDBFromFile(ctx, database, "accounts_csv", "../test_resources/accounts.csv")
	if err != nil {
		t.Fatalf("failed to create 'accounts_csv' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts_csv",
	}
	for _, src := range sourceNames {
		filename := strings.Split(src, "_")[0] + ".csv"
		f, err := os.Open(fmt.Sprintf("../test_resources/%s", filename))
		require.Nil(t, err)
		require.NotNil(t, f)
		defer f.Close()

		p := NewGenericCSVParser(ctx)
		rows, err := p.Parse(f)
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.NotNil(t, rows.Rows)
		require.NotNil(t, rows.Headings)

		// Step 3: insert the data into the database
		err = database.InsertRows(src, rows.Rows)
		require.Nil(t, err)
	}
}

func TestCreateDBFromExcelSingleSheet(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	ctx.Values().SetValue("excel.sheetname", "Sheet1")
	database := db.GetGremelDB()
	err := CreateDBFromFile(ctx, database, "accounts_excel", "../test_resources/accounts.xlsx")
	if err != nil {
		t.Fatalf("failed to create 'accounts_excel' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts_excel",
	}
	for _, src := range sourceNames {
		filename := strings.Split(src, "_")[0] + ".xlsx"
		f, err := os.Open(fmt.Sprintf("../test_resources/%s", filename))
		require.Nil(t, err)
		require.NotNil(t, f)
		defer f.Close()

		p := NewGenericExcelParser(ctx)
		rows, err := p.Parse(f)
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.NotNil(t, rows.Rows)
		require.NotNil(t, rows.Headings)

		// Step 3: insert the data into the database
		err = database.InsertRows(src, rows.Rows)
		require.Nil(t, err)
	}
}
