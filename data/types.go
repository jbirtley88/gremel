package data

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetType(v any) reflect.Kind {
	rv := reflect.ValueOf(v)
	return rv.Kind()
}

func ParseValue(value any) any {
	// Try int
	if v, err := strconv.ParseInt(fmt.Sprint(value), 10, 64); err == nil {
		return v
	}
	// Try float
	if v, err := strconv.ParseFloat(fmt.Sprint(value), 64); err == nil {
		return v
	}
	// Try bool
	if v, err := strconv.ParseBool(fmt.Sprint(value)); err == nil {
		return v
	}
	// Default to string
	return fmt.Sprint(value)
}
