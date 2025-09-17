package data

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

func GetType(v any) reflect.Kind {
	rv := reflect.ValueOf(v)
	return rv.Kind()
}

var reFloat = regexp.MustCompile(`^[+-]?(?:(?:\d+\.?\d*)|(?:\.\d+))(?:[eE][+-]?\d+)?$|^[+-]?(?:Inf|inf|NaN|nan)$`)

// IsFloat determines if a string is a valid float64 using regex
func IsFloat(s string) bool {
	// Regex pattern for valid float64 values:
	// Optional +/- sign, followed by digits, optional decimal point and more digits,
	// optional scientific notation (e/E followed by optional +/- and digits)
	// Also handles special cases like Inf, -Inf, NaN
	return reFloat.MatchString(s)
}

func InferValue(value any) any {
	// Try int
	sValue := fmt.Sprint(value)
	if v, err := strconv.ParseInt(sValue, 10, 64); err == nil {
		return v
	}
	// Try float
	if v, err := strconv.ParseFloat(sValue, 64); err == nil {
		return v
	}
	// Try bool
	if v, err := strconv.ParseBool(sValue); err == nil {
		return v
	}
	// Default to string
	return sValue
}
