package logparse

type LogFormat int

const (
	LogUnknown LogFormat = iota
	LogCLF
	LogCombined
	LogSyslog
)

func (lf LogFormat) String() string {
	switch lf {
	case LogCLF:
		return "Common Log Format"
	case LogCombined:
		return "Combined Log Format"
	case LogSyslog:
		return "Syslog"
	default:
		return "Unknown"
	}
}

func DetectLogFormat(line string) LogFormat {
	if _, err := ParseCLFLine(line); err == nil {
		return LogCLF
	}
	if _, err := ParseCombinedLogLine(line); err == nil {
		return LogCombined
	}
	if _, err := ParseSyslogLine(line); err == nil {
		return LogSyslog
	}
	return LogUnknown
}
