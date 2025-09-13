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

// CLFEntry represents a parsed Common Log Format line with an optional trailing latency field (milliseconds).
// Format (with optional latency_ms at the end):
// %h %l %u [%t] "%r" %>s %b [latency_ms]
// Example:
// 165.23.106.237 - alice [04/Sep/2025:19:12:36 -0700] "GET /logout?token=abcdef HTTP/1.1" 200 11281 850
type CLFEntry struct {
	Host       string
	Ident      string
	AuthUser   string
	Timestamp  int64 // Unix timestamp in millis
	Method     string
	Path       string
	Protocol   string
	Status     int
	Bytes      int // 0 if '-' or unparsable
	LatencyMs  int // optional trailing field; defaults to 0 if absent
	RawRequest string
	RawLine    string
}

// Regex to match CLF with optional latency field as the final token.
// Groups:
// 1 host, 2 ident, 3 authuser, 4 timestamp, 5 request, 6 status, 7 bytes, 8 latency(optional)
var clfRegex = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+\[([^\]]+)\]\s+"([^"]+)"\s+(\d{3})\s+(\S+)(?:\s+(\d+))?$`)

func ParseCLFLine(line string) (data.Row, error) {
	clfEntry, err := parseCLFLine(line)
	if err != nil {
		return nil, fmt.Errorf("ParseCLFLine(): %w", err)
	}

	return data.Row{
		"host":     clfEntry.Host,
		"path":     clfEntry.Path,
		"size":     clfEntry.Bytes,
		"ident":    clfEntry.Ident,
		"method":   clfEntry.Method,
		"status":   clfEntry.Status,
		"authuser": clfEntry.AuthUser,
		"proto":    clfEntry.Protocol,
		"latency":  clfEntry.LatencyMs,
		"time":     clfEntry.Timestamp}, nil
}

// ParseCLFLine parses a single CLF line with optional latency in ms.
func parseCLFLine(line string) (*CLFEntry, error) {
	m := clfRegex.FindStringSubmatch(line)
	if m == nil {
		return nil, errors.New("line does not match CLF pattern")
	}

	host := m[1]
	ident := m[2]
	authUser := m[3]
	tsStr := m[4]
	req := m[5]
	statusStr := m[6]
	bytesStr := m[7]
	latencyStr := ""
	if len(m) >= 9 {
		latencyStr = m[8]
	}

	// Parse timestamp like: 04/Sep/2025:19:12:36 -0700
	const tsLayout = "02/Jan/2006:15:04:05 -0700"
	ts, err := time.Parse(tsLayout, tsStr)
	if err != nil {
		return nil, fmt.Errorf("timestamp parse error: %w", err)
	}

	// Parse request: METHOD SP PATH SP PROTOCOL (best-effort)
	method, path, proto := "", "", ""
	reqParts := strings.SplitN(req, " ", 3)
	if len(reqParts) == 3 {
		method, path, proto = reqParts[0], reqParts[1], reqParts[2]
	} else if len(reqParts) == 2 {
		method, path = reqParts[0], reqParts[1]
	} else if len(reqParts) == 1 {
		method = reqParts[0]
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, fmt.Errorf("status parse error: %w", err)
	}

	// %b can be '-' meaning no bytes sent
	bytes := 0
	if bytesStr != "-" {
		if n, err := strconv.Atoi(bytesStr); err == nil {
			bytes = n
		}
	}

	latency := 0
	if latencyStr != "" {
		if n, err := strconv.Atoi(latencyStr); err == nil {
			latency = n
		}
	}

	return &CLFEntry{
		Host:       host,
		Ident:      ident,
		AuthUser:   authUser,
		Timestamp:  ts.UnixMilli(),
		Method:     method,
		Path:       path,
		Protocol:   proto,
		Status:     status,
		Bytes:      bytes,
		LatencyMs:  latency,
		RawRequest: req,
		RawLine:    line,
	}, nil
}
