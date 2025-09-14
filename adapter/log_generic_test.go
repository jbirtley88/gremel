package adapter

import (
	"bytes"
	"context"
	"os"
	"sort"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLFGeneric(t *testing.T) {
	f, err := os.Open("../test_resources/clf.log")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	ctx := data.NewGremelContext(context.TODO())
	// ctx.Values().SetValue("log.format", "clf")
	p := NewGenericLogParser(ctx)
	require.NotNil(t, p)

	expectedHeadings := []string{"authuser", "host", "ident", "latency", "method", "path", "proto", "size", "status", "time"}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
	assert.Equal(t, 1000, len(rows.Rows))
}

func TestCombinedGeneric(t *testing.T) {
	f, err := os.Open("../test_resources/combined.log")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	ctx := data.NewGremelContext(context.TODO())
	// ctx.Values().SetValue("log.format", "combined")
	p := NewGenericLogParser(ctx)
	require.NotNil(t, p)

	expectedHeadings := []string{"host", "ident", "latency", "referer", "request", "size", "status", "time", "user", "useragent"}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
	assert.Equal(t, 1000, len(rows.Rows))
}

func TestSyslogGeneric(t *testing.T) {
	f, err := os.Open("../test_resources/syslog.log")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	ctx := data.NewGremelContext(context.TODO())
	// ctx.Values().SetValue("log.format", "syslog")
	p := NewGenericLogParser(ctx)
	require.NotNil(t, p)

	expectedHeadings := []string{"host", "message", "pid", "process", "raw", "timestamp"}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
	assert.Equal(t, 10000, len(rows.Rows))
}

func TestGenericLogParser_UnsupportedFormat(t *testing.T) {
	ctx := data.NewGremelContext(context.TODO())
	ctx.Values().SetValue("log.format", "unknown_format")
	p := NewGenericLogParser(ctx)
	require.NotNil(t, p)

	buf := bytes.NewBuffer([]byte(
		`unrecognised log line
another unrecognised line`,
	))

	_, err := p.Parse(buf)
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "unrecognised log format")
}
