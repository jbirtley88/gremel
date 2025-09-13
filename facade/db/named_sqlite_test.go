package db

import (
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedSQLiteGremelDB_SeparateDatabases(t *testing.T) {
	// Create two named databases
	db1 := newNamedSQLiteGremelDB("test_db1")
	db2 := newNamedSQLiteGremelDB("test_db2")

	// Ensure they are not error databases
	_, isError1 := db1.(*ErrorGremelDB)
	_, isError2 := db2.(*ErrorGremelDB)
	require.False(t, isError1, "db1 should not be an error database")
	require.False(t, isError2, "db2 should not be an error database")

	defer func() {
		err1 := db1.Close()
		err2 := db2.Close()
		assert.NoError(t, err1, "Failed to close db1")
		assert.NoError(t, err2, "Failed to close db2")
	}()

	// Create different schemas in each database
	sampleRow1 := data.Row{
		"id":   1,
		"name": "test",
		"type": "A",
	}

	sampleRow2 := data.Row{
		"id":          1,
		"description": "test",
		"value":       42.5,
	}

	// Create schema in db1
	err := db1.CreateSchema("users", sampleRow1)
	require.NoError(t, err, "Failed to create schema in db1")

	// Create schema in db2
	err = db2.CreateSchema("products", sampleRow2)
	require.NoError(t, err, "Failed to create schema in db2")

	// Insert data into db1
	rows1 := []data.Row{
		{"id": 1, "name": "Alice", "type": "admin"},
		{"id": 2, "name": "Bob", "type": "user"},
	}
	err = db1.InsertRows("users", rows1)
	require.NoError(t, err, "Failed to insert into db1")

	// Insert data into db2
	rows2 := []data.Row{
		{"id": 1, "description": "Widget A", "value": 19.99},
		{"id": 2, "description": "Widget B", "value": 29.99},
	}
	err = db2.InsertRows("products", rows2)
	require.NoError(t, err, "Failed to insert into db2")

	// Verify that dropping schema in one database doesn't affect the other
	err = db1.DropSchema("users")
	assert.NoError(t, err, "Failed to drop schema in db1")

	// db2 should still have its table - we can insert more data to verify
	moreRows := []data.Row{
		{"id": 3, "description": "Widget C", "value": 39.99},
	}
	err = db2.InsertRows("products", moreRows)
	assert.NoError(t, err, "db2 should still be functional after db1 schema drop")
}

func TestNamedSQLiteGremelDB_SameName_SeparateInstances(t *testing.T) {
	t.Skip("Skipping test - requires shared database file handling")
	// Create two databases with the same name
	// They should still be separate instances but share the same logical database
	db1 := newNamedSQLiteGremelDB("shared_db")
	db2 := newNamedSQLiteGremelDB("shared_db")

	defer func() {
		db1.Close()
		db2.Close()
	}()

	// Create schema in db1
	sampleRow := data.Row{"id": 1, "name": "test"}
	err := db1.CreateSchema("shared_table", sampleRow)
	require.NoError(t, err, "Failed to create schema in db1")

	// Insert data via db1
	rows := []data.Row{
		{"id": 1, "name": "Record from db1"},
	}
	err = db1.InsertRows("shared_table", rows)
	require.NoError(t, err, "Failed to insert via db1")

	// Insert more data via db2 (should work since they share the same database)
	moreRows := []data.Row{
		{"id": 2, "name": "Record from db2"},
	}
	err = db2.InsertRows("shared_table", moreRows)
	assert.NoError(t, err, "Should be able to insert via db2 into shared database")
}

func TestNamedSQLiteGremelDB_DifferentNames_CompleteSeparation(t *testing.T) {
	// Create databases with different names - they should be completely separate
	dbA := newNamedSQLiteGremelDB("database_a")
	dbB := newNamedSQLiteGremelDB("database_b")

	defer func() {
		dbA.Close()
		dbB.Close()
	}()

	// Create same table name in both databases with same structure
	sampleRow := data.Row{"id": 1, "data": "test"}

	err := dbA.CreateSchema("common_table", sampleRow)
	require.NoError(t, err, "Failed to create schema in dbA")

	err = dbB.CreateSchema("common_table", sampleRow)
	require.NoError(t, err, "Failed to create schema in dbB")

	// Insert different data into each
	rowsA := []data.Row{
		{"id": 1, "data": "From Database A"},
	}
	err = dbA.InsertRows("common_table", rowsA)
	require.NoError(t, err, "Failed to insert into dbA")

	rowsB := []data.Row{
		{"id": 1, "data": "From Database B"},
	}
	err = dbB.InsertRows("common_table", rowsB)
	require.NoError(t, err, "Failed to insert into dbB")

	// Drop table in dbA
	err = dbA.DropSchema("common_table")
	assert.NoError(t, err, "Failed to drop schema in dbA")

	// dbB should still be able to insert data (proving they're separate)
	moreRowsB := []data.Row{
		{"id": 2, "data": "Still working in Database B"},
	}
	err = dbB.InsertRows("common_table", moreRowsB)
	assert.NoError(t, err, "dbB should still work after dbA table was dropped")
}
