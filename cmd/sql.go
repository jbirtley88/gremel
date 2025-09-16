package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jbirtley88/gremel/apiimpl"
	"github.com/jbirtley88/gremel/data"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Runs an interactive SQL shell",
	Run:   RunSQL,
}

var silentMode bool

func init() {
	rootCmd.AddCommand(sqlCmd)
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent", "q", false, "Silent mode - suppress output except for query results")
}

func RunSQL(cmd *cobra.Command, args []string) {
	if !silentMode {
		fmt.Println("Type '.help' or '?' for help")
	}
	err := runSQL(args)
	if err != nil {
		log.Errorf("%s: Error running SQL command: %v", cmd.Name(), err)
	}
}

func runSQL(args []string) error {
	// Read one line of text at a time from stdin
	reader := bufio.NewReader(os.Stdin)
	prompt := "gremel> "
	var sqlBuffer []string
	ctx := data.NewGremelContext(context.Background())
	ctx.Values().SetValue("silent", silentMode)
	for {
		// Use continuation prompt if we're building a multi-line SQL statement
		currentPrompt := prompt
		if len(sqlBuffer) > 0 {
			currentPrompt = "    ...> "
		}

		if !silentMode {
			fmt.Print(currentPrompt)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if !silentMode {
					fmt.Println("\nExiting...")
				}
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Process the command
		tokens := strings.Split(line, " ")
		firstToken := strings.ToLower(tokens[0])

		switch firstToken {
		case "exit", "quit", ".exit", ".quit", ".q":
			fmt.Println("Exiting...")
			return nil
		case ".mount":
			err := doMount(ctx, tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case ".help", "help", "?":
			err := doHelp(ctx, tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case ".tables":
			err := doTables(ctx, tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case ".schema":
			err := doSchema(ctx, tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case ".silent":
			err := doSilent(ctx, tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		default:
			// Handle SQL statements
			err := processSQLLine(ctx, line, &sqlBuffer)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				sqlBuffer = []string{} // Reset buffer on error
			}
		}
	}
}

// doSilent handles the .silent command
func doSilent(ctx data.GremelContext, tokens []string) error {
	if len(tokens) != 2 {
		return fmt.Errorf("usage: .silent on|off")
	}
	switch strings.ToLower(tokens[1]) {
	case "on":
		ctx.Values().SetValue("silent", true)
		silentMode = true
	case "off":
		ctx.Values().SetValue("silent", true)
		silentMode = false
		fmt.Println("Silent mode disabled")
	default:
		return fmt.Errorf("usage: .silent on|off")
	}
	return nil
}

func doHelp(ctx data.GremelContext, tokens []string) error {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	w.Write([]byte("Available commands:\n"))
	w.Write([]byte(".help\tShow this help message\n"))
	w.Write([]byte(".quit or .exit or .q\tExit the shell\n"))
	w.Write([]byte(".mount [tablename [<file_path>]]\tMount a data source\n"))
	w.Write([]byte(".tables\tList all tables\n"))
	w.Write([]byte(".schema <tablename>\tShow schema of a table\n"))
	w.Write([]byte(".headings on|off\tEnable or disable column headings\n"))
	w.Write([]byte(".silent on|off\tEnable or disable silent mode\n"))
	w.Write([]byte("SELECT ...;\tExecute a SQL SELECT statement\n"))
	w.Flush()
	return nil
}

// doMount handles the .mount command
func doMount(ctx data.GremelContext, tokens []string) error {
	switch len(tokens) {
	case 1:
		// .mount
		mountInfo, err := apiimpl.GetMount(ctx, "")
		if err != nil {
			return fmt.Errorf("error getting mount info: %v", err)
		}
		// TODO(john): better formatting
		for k, v := range mountInfo {
			fmt.Printf("%s: %v\n", k, v)
		}
		return nil

	case 2:
		// .mount tablename
		mountInfo, err := apiimpl.GetMount(ctx, tokens[1])
		if err != nil {
			return fmt.Errorf("error getting mount info: %v", err)
		}
		// TODO(john): better formatting
		for k, v := range mountInfo {
			fmt.Printf("%s: %v\n", k, v)
		}
		return nil

	case 3:
		// .mount NAME /path/to/file
		err := apiimpl.Mount(ctx, tokens[1], tokens[2])
		if err != nil {
			return fmt.Errorf("Error mounting file: %v\n", err)
		}
		return nil

	default:
		return fmt.Errorf("usage: .mount [tablename [<file_path>]]")
	}

	// TODO: Support for mounting http:// and https:// URLs as data sources
}

// doTables handles the .tables command
func doTables(ctx data.GremelContext, tokens []string) error {
	tables, err := apiimpl.GetTables(ctx)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error getting tables: %v\n", err))
	}
	for _, table := range tables {
		fmt.Println(table)
	}
	return nil
}

// doSchema handles the .schema command
func doSchema(ctx data.GremelContext, tokens []string) error {
	schema, err := apiimpl.GetSchema(ctx, tokens[1])
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error getting schema: %v\n", err))
		return err
	}
	for k, v := range schema {
		fmt.Printf("%s: %s\n", k, v)
	}
	return nil
}

// processSQLLine handles SQL statement parsing and buffering
func processSQLLine(ctx data.GremelContext, line string, sqlBuffer *[]string) error {
	// Sanitise the SQL input here to prevent SQL injection
	// cf. Bobby Tables: https://xkcd.com/327/
	// Remove SQL comments (-- comments outside of quotes)
	line = removeComments(line)

	// Check if this looks like a SQL statement (must start with SELECT for now)
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return nil // Empty line after comment removal
	}

	// If buffer is empty, check if this starts a valid SQL statement
	if len(*sqlBuffer) == 0 {
		firstWord := strings.ToUpper(strings.Split(trimmedLine, " ")[0])
		if firstWord != "SELECT" {
			return fmt.Errorf("only SELECT statements are supported currently")
		}
	}

	// If line contains semicolon, split at the first semicolon
	semicolonIndex := findSemicolonOutsideQuotes(line)
	if semicolonIndex != -1 {
		// Add everything up to the semicolon to buffer
		*sqlBuffer = append(*sqlBuffer, line[:semicolonIndex])

		// Execute the complete SQL statement
		completeSQL := strings.Join(*sqlBuffer, " ")
		err := executeSQL(ctx, completeSQL)
		if err != nil {
			*sqlBuffer = []string{} // Reset buffer
			return err
		}

		// Clear the buffer
		*sqlBuffer = []string{}
		return nil
	}

	// No semicolon found, add to buffer
	*sqlBuffer = append(*sqlBuffer, line)
	return nil
}

// removeComments removes SQL comments (--) that are outside of quotes
func removeComments(line string) string {
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

// findSemicolonOutsideQuotes finds the first semicolon outside of quotes
func findSemicolonOutsideQuotes(line string) int {
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

// separator returns a string of hyphens with the same length as the input string
func separator(s string) string {
	return strings.Repeat("-", len(s))
}

// executeSQL executes a complete SQL statement
func executeSQL(ctx data.GremelContext, sql string) error {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil
	}

	// TODO: Implement actual SQL execution using the database
	if ctx.Values().GetBool("verbose") {
		fmt.Printf("Executing SQL: %s\n", sql)
	}
	rows, columns, err := apiimpl.Query(data.NewGremelContext(context.Background()), sql)
	if err != nil {
		return fmt.Errorf("executeSQL(): %w", err)
	}

	// Print results in tabular format
	// Use tabwriter for aligned columns
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)

	// Column names
	if !silentMode {
		if len(columns) > 0 {
			for _, col := range columns {
				_, _ = w.Write([]byte(col + "\t"))
			}
			_, _ = w.Write([]byte("\n"))

			// Heading separator
			for _, col := range columns {
				_, _ = w.Write([]byte(separator(col) + "\t"))
			}
			_, _ = w.Write([]byte("\n"))
		}
	}

	// Process rows
	if len(rows) > 0 {
		for _, row := range rows {
			for _, col := range columns {
				valueAsString := fmt.Sprint(row[col])
				_, _ = w.Write([]byte(valueAsString + "\t"))
			}
			_, _ = w.Write([]byte("\n"))
		}
		_ = w.Flush()
	}

	if !silentMode {
		fmt.Printf("%d rows\n", len(rows))
	}
	return nil
}
