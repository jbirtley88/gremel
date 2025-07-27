package helper

import (
	"context"
	"strings"
)

type BaseSelector struct {
	Name string
	Ctx  context.Context
}

func NewBaseSelector(ctx context.Context) *BaseSelector {
	a := &BaseSelector{
		Name: "default",
		Ctx:  ctx,
	}
	return a
}

func (a *BaseSelector) GetName() string {
	return a.Name
}

func (a *BaseSelector) Select(needle string, haystack []map[string]any) ([]map[string]any, error) {
	if needle == "" {
		// Nothing to do
		return haystack, nil
	}

	// Figure out the columns of interest
	columnsOfInterest := make(map[string]bool, 0)
	for _, coi := range strings.Split(needle, ",") {
		coi = strings.TrimSpace(coi)
		if coi != "" {
			columnsOfInterest[coi] = true
		}
	}

	rowsOfInterest := make([]map[string]any, 0)

	for _, row := range haystack {
		// First make sure this is a row we care about
		anyColumnsOfInterest := false
		for coi := range columnsOfInterest {
			if _, exists := row[coi]; exists {
				anyColumnsOfInterest = true
				break
			}
		}
		if anyColumnsOfInterest {
			rowsOfInterest = append(rowsOfInterest, row)
		}
	}

	// Now extract the columns we care about
	result := make([]map[string]any, 0)
	for _, row := range rowsOfInterest {
		newRow := make(map[string]any)
		for coi := range columnsOfInterest {
			if value, exists := row[coi]; exists {
				newRow[coi] = value
			}
		}
		result = append(result, newRow)
	}
	return result, nil
}

func (a *BaseSelector) Where(rows []map[string]any, whereClause string) ([]map[string]any, error) {
	if whereClause == "" {
		// Nothing to do
		return rows, nil
	}

	result := make([]map[string]any, 0)
	for _, row := range rows {
		matches, err := Matches(a.Ctx, row, whereClause)
		if err != nil {
			return nil, err
		}
		if matches {
			result = append(result, row)
		}
	}
	return result, nil
}

func (a *BaseSelector) Order(input []map[string]any, by string) ([]map[string]any, error) {
	return Sort(input, by)
}
