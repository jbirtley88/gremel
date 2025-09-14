package db

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/jbirtley88/gremel/data"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteGremelDB struct {
	db           *sql.DB
	schemaByName map[string]data.Row
}

// NewSQLiteGremelDB creates a new in-memory SQLite database connection
func newSQLiteGremelDB() GremelDB {
	return newNamedSQLiteGremelDB("gremel")
}

// NewNamedSQLiteGremelDB creates a new named in-memory SQLite database connection
// Each named database is completely separate from others
func newNamedSQLiteGremelDB(dbName string) GremelDB {
	// Connect to named in-memory SQLite database
	// Using file::memory: syntax with cache=shared ensures each named database is separate
	// dbPath := "/var/tmp/" + dbName + ".db"
	// _ = os.Remove(dbPath)
	// _ = os.WriteFile(dbPath, []byte{}, 0644)
	// connectionString := fmt.Sprintf("file:%s?cache=shared", dbPath)
	connectionString := fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName)
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return NewErrorGremelDB(fmt.Errorf("failed to open in-memory SQLite database %q: %w", dbName, err))
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return NewErrorGremelDB(fmt.Errorf("failed to ping SQLite database %q: %w", dbName, err))
	}

	return &SQLiteGremelDB{
		db:           db,
		schemaByName: make(map[string]data.Row),
	}
}

// Close closes the database connection
func (db *SQLiteGremelDB) Close() error {
	return db.db.Close()
}

func (db *SQLiteGremelDB) getColumnType(value any) (string, error) {
	switch v := value.(type) {
	case int, int32, int64:
		return "INTEGER", nil
	case float32:
		if _, frac := math.Modf(float64(v)); frac == 0 {
			return "INTEGER", nil
		}
		return "REAL", nil
	case float64:
		if _, frac := math.Modf(v); frac == 0 {
			return "INTEGER", nil
		}
		return "REAL", nil
	case bool:
		return "BOOLEAN", nil
	case string:
		return "TEXT", nil
	default:
		return "", fmt.Errorf("unsupported data type: %T", value)
	}
}

func (db *SQLiteGremelDB) getCreateTableSQL(tableName string, row data.Row) (string, data.Row, error) {
	sqlLines := make([]string, 0)

	// Handle empty row case
	sqlLines = append(sqlLines, fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName))
	if len(row) == 0 {
		sqlLines = append(sqlLines, fmt.Sprintf("CREATE TABLE %s (", tableName))
		sqlLines = append(sqlLines, "    _placeholder INTEGER") // Add a placeholder column for empty tables
		sqlLines = append(sqlLines, ");")
		return strings.Join(sqlLines, "\n"), nil, nil
	}

	sqlLines = append(sqlLines, fmt.Sprintf("CREATE TABLE %s (", tableName))

	// Collect column definitions first
	schema := make(data.Row)
	columns := make([]string, 0, len(row))
	for fieldName, fieldValue := range row {
		columnType, err := db.getColumnType(fieldValue)
		schema[fieldName] = columnType
		if err != nil {
			return "", nil, fmt.Errorf("failed to get column type for field %q: %w", fieldName, err)
		}
		columns = append(columns, fmt.Sprintf("    %s %s", fieldName, columnType))
	}

	// Join columns with commas and add to SQL lines
	for i := range columns {
		if i < len(columns)-1 {
			columns[i] += ","
		}
		sqlLines = append(sqlLines, columns[i])
	}

	sqlLines = append(sqlLines, ");")
	return strings.Join(sqlLines, "\n"), schema, nil
}

func (db *SQLiteGremelDB) CreateSchema(tableName string, row data.Row) error {
	// Create the accounts table
	createTableSQL, schema, err := db.getCreateTableSQL(tableName, row)
	if err != nil {
		return fmt.Errorf("CreateSchema(%s): failed to generate CREATE TABLE SQL: %w", tableName, err)
	}

	// TODO(john): debug logger
	log.Printf("Creating table %s with SQL:\n%s\n", tableName, createTableSQL)
	_, err = db.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("CreateSchema(%s): failed to create table: %w", tableName, err)
	}

	// Stash the schema for later retrieval
	db.schemaByName[tableName] = schema

	// TODO(jb): Create indexes for better query performance

	return nil
}

func (db *SQLiteGremelDB) GetSchema(tableName string) (data.Row, error) {
	schema, exists := db.schemaByName[tableName]
	if !exists {
		return nil, fmt.Errorf("GetSchema(%s): schema not found", tableName)
	}
	return schema, nil
}

func (db *SQLiteGremelDB) DropSchema(tableName string) error {
	// Drop views first (in reverse order of dependency)
	dropStatements := []string{
		fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName),
	}

	for _, dropSQL := range dropStatements {
		if _, err := db.db.Exec(dropSQL); err != nil {
			return fmt.Errorf("DropSchema(%s): failed to drop schema: %w", tableName, err)
		}
	}

	return nil
}

func (db *SQLiteGremelDB) GetTables() ([]string, error) {
	var tables []string
	for t := range db.schemaByName {
		tables = append(tables, t)
	}
	return tables, nil
}

func (db *SQLiteGremelDB) InsertRows(tableName string, rows []data.Row) error {
	if len(rows) == 0 {
		return nil // Nothing to insert
	}

	// Prepare the insert statement
	fieldNames := make([]string, 0, len(rows[0]))
	placeholders := make([]string, 0, len(rows[0]))
	for fieldName := range rows[0] {
		fieldNames = append(fieldNames, fieldName)
		placeholders = append(placeholders, "?")
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := db.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("InsertRows(%s): failed to prepare insert statement: %w", tableName, err)
	}
	defer stmt.Close()

	// Insert each row
	for _, row := range rows {
		values := make([]any, 0, len(fieldNames))
		for _, fieldName := range fieldNames {
			values = append(values, row[fieldName])
		}

		if _, err := stmt.Exec(values...); err != nil {
			return fmt.Errorf("InsertRows(%s): failed to execute insert statement: %w", tableName, err)
		}
	}

	return nil
}

func (db *SQLiteGremelDB) Query(sqlQuery string) ([]data.Row, error) {
	rows, err := db.db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("Query(%s): failed to execute query: %w", sqlQuery, err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Query(%s): failed to get columns: %w", sqlQuery, err)
	}

	var results []data.Row

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("Query(%s): failed to scan row: %w", sqlQuery, err)
		}

		// Create a map and populate it with the column data
		rowMap := make(data.Row)
		for i, colName := range columns {
			val := columnValues[i]

			// Convert []byte to string for TEXT columns
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Query(%s): row iteration error: %w", sqlQuery, err)
	}

	return results, nil
}
