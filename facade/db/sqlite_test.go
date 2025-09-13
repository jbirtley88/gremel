package db

import (
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSQLiteGremelDB(t *testing.T) {
	t.Run("creates valid SQLite database connection", func(t *testing.T) {
		db := newSQLiteGremelDB()
		require.NotNil(t, db)

		// Should be able to cast to SQLiteGremelDB
		sqliteDB, ok := db.(*SQLiteGremelDB)
		assert.True(t, ok)
		assert.NotNil(t, sqliteDB.db)

		// Clean up
		err := db.Close()
		assert.NoError(t, err)
	})

	t.Run("database connection is working", func(t *testing.T) {
		db := newSQLiteGremelDB()
		require.NotNil(t, db)

		// Cast to access the underlying db for testing
		sqliteDB := db.(*SQLiteGremelDB)
		err := sqliteDB.db.Ping()
		assert.NoError(t, err)

		// Clean up
		err = db.Close()
		assert.NoError(t, err)
	})
}

func TestSQLiteGremelDB_Close(t *testing.T) {
	t.Run("closes database connection successfully", func(t *testing.T) {
		db := newSQLiteGremelDB()
		require.NotNil(t, db)

		err := db.Close()
		assert.NoError(t, err)
	})

	t.Run("closing already closed connection", func(t *testing.T) {
		db := newSQLiteGremelDB()
		require.NotNil(t, db)

		// Close once
		err := db.Close()
		assert.NoError(t, err)

		// Close again - SQLite handles this gracefully
		err = db.Close()
		assert.NoError(t, err)
	})
}

func TestSQLiteGremelDB_getColumnType(t *testing.T) {
	db := newSQLiteGremelDB().(*SQLiteGremelDB)
	defer db.Close()

	tests := []struct {
		name        string
		value       any
		expected    string
		expectError bool
	}{
		{
			name:     "int type",
			value:    42,
			expected: "INTEGER",
		},
		{
			name:     "int32 type",
			value:    int32(42),
			expected: "INTEGER",
		},
		{
			name:     "int64 type",
			value:    int64(42),
			expected: "INTEGER",
		},
		{
			name:     "float32 type",
			value:    float32(3.14),
			expected: "REAL",
		},
		{
			name:     "float64 type",
			value:    3.14159,
			expected: "REAL",
		},
		{
			name:     "bool type true",
			value:    true,
			expected: "BOOLEAN",
		},
		{
			name:     "bool type false",
			value:    false,
			expected: "BOOLEAN",
		},
		{
			name:     "string type",
			value:    "hello world",
			expected: "TEXT",
		},
		{
			name:        "unsupported type - complex",
			value:       complex(1, 2),
			expectError: true,
		},
		{
			name:        "unsupported type - slice",
			value:       []string{"a", "b"},
			expectError: true,
		},
		{
			name:        "unsupported type - map",
			value:       map[string]string{"key": "value"},
			expectError: true,
		},
		{
			name:        "unsupported type - nil",
			value:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := db.getColumnType(tt.value)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported data type")
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSQLiteGremelDB_getCreateTableSQL(t *testing.T) {
	db := newSQLiteGremelDB().(*SQLiteGremelDB)
	defer db.Close()

	t.Run("generates SQL for valid row", func(t *testing.T) {
		row := data.Row{
			"id":     1,
			"name":   "John Doe",
			"age":    30,
			"salary": 50000.50,
			"active": true,
		}

		sql, err := db.getCreateTableSQL("users", row)
		assert.NoError(t, err)
		assert.NotEmpty(t, sql)

		// Check that SQL contains CREATE TABLE statement
		assert.Contains(t, sql, "CREATE TABLE users (")

		// Check that all columns are present with correct types
		assert.Contains(t, sql, "id INTEGER")
		assert.Contains(t, sql, "name TEXT")
		assert.Contains(t, sql, "age INTEGER")
		assert.Contains(t, sql, "salary REAL")
		assert.Contains(t, sql, "active BOOLEAN")

		// Check that SQL ends properly
		assert.Contains(t, sql, ");")
	})

	t.Run("generates SQL for empty row", func(t *testing.T) {
		row := data.Row{}

		sql, err := db.getCreateTableSQL("empty_table", row)
		assert.NoError(t, err)
		assert.Contains(t, sql, "CREATE TABLE empty_table (")
		assert.Contains(t, sql, "_placeholder INTEGER") // Empty tables get a placeholder column
		assert.Contains(t, sql, ");")
	})

	t.Run("handles row with unsupported data type", func(t *testing.T) {
		row := data.Row{
			"id":          1,
			"unsupported": complex(1, 2),
		}

		sql, err := db.getCreateTableSQL("test_table", row)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get column type for field \"unsupported\"")
		assert.Contains(t, err.Error(), "unsupported data type")
		assert.Empty(t, sql)
	})

	t.Run("generates SQL with special characters in table name", func(t *testing.T) {
		row := data.Row{
			"field1": "value1",
		}

		sql, err := db.getCreateTableSQL("table_with_underscores", row)
		assert.NoError(t, err)
		assert.Contains(t, sql, "CREATE TABLE table_with_underscores (")
	})
}

func TestSQLiteGremelDB_CreateSchema(t *testing.T) {
	t.Run("creates schema successfully", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		row := data.Row{
			"id":    1,
			"name":  "John Doe",
			"email": "jb@example.com",
		}

		err := db.CreateSchema("users", row)
		assert.NoError(t, err)

		// Verify table was created by querying schema
		var tableName string
		err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
		assert.NoError(t, err)
		assert.Equal(t, "users", tableName)
	})

	t.Run("creates schema with different data types", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		row := data.Row{
			"int_field":    42,
			"string_field": "test",
			"float_field":  3.14,
			"bool_field":   true,
		}

		err := db.CreateSchema("mixed_types", row)
		assert.NoError(t, err)

		// Verify table was created
		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='mixed_types'").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("fails with unsupported data type", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		row := data.Row{
			"id":          1,
			"unsupported": complex(1, 2),
		}

		err := db.CreateSchema("test_table", row)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CreateSchema(test_table): failed to generate CREATE TABLE SQL")
		assert.Contains(t, err.Error(), "unsupported data type")
	})

	t.Run("handles table creation with same name twice", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		row := data.Row{
			"id": 1,
		}

		// Create table first time
		err := db.CreateSchema("duplicate_table", row)
		assert.NoError(t, err)

		// Try to create the same table again - should fail
		err = db.CreateSchema("duplicate_table", row)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CreateSchema(duplicate_table): failed to create table")
	})

	t.Run("creates schema with empty row", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		row := data.Row{}

		err := db.CreateSchema("empty_table", row)
		assert.NoError(t, err)

		// Verify table was created
		var tableName string
		err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='empty_table'").Scan(&tableName)
		assert.NoError(t, err)
		assert.Equal(t, "empty_table", tableName)
	})
}

