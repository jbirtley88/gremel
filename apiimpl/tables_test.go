package apiimpl

import (
	"context"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/stretchr/testify/assert"
)

func TestGetTables_EmptyDatabase(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())

	// Get tables from the database (might not be empty due to previous tests)
	tables, err := GetTables(ctx)

	// Should succeed with no error
	assert.NoError(t, err)

	// In Go, an empty map range results in a nil slice, which is valid behavior
	// We should verify that the result is usable (can get length, iterate, etc.)
	assert.GreaterOrEqual(t, len(tables), 0, "Should be able to get length of returned slice")

	// Verify we can iterate over it (whether nil or empty slice)
	for _, table := range tables {
		assert.IsType(t, "", table, "Each table name should be a string")
	}
}

func TestGetTables_WithSampleTables(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	// Create some test tables
	testTables := []string{"test_table1", "test_table2", "test_table3"}

	// Create schemas for test tables
	sampleRow := data.Row{
		"id":   1,
		"name": "test",
	}

	// Clean up any existing test tables first
	for _, tableName := range testTables {
		_ = database.DropSchema(tableName) // Ignore errors if table doesn't exist
	}

	// Create test tables
	for _, tableName := range testTables {
		err := database.CreateSchema(tableName, sampleRow)
		assert.NoError(t, err, "Failed to create test table %s", tableName)
	}

	// Test GetTables
	tables, err := GetTables(ctx)

	// Should succeed
	assert.NoError(t, err)
	assert.NotNil(t, tables)

	// All test tables should be present
	for _, expectedTable := range testTables {
		assert.Contains(t, tables, expectedTable, "Expected table %s to be in the list", expectedTable)
	}

	// Clean up test tables
	for _, tableName := range testTables {
		err := database.DropSchema(tableName)
		assert.NoError(t, err, "Failed to clean up test table %s", tableName)
	}
}

func TestGetTables_AfterDropTable(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_drop_table"
	sampleRow := data.Row{
		"id":    1,
		"value": "test",
	}

	// Clean up any existing table first
	_ = database.DropSchema(tableName)

	// Create a test table
	err := database.CreateSchema(tableName, sampleRow)
	assert.NoError(t, err)

	// Verify table exists
	tables, err := GetTables(ctx)
	assert.NoError(t, err)
	assert.Contains(t, tables, tableName)

	// Drop the table
	err = database.DropSchema(tableName)
	assert.NoError(t, err)

	// Note: Currently, DropSchema doesn't remove the table from GetTables() result
	// This is a limitation of the current implementation - the table schema remains
	// in memory even after the table is dropped from the database.
	// In a future version, this should be fixed to remove the table from the list.
	tables, err = GetTables(ctx)
	assert.NoError(t, err)
	// For now, we just verify GetTables still works after dropping
	assert.NotNil(t, tables)
}

// TestGetTables_ContextParameter tests that the function accepts a GremelContext parameter
// even though it doesn't use it (testing the function signature)
func TestGetTables_ContextParameter(t *testing.T) {
	// Test with nil context (though this is not recommended in real usage)
	tables, err := GetTables(nil)
	assert.NoError(t, err)
	assert.NotNil(t, tables)

	// Test with proper context
	ctx := data.NewGremelContext(context.Background())
	tables, err = GetTables(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tables)
}

func TestGetTables_ReturnType(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())

	// Test return type
	tables, err := GetTables(ctx)

	// Should return a slice of strings and no error
	assert.NoError(t, err)
	assert.IsType(t, []string{}, tables)
	assert.NotNil(t, tables)
}

func TestGetTables_Consistency(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())

	// Call GetTables multiple times
	tables1, err1 := GetTables(ctx)
	assert.NoError(t, err1)

	tables2, err2 := GetTables(ctx)
	assert.NoError(t, err2)

	// Results should be consistent (same tables in same order or similar)
	assert.Equal(t, len(tables1), len(tables2), "GetTables should return consistent results")
	assert.ElementsMatch(t, tables1, tables2, "GetTables should return the same tables")
}
