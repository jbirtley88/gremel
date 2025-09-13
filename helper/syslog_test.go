package helper

import (
	"testing"
	"time"
)

func TestParseSyslogLine_ValidLines(t *testing.T) {
	tests := []struct {
		line      string
		want      SyslogEntry
		wantError bool
	}{
		{
			line: "2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd[762891]: Received disconnect from 8.8.4.4 port 5187: timeout [preauth]",
			want: SyslogEntry{
				Host:    "8.8.4.4",
				Process: "sshd",
				PID:     762891,
				Message: "Received disconnect from 8.8.4.4 port 5187: timeout [preauth]",
			},
			wantError: false,
		},
		{
			line: "2025-09-01T16:06:38.151+01:00 10.0.2.5 systemd[1]: Finished logrotate.service - Secure Shell Daemon..",
			want: SyslogEntry{
				Host:    "10.0.2.5",
				Process: "systemd",
				PID:     1,
				Message: "Finished logrotate.service - Secure Shell Daemon..",
			},
			wantError: false,
		},
		{
			line: "2025-09-05T04:40:58.520+00:00 192.168.0.10 kernel: [2025-09-05T04:40:58.520+00:00] ALERT eth{eth_num}: Link is Down",
			want: SyslogEntry{
				Host:    "192.168.0.10",
				Process: "kernel",
				PID:     0,
				Message: "[2025-09-05T04:40:58.520+00:00] ALERT eth{eth_num}: Link is Down",
			},
			wantError: false,
		},
		{
			line: "2025-09-03T05:35:43.054+01:00 172.16.5.22 myapp[175432]: [DEBUG] API error: not found",
			want: SyslogEntry{
				Host:    "172.16.5.22",
				Process: "myapp",
				PID:     175432,
				Message: "[DEBUG] API error: not found",
			},
			wantError: false,
		},
		{
			line: "2025-09-09T21:21:51.120+00:00 10.0.2.5 kubelet[449549]: E0909 21:21:51 449549 log.go:940] \"RuntimeConfig from runtime service failed\" err=\"rpc error: code = Unimplemented desc = unknown method RuntimeConfig for service runtime.v1.RuntimeService\"",
			want: SyslogEntry{
				Host:    "10.0.2.5",
				Process: "kubelet",
				PID:     449549,
				Message: "E0909 21:21:51 449549 log.go:940] \"RuntimeConfig from runtime service failed\" err=\"rpc error: code = Unimplemented desc = unknown method RuntimeConfig for service runtime.v1.RuntimeService\"",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		got, err := ParseSyslogLine(tt.line)
		if (err != nil) != tt.wantError {
			t.Errorf("ParseSyslogLine error = %v, wantError %v", err, tt.wantError)
			continue
		}
		if err == nil {
			if got["host"] != tt.want.Host {
				t.Errorf("Host = %q, want %q", got["host"], tt.want.Host)
			}
			if got["process"] != tt.want.Process {
				t.Errorf("Process = %q, want %q", got["process"], tt.want.Process)
			}
			if got["pid"] != tt.want.PID {
				t.Errorf("PID = %d, want %d", got["pid"], tt.want.PID)
			}
			if got["message"] != tt.want.Message {
				t.Errorf("Message = %q, want %q", got["message"], tt.want.Message)
			}
			if got["raw"] != tt.line {
				t.Errorf("RawLine = %q, want %q", got["raw"], tt.line)
			}
			// Timestamp check: only verify it parses to a valid time instance
			if got["timestamp"].(time.Time).IsZero() {
				t.Errorf("Timestamp should be parsed but got zero value")
			}
		}
	}
}

func TestParseSyslogLine_InvalidLines(t *testing.T) {
	invalidLines := []string{
		"", // empty line
		"Not a syslog line",
		"2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd: Received disconnect", // no PID brackets
		"2025-09-05T15:45:05.396+00:00",                                   // missing fields
		"badtimestamp 8.8.4.4 sshd[762891]: Message",                      // bad timestamp
		"2025-09-05T15:45:05.396+00:00 8.8.4.4 kernel: Message",           // kernel line without PID
	}

	for _, line := range invalidLines {
		got, err := ParseSyslogLine(line)
		if err == nil {
			t.Errorf("ParseSyslogLine(%q) expected error, got %+v", line, got)
		}
	}
}

func TestParseSyslogLine_TimestampVariants(t *testing.T) {
	lines := []struct {
		line     string
		expected string // Expected time layout
	}{
		{
			line:     "2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd[1]: msg",
			expected: "2006-01-02T15:04:05.000-07:00",
		},
		{
			line:     "2025-09-05T15:45:05.396+01:00 8.8.4.4 sshd[1]: msg",
			expected: "2006-01-02T15:04:05.000+07:00",
		},
		{
			line:     "2025-09-05T15:45:05.396Z 8.8.4.4 sshd[1]: msg",
			expected: "2006-01-02T15:04:05.000Z",
		},
	}
	for _, l := range lines {
		got, err := ParseSyslogLine(l.line)
		if err != nil {
			t.Errorf("ParseSyslogLine(%q) unexpected error: %v", l.line, err)
			continue
		}
		if got["timestamp"].(time.Time).IsZero() {
			t.Errorf("ParseSyslogLine(%q) failed to parse timestamp", l.line)
		}
	}
}

func TestParseSyslogLine_PIDParsing(t *testing.T) {
	lines := []struct {
		line string
		pid  int
	}{
		{
			line: "2025-09-05T15:45:05.396+00:00 8.8.4.4 sshd[999]: msg",
			pid:  999,
		},
		{
			line: "2025-09-05T15:45:05.396+00:00 8.8.4.4 systemd[0]: msg",
			pid:  0,
		},
	}
	for _, l := range lines {
		got, err := ParseSyslogLine(l.line)
		if err != nil {
			t.Errorf("ParseSyslogLine(%q) unexpected error: %v", l.line, err)
			continue
		}
		if got["pid"] != l.pid {
			t.Errorf("ParseSyslogLine(%q) PID = %d, want %d", l.line, got["pid"], l.pid)
		}
	}
}

func TestParseSyslogLine_KernelNoPID(t *testing.T) {
	line := "2025-09-05T04:40:58.520+00:00 192.168.0.10 kernel: [2025-09-05T04:40:58.520+00:00] ALERT eth{eth_num}: Link is Down"
	got, err := ParseSyslogLine(line)
	if err != nil {
		t.Errorf("ParseSyslogLine(%q) unexpected error: %v", line, err)
	}
	if got["pid"] != 0 {
		t.Errorf("ParseSyslogLine(%q) kernel PID = %d, want 0", line, got["pid"])
	}
}

// Example of how to test ParseSyslogLine with a table of real samples
func TestParseSyslogLine_RealSamples(t *testing.T) {
	lines := []string{
		"2025-09-12T23:11:42.376+00:00 8.8.4.4 nginx[425545]: worker process started",
		"2025-09-10T19:50:56.067+01:00 8.8.4.4 sshd[734353]: Received disconnect from 8.8.4.4 port 20240: connection refused [preauth]",
	}
	for _, line := range lines {
		got, err := ParseSyslogLine(line)
		if err != nil {
			t.Errorf("ParseSyslogLine(%q) unexpected error: %v", line, err)
		}
		if got["host"] == "" || got["process"] == "" || got["message"] == "" {
			t.Errorf("ParseSyslogLine(%q) incomplete parse: %+v", line, got)
		}
	}
}
