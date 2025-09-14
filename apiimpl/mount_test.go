package apiimpl

import (
	"context"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/stretchr/testify/assert"
)

func TestGetHighLatencyDatacenter(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	// Step 2: create the 'weblogs' table from the logfile
	err := MountFile(ctx, "weblogs", "../test_resources/weblogs.log")
	if err != nil {
		t.Fatalf("failed to create 'weblogs' table in DB: %v", err)
	}

	// Step 3: create the 'ipaddresses' table from the excel spreadsheet
	ctx.Values().SetValue("excel.sheetname", "ipaddresses")
	err = MountFile(ctx, "ipaddresses", "../test_resources/ipaddresses.xlsx")
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
