package logparse

import (
	"fmt"

	"github.com/jbirtley88/gremel/data"
)

func ParseLine(line string) (data.Row, error) {
	if row, err := ParseCLFLine(line); err == nil {
		return row, nil
	}
	if row, err := ParseCombinedLogLine(line); err == nil {
		return row, nil
	}
	if row, err := ParseSyslogLine(line); err == nil {
		return row, nil
	}
	return nil, fmt.Errorf("ParseLine: unrecognised log format")
}
