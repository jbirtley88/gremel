package data

// Selector is used to filter and extract data items from structured dinput data in the form
// of a []map[string]any.

// Or, as it would be in JSON:
//
// [
//
//	{"foo": 0.12, "bar": "OFF", "baz": "ACTIVE"},
//	{"foo": 0.67, "bar": "ON", "baz": "CLEAR"}
//
// ]
//
// This input format was chosen because it is a common format for structured data, e.g.:
//
//   - CSV files with headings
//   - JSON arrays of objects
//   - most log files
//
// The syntax of the whereClause is similar to a WHERE clause in SQL, with the caveat that
// variable references are enclosed in curly braces, e.g.:
//
//	{foo} > 0.45 AND {bar} == "ON" OR {baz} IN ["ACTIVE", "CLEAR"]
//
// It's very lightweight and, if you've ever done SQL, you'll find it familiar.
//
// See github.com/jbirtley88/grools for more details
type Selector interface {
	Filter(input []map[string]any, whereClause string) ([]map[string]any, error)
	Order(input []map[string]any, by string) ([]map[string]any, error)
}
