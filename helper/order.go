package helper

import (
	"sort"
	"strings"
)

func Sort(rows []map[string]any, by string) ([]map[string]any, error) {
	if by == "" {
		// Nothing to do
		return rows, nil
	}

	descendingOrder := false
	fields := strings.Split(by, " ")
	if len(fields) > 1 && fields[1] == "desc" {
		descendingOrder = true
		by = fields[0]
	}
	sort.Slice(rows, func(i, j int) bool {
		// Assuming 'by' is a key in the maps, we compare the values
		// If the value is a string, we compare as strings
		// If the value is a number, we compare as float64
		valI, okI := rows[i][by]
		valJ, okJ := rows[j][by]
		if !okI || !okJ {
			// If either value is missing, we can't compare, so keep original order
			return false
		}
		switch vI := valI.(type) {
		case string:
			vJ, ok := valJ.(string)
			if !ok {
				return false // Can't compare different types
			}
			if descendingOrder {
				return vI > vJ
			}
			return vI < vJ
		case float64:
			vJ, ok := valJ.(float64)
			if !ok {
				return false // Can't compare different types
			}
			if descendingOrder {
				return vI > vJ
			}
			return vI < vJ
		case int:
			vJ, ok := valJ.(int)
			if !ok {
				return false // Can't compare different types
			}
			if descendingOrder {
				return vI > vJ
			}
			return vI < vJ
		default:
			// If the type is not supported, we can't compare, so keep original order
			return false
		}
	})
	// Return the sorted rows
	return rows, nil
}
