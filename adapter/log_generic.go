package adapter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/helper"
	"github.com/xojoc/logparse"
)

// GenericLogParser is a blunt but effective instrument
//
// It:
//
//   - loads the log (currently, only CLF or COMBINED are supported)
//   - parses the rows
//   - headings are fixed as per spec - https://www.chiark.greenend.org.uk/ucgi/~sret1/analog/olddocs.pl?version=5.23&file=logfmt.html
type GenericLogParser struct {
	BaseAdapter
}

func NewGenericLogParser(ctx data.GremelContext) data.Parser {
	p := &GenericLogParser{
		BaseAdapter: *NewBaseAdapter("log", ctx),
	}
	return p
}

func (p *GenericLogParser) Parse(input io.Reader) (*data.RowList, error) {
	logFormat := "clf"
	if p.Ctx != nil {
		if sn, ok := p.Ctx.Values().GetValue("log.format").(string); ok {
			logFormat = sn
		}
	}
	switch logFormat {
	case "clf":
		return p.parseCLF(input)
	case "combined":
		return p.parseCombined(input)
	case "syslog":
		return p.parseSyslog(input)
	}

	return nil, fmt.Errorf("Parse(%s): unsupported log format: %s", p.GetName(), logFormat)
}

func (p *GenericLogParser) parseCLF(input io.Reader) (*data.RowList, error) {
	var rows []data.Row

	scanner := bufio.NewScanner(input)
	// optionally, resize scanner's capacity for lines over 64K
	// const maxCapacity int = longLineLen  // your required line length
	// buf := make([]byte, maxCapacity)
	// scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		entry, err := logparse.Common(scanner.Text())
		if err != nil {
			log.Printf("Parse(%s): error parsing log entry: %v", p.GetName(), err)
			continue
		}
		// COnvert the parsed entry into a data.Row
		row := make(data.Row)
		row["host"] = entry.Host
		row["ident"] = "-"
		row["user"] = entry.User
		row["time"] = entry.Time
		row["request"] = entry.Request
		row["status"] = entry.Status
		row["bytes"] = entry.Bytes
		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Parse(%s): error parsing CLF log: %v", p.GetName(), err)
	}
	return data.NewRowList(rows, p.GetHeadings(rows), nil), nil
}

func (p *GenericLogParser) parseCombined(input io.Reader) (*data.RowList, error) {
	var rows []data.Row

	scanner := bufio.NewScanner(input)
	// optionally, resize scanner's capacity for lines over 64K
	// const maxCapacity int = longLineLen  // your required line length
	// buf := make([]byte, maxCapacity)
	// scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		entry, err := logparse.Common(scanner.Text())
		if err != nil {
			log.Printf("Parse(%s): error parsing log entry: %v", p.GetName(), err)
			continue
		}
		// COnvert the parsed entry into a data.Row
		row := make(data.Row)
		row["host"] = entry.Host
		row["ident"] = "-"
		row["user"] = entry.User
		row["time"] = entry.Time
		row["request"] = entry.Request
		row["status"] = entry.Status
		row["bytes"] = entry.Bytes
		row["referer"] = entry.Referer
		row["useragent"] = entry.UserAgent
		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Parse(%s): error parsing CLF log: %v", p.GetName(), err)
	}
	return data.NewRowList(rows, p.GetHeadings(rows), nil), nil
}

var syslogRegex = regexp.MustCompile(
	`^([0-9T:\.\+\-]+)\s+([0-9\.]+)\s+(\w+)\[(\d+)\]:\s+(.*)$`,
)

func (p *GenericLogParser) parseSyslog(input io.Reader) (*data.RowList, error) {
	var rows []data.Row
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		row, err := helper.ParseSyslogLine(scanner.Text())
		if err != nil {
			log.Printf("Parse(%s): error parsing syslog line: %v", p.GetName(), err)
			continue
		}
		rows = append(rows, row)
	}

	return data.NewRowList(rows, p.GetHeadings(rows), nil), nil
}
