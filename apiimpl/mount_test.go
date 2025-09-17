package apiimpl

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/jbirtley88/gremel/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMountLocalFile(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())

	// Step 2: create the 'weblogs' table from the logfile
	err := Mount(ctx, "weblogs", "../test_resources/weblogs.log")
	require.Nil(t, err)
}

func TestMountCSVParsesFloatCorrectly(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())

	// Step 2: create the 'weblogs' table from the logfile
	err := Mount(ctx, "weblogs", "../test_resources/accounts.csv")
	require.Nil(t, err)
}

func TestMountFileUrl(t *testing.T) {
	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())

	// Step 2: create the 'weblogs' table from the logfile
	fileUrl, err := filepath.Abs("../test_resources/weblogs.log")
	require.Nil(t, err)
	err = Mount(ctx, "weblogs", fileUrl)
	require.Nil(t, err)
}

// mockHttpHelper implements helper.HttpHelper for testing
type mockHttpHelper struct {
	responseCode int
	responseBody string
	shouldError  bool
	errorMsg     string
}

// Compile-time check to ensure mockHttpHelper implements helper.HttpHelper
var _ helper.HttpHelper = (*mockHttpHelper)(nil)

func (m *mockHttpHelper) Get(url string) (code int, body io.ReadCloser, err error) {
	if m.shouldError {
		return 0, nil, errors.New(m.errorMsg)
	}
	return m.responseCode, io.NopCloser(strings.NewReader(m.responseBody)), nil
}

func TestMountHttpUrl(t *testing.T) {
	// Read the contents of people.json to use as mock response
	peopleJsonPath := "../test_resources/people.json"
	peopleJsonBytes, err := os.ReadFile(peopleJsonPath)
	require.NoError(t, err)

	// Create a mock HttpHelper
	originalHttpHelper := httpHelper
	defer func() {
		// Restore original httpHelper after test
		httpHelper = originalHttpHelper
	}()

	// Replace with mock that returns people.json content
	httpHelper = &mockHttpHelper{
		responseCode: 200,
		responseBody: string(peopleJsonBytes),
		shouldError:  false,
	}

	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	// Step 2: mount an HTTP URL that should return people.json data
	testUrl := "https://api.example.com/people"
	err = Mount(ctx, "people", testUrl)
	require.NoError(t, err)

	// Step 3: verify the mount was registered
	mountInfo, err := GetMount(ctx, "people")
	require.NoError(t, err)
	assert.Equal(t, testUrl, mountInfo["people"])

	// Step 4: verify that data was actually loaded by querying the table
	rows, columns, err := database.Query("SELECT COUNT(*) as count FROM people")
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Contains(t, columns, "count")
	// people.json has 1000 records
	assert.Equal(t, int64(1000), rows[0]["count"])

	// Step 5: verify some actual data content
	rows, columns, err = database.Query("SELECT id, fullname, email FROM people WHERE id = 1")
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, []string{"id", "fullname", "email"}, columns)
	assert.Equal(t, int64(1), rows[0]["id"]) // ID should be int64 after proper type inference
	assert.Equal(t, "Marcellina Benedicto", rows[0]["fullname"])
	assert.Equal(t, "mbenedicto0@earthlink.net", rows[0]["email"])
}

func TestMountHttpUrlError(t *testing.T) {
	// Create a mock HttpHelper that returns an error
	originalHttpHelper := httpHelper
	defer func() {
		// Restore original httpHelper after test
		httpHelper = originalHttpHelper
	}()

	// Replace with mock that returns an error
	httpHelper = &mockHttpHelper{
		shouldError: true,
		errorMsg:    "connection refused",
	}

	// Step 1: fire up the DB
	ctx := data.NewGremelContext(context.Background())

	// Step 2: try to mount an HTTP URL that will fail
	testUrl := "https://example.com/nonexistent.json"
	err := Mount(ctx, "test_table", testUrl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

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
