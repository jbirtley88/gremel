package logparse

import (
	"testing"

	"github.com/jbirtley88/gremel/data"
)

func TestParseLines(t *testing.T) {
	tests := []struct {
		line      string
		wantError bool
		expected  data.Row
	}{
		{
			line:      `2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd[762891]: Received disconnect from 8.8.4.4 port 5187: timeout [preauth]`,
			wantError: false,
			expected: data.Row{
				"host": "8.8.4.4",
				"proc": "sshd",
				"pid":  762891,
				"msg":  "Received disconnect from 8.8.4.4 port 5187: timeout [preauth]",
			},
		},
		{
			line:      `205.15.228.48 - dan [12/Sep/2025:21:03:41 +0000] "POST /search HTTP/1.1" 200 3184 670`,
			wantError: false,
			expected: data.Row{
				"host":      "205.15.228.48",
				"user":      "dan",
				"timestamp": "12/Sep/2025:21:03:41 +0000",
				"method":    "POST",
				"path":      "/search",
				"protocol":  "HTTP/1.1",
				"status":    200,
				"size":      3184,
			},
		},
		{
			line:      `139.83.32.41 - - [08/Sep/2025:17:35:20 +0530] "POST /posts?cat=books HTTP/1.1" 200 2572 "https://google.com" "Googlebot/2.1 (+http://www.google.com/bot.html)" 859`,
			wantError: false,
			expected: data.Row{
				"host":       "139.83.32.41",
				"user":       "-",
				"timestamp":  "08/Sep/2025:17:35:20 +0530",
				"method":     "POST",
				"path":       "/posts?cat=books",
				"protocol":   "HTTP/1.1",
				"status":     200,
				"size":       2572,
				"referrer":   "https://google.com",
				"user_agent": "Googlebot/2.1 (+http://www.google.com/bot.html)",
				"duration":   859,
			},
		},
	}

	for _, tt := range tests {
		got, err := ParseLine(tt.line)
		if (err != nil) != tt.wantError {
			t.Errorf("ParseLine() error = %v, wantError %v", err, tt.wantError)
			continue
		}
		if err == nil && got == nil {
			t.Errorf("ParseLine() = %v, want non-nil", got)
		}
	}
}
