package adapter

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/jbirtley88/gremel/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Deal with the case where the top-level is a list:
// [
//
//	{"foo": 1, "bar": 2, ...},
//	{"foo": 1, "bar": 2, ...},
//
// ]
func TestJsonGenericWhenObjectListIsTopLevel(t *testing.T) {
	f, err := os.Open("../test_resources/accounts_top_level.json")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	p := NewGenericJsonParser(data.NewGremelContext(context.TODO()))
	require.NotNil(t, p)

	expectedHeadings := []string{
		"email",
		"id",
		"mac_address",
		"percent",
		"username",
	}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
}

func TestJsonGenericWhenObjectListIsNested(t *testing.T) {
	f, err := os.Open("../test_resources/accounts_nested.json")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	p := NewGenericJsonParser(data.NewGremelContext(context.TODO()))
	require.NotNil(t, p)

	expectedHeadings := []string{
		"email",
		"id",
		"mac_address",
		"percent",
		"username",
	}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
}

func TestJsonGenericWhenObjectListIsNestedAndRootIsGiven(t *testing.T) {
	f, err := os.Open("../test_resources/accounts_nested.json")
	require.Nil(t, err)
	require.NotNil(t, f)
	defer f.Close()

	ctx := data.NewGremelContext(context.TODO(), data.NewMetadata().SetValue("json_root", "data.list"))
	p := NewGenericJsonParser(ctx)
	require.NotNil(t, p)

	expectedHeadings := []string{
		"email",
		"id",
		"mac_address",
		"percent",
		"username",
	}

	rows, err := p.Parse(f)
	require.Nil(t, err)
	require.NotNil(t, rows)
	require.NotNil(t, rows.Rows)
	require.NotNil(t, rows.Headings)

	sort.Strings(rows.Headings)
	assert.Equal(t, expectedHeadings, rows.Headings)
	assert.Equal(t, 1000, len(rows.Rows), "Should be 10 rows")
}
