package logparse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jbirtley88/gremel/data"
)

// LogEntry represents one parsed Apache combined log line with an optional trailing latency (ms).
type LogEntry struct {
	Host        string `json:"host"`
	Ident       string `json:"ident"`
	User        string `json:"user"`
	Timestamp   int64  `json:"timestamp"`
	Request     string `json:"request"`  // full raw request line e.g. `GET /path HTTP/1.1`
	Method      string `json:"method"`   // parsed from Request, if available
	Path        string `json:"path"`     // parsed from Request, if available
	Protocol    string `json:"protocol"` // parsed from Request, if available
	Status      int    `json:"status"`
	Bytes       int    `json:"bytes"` // '-' becomes 0
	Referer     string `json:"referer"`
	UserAgent   string `json:"user_agent"`
	LatencyMS   int    `json:"latency_ms"` // optional final field, defaults to 0 if missing
	RawLine     string `json:"raw_line"`
	ParseErrMsg string `json:"parse_error,omitempty"`
}

// Strict Apache Combined + optional trailing integer latency (ms).
// Pattern: %h %l %u [%t] "%r" %>s %b "%{Referer}i" "%{User-agent}i" [latency_ms]
// Example:
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com" "Mozilla/4.08" 123
var combinedWithOptionalLatency = regexp.MustCompile(
	`^(\S+)\s+(\S+)\s+(\S+)\s+\[([^\]]+)\]\s+"([^"]*)"\s+(\d{3})\s+(\S+)\s+"([^"]*)"\s+"([^"]*)"(?:\s+(\d+))?\s*$`,
)

func ParseCombinedLogLine(line string) (data.Row, error) {
	combinedEntry, err := parseCombinedLogLine(line)
	if err != nil {
		return nil, fmt.Errorf("ParseCombinedLogLine(): %w", err)
	}

	return data.Row{
		"host":      combinedEntry.Host,
		"user":      combinedEntry.User,
		"ident":     combinedEntry.Ident,
		"size":      combinedEntry.Bytes,
		"status":    combinedEntry.Status,
		"referer":   combinedEntry.Referer,
		"request":   combinedEntry.Request,
		"latency":   combinedEntry.LatencyMS,
		"time":      combinedEntry.Timestamp,
		"useragent": combinedEntry.UserAgent,
	}, nil
}

// ParseCombinedLogLine parses a single Apache combined log line with an optional trailing latency field.
// If the latency field is absent, LatencyMS will be 0.
// Timestamp layout is [02/Jan/2006:15:04:05 -0700].
func parseCombinedLogLine(line string) (*LogEntry, error) {
	m := combinedWithOptionalLatency.FindStringSubmatch(line)
	if m == nil {
		return nil, errors.New("line does not match Apache combined format with optional latency")
	}

	host := m[1]
	ident := m[2]
	user := m[3]
	timeStr := m[4]
	request := m[5]
	statusStr := m[6]
	bytesStr := m[7]
	referer := m[8]
	userAgent := m[9]
	latencyStr := m[10] // may be empty

	// Parse timestamp: [02/Jan/2006:15:04:05 -0700]
	// m[4] already excludes the surrounding brackets.
	ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr)
	if err != nil {
		// Keep zero value but return error to caller; they may want to inspect RawLine/ParseErrMsg.
		return &LogEntry{
			Host:        host,
			Ident:       ident,
			User:        user,
			Timestamp:   0,
			Request:     request,
			Status:      atoiDefault(statusStr, 0),
			Bytes:       parseBytes(bytesStr),
			Referer:     referer,
			UserAgent:   userAgent,
			LatencyMS:   atoiDefault(latencyStr, 0),
			RawLine:     line,
			ParseErrMsg: fmt.Sprintf("failed to parse timestamp: %v", err),
		}, nil
	}

	status := atoiDefault(statusStr, 0)
	bytes := parseBytes(bytesStr)
	latency := 0
	if latencyStr != "" {
		latency = atoiDefault(latencyStr, 0)
	}

	method, path, proto := splitRequest(request)

	entry := &LogEntry{
		Host:      host,
		Ident:     ident,
		User:      user,
		Timestamp: ts.UnixMilli(),
		Request:   request,
		Method:    method,
		Path:      path,
		Protocol:  proto,
		Status:    status,
		Bytes:     bytes,
		Referer:   referer,
		UserAgent: userAgent,
		LatencyMS: latency,
		RawLine:   line,
	}
	return entry, nil
}

func atoiDefault(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// parseBytes handles %b field which can be "-" meaning no bytes sent.
func parseBytes(s string) int {
	if s == "-" {
		return 0
	}
	return atoiDefault(s, 0)
}

// splitRequest extracts method, path, and protocol from the "%r" field.
// If the request is empty or malformed, returns empty components and keeps Request raw.
func splitRequest(req string) (method, path, proto string) {
	// Typical: "GET /path HTTP/1.1"
	parts := strings.SplitN(req, " ", 3)
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	}
	// Some servers might emit "-" or empty for request
	return "", "", ""
}
