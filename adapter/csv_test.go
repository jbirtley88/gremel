package adapter

import (
	"context"
	"os"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCSVHappyPath(t *testing.T) {
	ctx := data.NewGremelContext(context.TODO())
	f, err := os.Open("../test_resources/accounts.csv")
	require.Nil(t, err)
	require.NotNil(t, f)

	adapter := NewCSVAdapter("test_csv_adapter", ctx, f)
	require.NotNil(t, adapter, "Expected adapter to be registered")

	rows, headings, err := adapter.Load()
	require.Nil(t, err, "Expected no error when loading CSV data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, headings, "Expected some headings to be loaded")
	assert.Equal(t, 1000, len(rows), "Expected 1000 rows to be loaded")
	assert.Equal(t, 5, len(headings), "Expected 5 headings to be loaded")
}
