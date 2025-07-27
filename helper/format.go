package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func TextOutput(rows []map[string]any, headings []string, writeTo io.Writer) error {
	byteWriter := bytes.NewBuffer(make([]byte, 0))
	w := tabwriter.NewWriter(byteWriter, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)

	for _, heading := range headings {
		if _, err := w.Write([]byte(heading + "\t  ")); err != nil {
			return fmt.Errorf("error writing headings: %w", err)
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return fmt.Errorf("error writing newline after headings: %w", err)
	}

	for _, row := range rows {
		for _, heading := range headings {
			if value, exists := row[heading]; exists {
				if _, err := w.Write([]byte(formatValue(value) + "\t  ")); err != nil {
					return fmt.Errorf("error writing row value: %w", err)
				}
			} else {
				if _, err := w.Write([]byte("\t  ")); err != nil {
					return fmt.Errorf("error writing empty column: %w", err)
				}
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return fmt.Errorf("error writing newline after row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	writeTo.Write(byteWriter.Bytes())
	return nil
}

func TableOutput(rows []map[string]any, headings []string, writeTo io.Writer) error {
	byteWriter := bytes.NewBuffer(make([]byte, 0))
	w := tabwriter.NewWriter(byteWriter, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns|tabwriter.AlignRight|tabwriter.TabIndent)

	for _, heading := range headings {
		if _, err := w.Write([]byte(heading + "\t  ")); err != nil {
			return fmt.Errorf("error writing headings: %w", err)
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return fmt.Errorf("error writing newline after headings: %w", err)
	}

	for _, row := range rows {
		for _, heading := range headings {
			if value, exists := row[heading]; exists {
				if _, err := w.Write([]byte(formatValue(value) + "\t  ")); err != nil {
					return fmt.Errorf("error writing row value: %w", err)
				}
			} else {
				if _, err := w.Write([]byte("\t  ")); err != nil {
					return fmt.Errorf("error writing empty column: %w", err)
				}
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return fmt.Errorf("error writing newline after row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	writeTo.Write(byteWriter.Bytes())
	return nil
}

func CSVOutput(rows []map[string]any, headings []string, writeTo io.Writer) error {
	byteWriter := bytes.NewBuffer(make([]byte, 0))

	// Write headings
	for i, heading := range headings {
		if i > 0 {
			byteWriter.WriteString(",")
		}
		byteWriter.WriteString(heading)
	}
	byteWriter.WriteString("\n")

	// Write rows
	for _, row := range rows {
		for i, heading := range headings {
			if i > 0 {
				byteWriter.WriteString(",")
			}
			if value, exists := row[heading]; exists {
				formatted := formatValue(value)
				if strings.Contains(formatted, ",") || strings.Contains(formatted, "\"") {
					formatted = "\"" + strings.ReplaceAll(formatted, "\"", "\"\"") + "\""
				}
				byteWriter.WriteString(formatted)
			}
		}
		byteWriter.WriteString("\n")
	}

	writeTo.Write(byteWriter.Bytes())
	return nil
}

func JSONOutput(rows []map[string]any, writeTo io.Writer) error {
	encoder := json.NewEncoder(writeTo)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rows); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}
	return nil
}

func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}