func TestSQLiteGremelDB_DropSchema(t *testing.T) {
	t.Run("drops existing table successfully", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		// First create a table
		row := data.Row{
			"id":   1,
			"name": "test",
		}
		err := db.CreateSchema("test_table", row)
		require.NoError(t, err)

		// Verify table exists
		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// Drop the table
		err = db.DropSchema("test_table")
		assert.NoError(t, err)

		// Verify table no longer exists
		err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("drops non-existing table successfully", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		// Try to drop a table that doesn't exist - should not error due to IF EXISTS
		err := db.DropSchema("non_existing_table")
		assert.NoError(t, err)
	})

	t.Run("drops multiple tables with same name pattern", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)
		defer db.Close()

		// Create multiple tables
		for i := 1; i <= 3; i++ {
			row := data.Row{
				"id": i,
			}
			tableName := "test_table"
			if i > 1 {
				// Only create one table since we can't have duplicates
				continue
			}
			err := db.CreateSchema(tableName, row)
			require.NoError(t, err)
		}

		// Drop the table
		err := db.DropSchema("test_table")
		assert.NoError(t, err)

		// Verify table no longer exists
		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestSQLiteGremelDB_Integration(t *testing.T) {
	t.Run("complete lifecycle - create, use, drop", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)

		// Define test data
		tableName := "integration_test"
		row := data.Row{
			"id":       1,
			"username": "testuser",
			"balance":  100.50,
			"active":   true,
		}

		// Create schema
		err := db.CreateSchema(tableName, row)
		require.NoError(t, err)

		// Insert test data
		insertSQL := "INSERT INTO integration_test (id, username, balance, active) VALUES (?, ?, ?, ?)"
		_, err = db.db.Exec(insertSQL, 1, "testuser", 100.50, true)
		assert.NoError(t, err)

		// Query the data back
		var id int
		var username string
		var balance float64
		var active bool

		selectSQL := "SELECT id, username, balance, active FROM integration_test WHERE id = ?"
		err = db.db.QueryRow(selectSQL, 1).Scan(&id, &username, &balance, &active)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, 100.50, balance)
		assert.True(t, active)

		// Drop schema
		err = db.DropSchema(tableName)
		assert.NoError(t, err)

		// Verify table is gone
		var count int
		err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("multiple tables management", func(t *testing.T) {
		db := newSQLiteGremelDB().(*SQLiteGremelDB)

		// Create multiple tables
		tables := map[string]data.Row{
			"users": {
				"id":   1,
				"name": "John",
			},
			"products": {
				"id":    1,
				"title": "Widget",
				"price": 9.99,
			},
			"orders": {
				"id":      1,
				"user_id": 1,
				"total":   29.97,
			},
		}

		// Create all schemas
		for tableName, row := range tables {
			err := db.CreateSchema(tableName, row)
			require.NoError(t, err)
		}

		// Verify all tables exist
		for tableName := range tables {
			var count int
			err := db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "Table %s should exist", tableName)
		}

		// Drop all schemas
		for tableName := range tables {
			err := db.DropSchema(tableName)
			assert.NoError(t, err)
		}

		// Verify all tables are gone
		var totalCount int
		err := db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&totalCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, totalCount)
	})
}
