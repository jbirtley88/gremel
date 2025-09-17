package helper

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jbirtley88/gremel/data"
)

func NormaliseNumbers(rows []data.Row) []data.Row {
	// Normalise the types in each row
	// The default JSON parser, when decoding into a map[string]any, treats
	// all numbers as float64 (for speed and safety)
	// We want to detect integers vs floats, particularly the edge case of
	// numbers that are mathematically integers (e.g. 1234.0)
	//
	// We need to scan each column in each row and if any row value for that
	// column is a float that is not a whole number, we need to convert *all*
	// the row values in that column to float64

	// First pass: identify columns that contain non-whole-number floats
	floatColumns := make(map[string]bool)
	for i := range rows {
		for k, v := range rows[i] {
			if floatVal, isFloat := v.(float64); isFloat {
				if _, frac := math.Modf(floatVal); frac != 0 {
					floatColumns[k] = true
				}
			}
		}
	}

	// Second pass: normalize types, forcing float64 for columns identified above
	for i := range rows {
		for k, v := range rows[i] {
			if floatColumns[k] {
				// Force this column to be float64
				if f, ok := v.(float64); ok {
					rows[i][k] = f
				} else {
					// Convert other types to float64 if possible
					if parsed, err := strconv.ParseFloat(fmt.Sprint(v), 64); err == nil {
						rows[i][k] = parsed
					} else {
						// Keep as original value if can't convert to float
						rows[i][k] = v
					}
				}
			} else {
				// Use normal type inference for columns without non-whole floats
				// For float64 values that are whole numbers, convert to int64
				if f, ok := v.(float64); ok {
					if f == float64(int64(f)) && f >= float64(int64(-9223372036854775808)) && f <= float64(int64(9223372036854775807)) {
						rows[i][k] = int64(f)
					} else {
						rows[i][k] = f
					}
				} else {
					// For other types, use normal inference
					rows[i][k] = data.InferValue(v)
				}
			}
		}
	}
	return rows
}
