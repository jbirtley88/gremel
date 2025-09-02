package data

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockParser is a test implementation of the Parser interface
type mockParser struct {
	*BaseParser
	returnError bool
}

func newMockParser(name string, ctx context.Context) Parser {
	return &mockParser{
		BaseParser: NewBaseParser(ctx, name),
	}
}

func newMockParserWithError(name string, ctx context.Context) Parser {
	return &mockParser{
		BaseParser:  NewBaseParser(ctx, name),
		returnError: true,
	}
}

func TestGetParserRegistry(t *testing.T) {
	registry := GetParserRegistry()

	require.NotNil(t, registry)
	assert.IsType(t, &parserRegistry{}, registry)
}

func TestGetParserRegistry_ReturnsSameInstance(t *testing.T) {
	registry1 := GetParserRegistry()
	registry2 := GetParserRegistry()

	assert.Same(t, registry1, registry2, "GetParserRegistry should return the same instance")
}

func TestParserRegistry_Register(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	err := registry.Register("test", newMockParser)

	assert.NoError(t, err)
	assert.Contains(t, registry.parsersByName, "test")
}

func TestParserRegistry_Register_CaseInsensitive(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	err := registry.Register("TestParser", newMockParser)

	assert.NoError(t, err)
	assert.Contains(t, registry.parsersByName, "testparser")
	assert.NotContains(t, registry.parsersByName, "TestParser")
}

func TestParserRegistry_Register_OverwriteExisting(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	// Register first parser
	err1 := registry.Register("test", newMockParser)
	assert.NoError(t, err1)

	// Register second parser with same name
	err2 := registry.Register("test", newMockParserWithError)
	assert.NoError(t, err2)

	// Get parser and verify it's the second one
	ctx := context.Background()
	parser := registry.Get("test", ctx)
	mockParser, ok := parser.(*mockParser)
	require.True(t, ok)
	assert.True(t, mockParser.returnError, "Should be the second parser with returnError=true")
}

func TestParserRegistry_Get_ExistingParser(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	err := registry.Register("test", newMockParser)
	require.NoError(t, err)

	ctx := context.Background()
	parser := registry.Get("test", ctx)

	require.NotNil(t, parser)
	assert.Equal(t, "test", parser.GetName())

	// Verify it's not an error parser
	_, isErrorParser := parser.(*ParserError)
	assert.False(t, isErrorParser, "Should not return a ParserError for existing parser")
}

func TestParserRegistry_Get_CaseInsensitive(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	err := registry.Register("TestParser", newMockParser)
	require.NoError(t, err)

	ctx := context.Background()

	tests := []string{"testparser", "TESTPARSER", "TestParser", "tEsTpArSeR"}
	for _, testName := range tests {
		t.Run("name_"+testName, func(t *testing.T) {
			parser := registry.Get(testName, ctx)

			require.NotNil(t, parser)
			assert.Equal(t, testName, parser.GetName())

			// Verify it's not an error parser
			_, isErrorParser := parser.(*ParserError)
			assert.False(t, isErrorParser, "Should not return a ParserError for existing parser")
		})
	}
}

func TestParserRegistry_Get_NonExistentParser(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	ctx := context.Background()
	parser := registry.Get("nonexistent", ctx)

	require.NotNil(t, parser)

	// Verify it returns a ParserError
	errorParser, isErrorParser := parser.(*ParserError)
	require.True(t, isErrorParser, "Should return a ParserError for non-existent parser")

	assert.NotNil(t, errorParser.Err)
	assert.Contains(t, errorParser.Err.Error(), "Parser 'nonexistent' not found")
}

func TestParserRegistry_Get_EmptyName(t *testing.T) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	ctx := context.Background()
	parser := registry.Get("", ctx)

	require.NotNil(t, parser)

	// Verify it returns a ParserError
	errorParser, isErrorParser := parser.(*ParserError)
	require.True(t, isErrorParser, "Should return a ParserError for empty name")

	assert.NotNil(t, errorParser.Err)
	assert.Contains(t, errorParser.Err.Error(), "Parser '' not found")
}

func TestParserRegistry_Integration(t *testing.T) {
	// Test the global registry instance
	ctx := context.Background()

	// Register a parser
	err := GetParserRegistry().Register("integration-test", newMockParser)
	require.NoError(t, err)

	// Get the parser
	parser := GetParserRegistry().Get("integration-test", ctx)
	require.NotNil(t, parser)

	assert.Equal(t, "integration-test", parser.GetName())

	// Verify it's not an error parser
	_, isErrorParser := parser.(*ParserError)
	assert.False(t, isErrorParser, "Should not return a ParserError for registered parser")
}

func TestParserRegistry_Integration_CaseInsensitive(t *testing.T) {
	// Test the global registry instance with case insensitive names
	ctx := context.Background()

	// Register with mixed case
	err := GetParserRegistry().Register("Integration-Test-Case", newMockParser)
	require.NoError(t, err)

	// Get with different case
	parser := GetParserRegistry().Get("integration-test-case", ctx)
	require.NotNil(t, parser)

	assert.Equal(t, "integration-test-case", parser.GetName())

	// Verify it's not an error parser
	_, isErrorParser := parser.(*ParserError)
	assert.False(t, isErrorParser, "Should not return a ParserError for registered parser")
}

func TestParserConstructor_FunctionSignature(t *testing.T) {
	// Test that our mock constructor matches the expected signature
	var constructor ParserConstructor = newMockParser

	ctx := context.Background()
	parser := constructor("test", ctx)

	require.NotNil(t, parser)
	assert.Equal(t, "test", parser.GetName())
}

// Benchmark tests to verify performance
func BenchmarkParserRegistry_Register(b *testing.B) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Register("benchmark", newMockParser)
	}
}

func BenchmarkParserRegistry_Get_Found(b *testing.B) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}
	registry.Register("benchmark", newMockParser)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get("benchmark", ctx)
	}
}

func BenchmarkParserRegistry_Get_NotFound(b *testing.B) {
	registry := &parserRegistry{
		parsersByName: make(map[string]ParserConstructor),
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get("nonexistent", ctx)
	}
}
