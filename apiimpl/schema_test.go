package apiimpl

import (
	"context"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/stretchr/testify/assert"
)

func TestGetSchema_ValidTable(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_schema_table"
	sampleRow := data.Row{
		"id":       123,
		"username": "testuser",
		"email":    "test@example.com",
		"active":   true,
	}

	// Clean up any existing table first
	_ = database.DropSchema(tableName)

	// Create a test table with schema
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Test GetSchema
	schema, err := GetSchema(ctx, tableName)

	// Should succeed
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Schema should match the sample row structure
	// The schema contains the field names and their inferred types
	for key := range sampleRow {
		assert.Contains(t, schema, key, "Schema should contain field %s", key)
	}

	// Clean up
	_ = database.DropSchema(tableName)
}

func TestGetSchema_NonExistentTable(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())

	nonExistentTable := "table_that_does_not_exist_12345"

	// Test GetSchema for non-existent table
	schema, err := GetSchema(ctx, nonExistentTable)

	// Should return an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schema not found", "Error should indicate schema not found")

	// Schema should be empty or nil
	if schema != nil {
		assert.Equal(t, 0, len(schema), "Schema should be empty for non-existent table")
	}
}

func TestGetSchema_EmptyTableName(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())

	// Test GetSchema with empty table name
	schema, err := GetSchema(ctx, "")

	// Should return an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schema not found", "Error should indicate schema not found")

	// Schema should be empty
	if schema != nil {
		assert.Equal(t, 0, len(schema), "Schema should be empty for empty table name")
	}
}

func TestGetSchema_MultipleFieldTypes(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_multi_type_schema"
	sampleRow := data.Row{
		"string_field": "test string",
		"int_field":    42,
		"float_field":  3.14159,
		"bool_field":   true,
		// Note: nil fields are not supported in CreateSchema as they can't be typed
	}

	// Clean up any existing table first
	_ = database.DropSchema(tableName)

	// Create a test table with various field types
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Test GetSchema
	schema, err := GetSchema(ctx, tableName)

	// Should succeed
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Verify all fields are present in schema
	expectedFields := []string{"string_field", "int_field", "float_field", "bool_field"}
	for _, field := range expectedFields {
		assert.Contains(t, schema, field, "Schema should contain field %s", field)
	}

	// Clean up
	_ = database.DropSchema(tableName)
}

func TestGetSchema_SchemaContent(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_schema_content"
	sampleRow := data.Row{
		"id":     123,
		"name":   "test",
		"score":  95.5,
		"active": true,
	}

	// Clean up and create table
	_ = database.DropSchema(tableName)
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Test GetSchema
	schema, err := GetSchema(ctx, tableName)

	// Should succeed
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// The schema should contain field names mapped to SQL types
	assert.Greater(t, len(schema), 0, "Schema should not be empty")

	// Check that the schema contains the expected field names
	assert.Contains(t, schema, "id", "Schema should contain id field")
	assert.Contains(t, schema, "name", "Schema should contain name field")
	assert.Contains(t, schema, "score", "Schema should contain score field")
	assert.Contains(t, schema, "active", "Schema should contain active field")

	// Check that the values are SQL type strings
	if idType, exists := schema["id"]; exists {
		assert.IsType(t, "", idType, "Schema field type should be string")
		assert.Contains(t, []string{"INTEGER", "INT"}, idType, "id field should be INTEGER type")
	}

	if nameType, exists := schema["name"]; exists {
		assert.IsType(t, "", nameType, "Schema field type should be string")
		assert.Equal(t, "TEXT", nameType, "name field should be TEXT type")
	}

	if scoreType, exists := schema["score"]; exists {
		assert.IsType(t, "", scoreType, "Schema field type should be string")
		assert.Equal(t, "REAL", scoreType, "score field should be REAL type")
	}

	if activeType, exists := schema["active"]; exists {
		assert.IsType(t, "", activeType, "Schema field type should be string")
		assert.Equal(t, "BOOLEAN", activeType, "active field should be BOOLEAN type")
	}

	// Clean up
	_ = database.DropSchema(tableName)
}

func TestGetSchema_ContextParameter(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_context_schema"
	sampleRow := data.Row{"field": "value"}

	// Clean up and create table
	_ = database.DropSchema(tableName)
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Test with nil context (though not recommended in practice)
	schema, err := GetSchema(nil, tableName)
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Test with proper context
	schema, err = GetSchema(ctx, tableName)
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	// Clean up
	_ = database.DropSchema(tableName)
}

func TestGetSchema_ReturnType(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_return_type_schema"
	sampleRow := data.Row{"test_field": "test_value"}

	// Clean up and create table
	_ = database.DropSchema(tableName)
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Test return type
	schema, err := GetSchema(ctx, tableName)

	// Should return data.Row and no error
	assert.NoError(t, err)
	assert.IsType(t, data.Row{}, schema)
	assert.NotNil(t, schema)

	// Clean up
	_ = database.DropSchema(tableName)
}

func TestGetSchema_ConsistencyAcrossMultipleCalls(t *testing.T) {
	// Create a context
	ctx := data.NewGremelContext(context.Background())
	database := db.GetGremelDB()

	tableName := "test_consistency_schema"
	sampleRow := data.Row{
		"field1": "value1",
		"field2": 42,
		"field3": true,
	}

	// Clean up and create table
	_ = database.DropSchema(tableName)
	err := database.CreateSchema(tableName, []data.Row{sampleRow})
	assert.NoError(t, err)

	// Call GetSchema multiple times
	schema1, err1 := GetSchema(ctx, tableName)
	assert.NoError(t, err1)

	schema2, err2 := GetSchema(ctx, tableName)
	assert.NoError(t, err2)

	schema3, err3 := GetSchema(ctx, tableName)
	assert.NoError(t, err3)

	// All calls should return identical results
	assert.Equal(t, schema1, schema2, "Schema should be consistent across multiple calls")
	assert.Equal(t, schema2, schema3, "Schema should be consistent across multiple calls")
	assert.Equal(t, schema1, schema3, "Schema should be consistent across multiple calls")

	// Clean up
	_ = database.DropSchema(tableName)
}
