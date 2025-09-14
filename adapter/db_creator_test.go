package adapter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTablesFromJSON(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()
	err := CreateTableFromFile(ctx, database, "accounts", "json", "../test_resources/accounts_nested.json")
	if err != nil {
		t.Fatalf("failed to create 'accounts' table in DB: %v", err)
	}
	err = CreateTableFromFile(ctx, database, "people", "json", "../test_resources/people.json")
	if err != nil {
		t.Fatalf("failed to create 'people' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts",
		"people",
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

func TestCreateTablesFromCSV(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()
	err := CreateTableFromFile(ctx, database, "accounts", "csv", "../test_resources/accounts.csv")
	if err != nil {
		t.Fatalf("failed to create 'accounts' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts",
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

func TestCreateTablesFromExcelSingleSheet(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	ctx.Values().SetValue("excel.sheetname", "Sheet1")
	database := db.GetGremelDB()
	err := CreateTableFromFile(ctx, database, "accounts", "xlsx", "../test_resources/accounts.xlsx")
	if err != nil {
		t.Fatalf("failed to create 'accounts' table in DB: %v", err)
	}

	// Step 2: populate the data
	sourceNames := []string{
		"accounts",
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

func TestGetHighLatencyDatacenter(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	// Step 2: create the 'weblogs' table from the logfile
	// ctx.Values().SetValue("log.format", "combined")
	err := CreateTableFromFile(ctx, database, "weblogs", "log", "../test_resources/weblogs.log")
	if err != nil {
		t.Fatalf("failed to create 'weblogs' table in DB: %v", err)
	}

	// Step 3: create the 'ipaddresses' table from the excel spreadsheet
	ctx.Values().SetValue("excel.sheetname", "ipaddresses")
	err = CreateTableFromFile(ctx, database, "ipaddresses", "xlsx", "../test_resources/ipaddresses.xlsx")
	if err != nil {
		t.Fatalf("failed to create 'ipaddresses' table in DB: %v", err)
	}

	sqlQuery := `SELECT
  i.datacenter,
  COUNT(DISTINCT CASE WHEN CAST(w.latency AS INTEGER) > 2000 THEN i.ip END) AS "latency>2000"
FROM ipaddresses AS i
LEFT JOIN weblogs AS w
  ON w.host = i.ip
WHERE w.request LIKE  'GET /api/foo%'
GROUP BY i.datacenter
ORDER BY i.datacenter`
	rows, columns, err := database.Query(sqlQuery)
	if err != nil {
		t.Fatalf("failed to execute query: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected some results, got none")
	}
	assert.Equal(t, int64(237), rows[0]["latency>2000"])
	assert.Equal(t, "datacenter1", rows[0]["datacenter"])
	assert.Equal(t, int64(0), rows[1]["latency>2000"])
	assert.Equal(t, int64(0), rows[2]["latency>2000"])
	assert.Equal(t, int64(0), rows[3]["latency>2000"])
	assert.Equal(t, []string{"datacenter", "latency>2000"}, columns)
}
