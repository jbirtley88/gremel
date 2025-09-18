package helper

import (
	"reflect"
	"testing"

	"github.com/jbirtley88/gremel/data"
)

func TestDeriveSchema(t *testing.T) {
	// Test: No rows provided
	_, err := DeriveSchema([]data.Row{})
	if err == nil {
		t.Error("expected error for no rows, got nil")
	}

	// Test: Single row, simple types
	row := data.Row{"a": int64(1), "b": 2.5, "c": "foo", "d": true}
	schema, err := DeriveSchema([]data.Row{row})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	want := map[string]reflect.Kind{
		"a": reflect.Int64,
		"b": reflect.Float64,
		"c": reflect.String,
		"d": reflect.Bool,
	}
	for k, v := range want {
		if schema[k] != v {
			t.Errorf("expected %s to be %v, got %v", k, v, schema[k])
		}
	}

	// Test: Type promotion int -> float
	rows := []data.Row{
		{"x": int64(1)},
		{"y": 2},   // a field not seen before
		{"y": 2.5}, // upgrade to float
	}
	schema, err = DeriveSchema(rows)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if schema["x"] != reflect.Int64 {
		t.Errorf("expected x to be Int64, got %v", schema["x"])
	}
	if schema["y"] != reflect.Float64 {
		t.Errorf("expected y to be Float64, got %v", schema["y"])
	}

	// Test: Type conflict promotes to string
	rows = []data.Row{
		{"y": int64(1)},
		{"y": "bar"},
	}
	schema, err = DeriveSchema(rows)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if schema["y"] != reflect.String {
		t.Errorf("expected y to be String, got %v", schema["y"])
	}

	// Test: Multiple columns, mixed types
	rows = []data.Row{
		{"a": int64(1), "b": "foo"},
		{"a": 2.5, "b": "bar"},
		{"a": "baz", "b": "qux"},
	}
	schema, err = DeriveSchema(rows)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if schema["a"] != reflect.String {
		t.Errorf("expected a to be String, got %v", schema["a"])
	}
	if schema["b"] != reflect.String {
		t.Errorf("expected b to be String, got %v", schema["b"])
	}
}
