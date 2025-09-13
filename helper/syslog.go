package helper

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jbirtley88/gremel/data"
)

// SyslogEntry is a struct representing a syslog line
type SyslogEntry struct {
	Timestamp time.Time
	Host      string
	Process   string
	PID       int
	Message   string
	RawLine   string // original line for reference/debug
}

// Regex to match syslog lines - updated to handle various formats
var syslogRegex = regexp.MustCompile(
	`^([0-9T:\.\+\-Z]+)\s+([^\s]+)\s+([^\[\s:]+)(?:\[(\d+)\])?\:\s+(.*)$`,
)

func ParseSyslogLine(line string) (data.Row, error) {
	matches := syslogRegex.FindStringSubmatch(line)
	if len(matches) < 6 {
		return nil, fmt.Errorf("no syslog match")
	}

	tsStr := matches[1]
	host := matches[2]
	process := matches[3]
	pidStr := matches[4] // May be empty for kernel entries
	msg := matches[5]

	// PID brackets are required for all processes except "kernel"
	if process != "kernel" && pidStr == "" {
		return nil, fmt.Errorf("non-kernel process %q requires PID brackets", process)
	}

	// Kernel entries without PID should have the format: [timestamp] LEVEL message
	if process == "kernel" && pidStr == "" {
		kernelMsgPattern := regexp.MustCompile(`^\[[0-9T:\.\+\-Z]+\]\s+[A-Z]+\s+`)
		if !kernelMsgPattern.MatchString(msg) {
			return nil, fmt.Errorf("kernel message without PID must have format '[timestamp] LEVEL message', got: %q", msg)
		}
	}

	// Try to parse timestamp with different formats
	var timestamp time.Time
	var err error

	// Common timestamp formats in syslog
	timeFormats := []string{
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000+07:00",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range timeFormats {
		timestamp, err = time.Parse(format, tsStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp %q: %w", tsStr, err)
	}

	pid := 0
	if pidStr != "" {
		pid, err = strconv.Atoi(pidStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PID %q: %w", pidStr, err)
		}
	}

	row := make(data.Row)
	row["timestamp"] = timestamp
	row["host"] = host
	row["process"] = process
	row["pid"] = pid
	row["message"] = msg
	row["raw"] = line
	return row, nil
}

func main() {
	f, err := os.Open("/var/log/syslog")
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := ParseSyslogLine(line)
		if err != nil {
			fmt.Printf("Could not parse line: %s\nError: %v\n", strings.TrimSpace(line), err)
			continue
		}
		fmt.Printf("%+v\n", entry)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
