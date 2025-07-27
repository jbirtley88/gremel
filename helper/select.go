package helper

// Select extracts named columns from the provided rows.
//
// In AuthFairy, it executes after Filter(), and before Order()
func Select(needle string, haystack []map[string]any) ([]map[string]any, error) {
	if needle == "" {
		// Nothing to do
		return haystack, nil
	}

	return haystack, nil
}
