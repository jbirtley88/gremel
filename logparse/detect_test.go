package logparse

import (
	"testing"
)

func TestDetectLines(t *testing.T) {
	tests := []struct {
		line      string
		want      LogFormat
		wantError bool
	}{
		{
			line:      `2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd[762891]: Received disconnect from 8.8.4.4 port 5187: timeout [preauth]`,
			want:      LogSyslog,
			wantError: false,
		},
		{
			line:      `205.15.228.48 - dan [12/Sep/2025:21:03:41 +0000] "POST /search HTTP/1.1" 200 3184 670`,
			want:      LogCLF,
			wantError: false,
		},
		{
			line:      `139.83.32.41 - - [08/Sep/2025:17:35:20 +0530] "POST /posts?cat=books HTTP/1.1" 200 2572 "https://google.com" "Googlebot/2.1 (+http://www.google.com/bot.html)" 859`,
			want:      LogCombined,
			wantError: false,
		},
	}

	for _, tt := range tests {
		got := DetectLogFormat(tt.line)
		if got != tt.want {
			t.Errorf("DetectLogFormat() = %v, want %v", got, tt.want)
		}
	}
}

func TestDetectInvalid(t *testing.T) {
	tests := []struct {
		line      string
		want      LogFormat
		wantError bool
	}{
		{
			line:      "This is not a valid log line",
			want:      LogUnknown,
			wantError: false,
		},
		{
			line:      "",
			want:      LogUnknown,
			wantError: false,
		},
	}

	for _, tt := range tests {
		got := DetectLogFormat(tt.line)
		if got != tt.want {
			t.Errorf("DetectLogFormat() = %v, want %v", got, tt.want)
		}
	}
}
