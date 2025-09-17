package data

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestInferValue(t *testing.T) {
	tests := []struct {
		input       any
		expected    any
		isNaN       bool
		description string
	}{
		{"123", int64(123), false, "positive integer"},
		{"-456", int64(-456), false, "negative integer"},
		{"78.9", float64(78.9), false, "positive float"},
		{"-0.12", float64(-0.12), false, "negative float"},
		{"3.14e2", float64(314), false, "scientific notation"},
		{"-2.5E-1", float64(-0.25), false, "negative scientific notation"},
		{"NaN", math.NaN(), true, "NaN should parse as float64"},
		{"hello", "hello", false, "string"},
		{"123abc", "123abc", false, "mixed string"},
		{123, int64(123), false, "integer input"},
		{45.67, float64(45.67), false, "float input"},
	}

	for _, tt := range tests {
		result := InferValue(tt.input)

		if tt.isNaN {
			// Special case for NaN values
			resultFloat, ok := result.(float64)
			if !ok {
				t.Errorf("InferValue(%v) = %v (type %T), expected float64 NaN", tt.input, result, result)
				continue
			}
			if !math.IsNaN(resultFloat) {
				t.Errorf("InferValue(%v) = %v, expected NaN", tt.input, result)
			}
		} else {
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InferValue(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestIsFloat(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		// Valid floats
		{"123.456", true},
		{"123", true},
		{".456", true},
		{"123.", true},
		{"+123.456", true},
		{"-123.456", true},
		{"1.23e10", true},
		{"1.23E10", true},
		{"1.23e-10", true},
		{"1.23E+10", true},
		{"0", true},
		{"0.0", true},
		{"-0", true},
		{"+0", true},

		// Special values
		{"Inf", true},
		{"-Inf", true},
		{"+Inf", true},
		{"inf", true},
		{"-inf", true},
		{"+inf", true},
		{"NaN", true},
		{"nan", true},

		// Invalid floats
		{"", false},
		{"abc", false},
		{"12.34.56", false},
		{"12e", false},
		{"12e+", false},
		{"12ee10", false},
		{".e10", false},
		{"e10", false},
		{"12..34", false},
		{"++12", false},
		{"--12", false},
		{"12-34", false},
		{"12+34", false},
		{" 123 ", false}, // spaces not allowed
		{"∞", false},     // Unicode infinity not allowed
	}

	for _, tc := range testCases {
		result := IsFloat(tc.input)
		if result != tc.expected {
			t.Errorf("IsFloat(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestParseJsonInferTypes(t *testing.T) {
	rawJson := `[
	{
		"id": 101,
		"name": "Alice",
		"balance": 2500.0
	},
	{
		"id": 102,
		"name": "Bob",
		"balance": 2501.99
	}
	]`
	// Step 1: Parse the JSON
	var result []map[string]any
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(rawJson)))
	err := decoder.Decode(&result)
	if err != nil {
		t.Fatalf("ParseJSON error: %v", err)
	}
	for i := range result {
		for k, v := range result[i] {
			result[i][k] = InferValue(v)
		}
	}
	// Step 2: Verify the types
	if _, ok := result[0]["id"].(int64); !ok {
		t.Errorf("Expected id to be int64, got %T", result[0]["id"])
	}
	if _, ok := result[0]["name"].(string); !ok {
		t.Errorf("Expected name to be string, got %T", result[0]["name"])
	}
	if _, ok := result[0]["balance"].(int64); !ok {
		t.Errorf("Expected balance to be int64, got %T", result[0]["balance"])
	}
	if _, ok := result[1]["id"].(int64); !ok {
		t.Errorf("Expected id to be int64, got %T", result[1]["id"])
	}
	if _, ok := result[1]["balance"].(float64); !ok {
		t.Errorf("Expected balance to be float64, got %T", result[1]["balance"])
	}
}
