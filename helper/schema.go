package helper

import (
	"fmt"
	"math"
	"reflect"

	"github.com/jbirtley88/gremel/data"
)

func DeriveSchema(rows []data.Row) (map[string]reflect.Kind, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("deriveSchema: no rows provided")
	}

	schema := make(map[string]reflect.Kind)
	for _, row := range rows {
		for fieldName, fieldValue := range row {
			columnType, err := GetColumnType(fieldValue)
			if err != nil {
				return nil, fmt.Errorf("DeriveSchema: failed to get column type for field %q: %w", fieldName, err)
			}

			// There are a couple of rules here:
			// 1. If the column is already in the schema, and the type is the same, do nothing
			// 2. If the column is already in the schema as a Float64, and the new type is Int64, do nothing
			// 3. If the column is already in the schema as an Int64, and the new type is Float64, promote to Float64
			// 4. If the column is already in the schema as something else, and the new type is different, promote to String
			if existingType, exists := schema[fieldName]; exists {
				if existingType == columnType {
					// Same type, do nothing
					continue
				}
				if existingType == reflect.Float64 && columnType == reflect.Int64 {
					// Float64 is more general than Int64, do nothing
					continue
				}
				if existingType == reflect.Int64 && columnType == reflect.Float64 {
					// Promote to Float64
					schema[fieldName] = reflect.Float64
					continue
				}
				// Different types, promote to String
				schema[fieldName] = reflect.String
				continue
			}
			// New column, add to schema
			schema[fieldName] = columnType
		}
	}
	return schema, nil
}

func GetColumnType(value any) (reflect.Kind, error) {
	switch v := value.(type) {
	case int, int32, int64:
		return reflect.Int64, nil
	case float32:
		if _, frac := math.Modf(float64(v)); frac == 0 {
			return reflect.Int64, nil
		}
		return reflect.Float64, nil
	case float64:
		if _, frac := math.Modf(v); frac == 0 {
			return reflect.Int64, nil
		}
		return reflect.Float64, nil
	case bool:
		return reflect.Bool, nil
	case string:
		return reflect.String, nil
	default:
		return reflect.Invalid, fmt.Errorf("unsupported data type: %T", value)
	}
}
