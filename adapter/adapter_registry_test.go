package adapter

import (
	"context"
	"os"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndGetCSV(t *testing.T) {
	err := GetAdapterRegistry().Register("csv", NewCSVAdapter)
	require.Nil(t, err, "Expected no error when registering adapter")

	f, err := os.Open("../test_resources/people.csv")
	require.Nil(t, err)
	require.NotNil(t, f)

	csvAdapter := GetAdapterRegistry().Get("csv", data.NewGremelContext(context.TODO()), f)
	require.NotNil(t, csvAdapter, "Expected adapter to be registered")
	rows, headings, err := csvAdapter.Load()
	require.Nil(t, err, "Expected no error when loading CSV data")
	require.NotEmpty(t, rows, "Expected some rows to be loaded")
	require.NotEmpty(t, headings, "Expected some headings to be loaded")
	assert.Equal(t, 1000, len(rows), "Expected 1000 rows to be loaded")
	assert.Equal(t, 6, len(headings), "Expected 6 headings to be loaded")
}
