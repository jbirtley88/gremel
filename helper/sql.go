package helper

import (
	"fmt"
	"strings"

	"github.com/jbirtley88/gremel/data"
)

// processNextSQLLine handles SQL statement parsing and buffering
func ProcessNextSQLLine(ctx data.GremelContext, line string, sqlBuffer *[]string) (string, error) {
	// Sanitise the SQL input here to prevent SQL injection
	// cf. Bobby Tables: https://xkcd.com/327/
	// Remove SQL comments (-- comments outside of quotes)
	line = RemoveComments(line)

	// Check if this looks like a SQL statement (must start with SELECT for now)
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return "", nil // Empty line after comment removal
	}

	// If buffer is empty, check if this starts a valid SQL statement
	if len(*sqlBuffer) == 0 {
		firstWord := strings.ToUpper(strings.Split(trimmedLine, " ")[0])
		if firstWord != "SELECT" {
			return "", fmt.Errorf("only SELECT statements are supported currently")
		}
	}

	// If line contains semicolon, split at the first semicolon
	semicolonIndex := FindSemicolonOutsideQuotes(line)
	if semicolonIndex != -1 {
		// Add everything up to the semicolon to buffer
		*sqlBuffer = append(*sqlBuffer, line[:semicolonIndex])

		// Execute the complete SQL statement
		completeSQL := strings.Join(*sqlBuffer, " ")

		// Clear the buffer
		*sqlBuffer = []string{}
		return completeSQL, nil
	}

	// No semicolon found, add to buffer
	*sqlBuffer = append(*sqlBuffer, line)
	return "", nil
}

// RemoveComments removes SQL comments (--) that are outside of quotes
func RemoveComments(line string) string {
	var result strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		// Handle quote tracking
		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		// Check for comment start
		if !inSingleQuote && !inDoubleQuote && char == '-' && i+1 < len(line) && line[i+1] == '-' {
			// Found comment outside quotes, ignore rest of line
			break
		}

		result.WriteByte(char)
	}

	return result.String()
}

// FindSemicolonOutsideQuotes finds the first semicolon outside of quotes
func FindSemicolonOutsideQuotes(line string) int {
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		// Handle quote tracking
		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		} else if char == ';' && !inSingleQuote && !inDoubleQuote {
			return i
		}
	}

	return -1 // No semicolon found outside quotes
}
