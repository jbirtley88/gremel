package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jbirtley88/gremel/data"

	"github.com/spf13/viper"
)

// GenericJsonParser is a blunt but effective instrument
//
// It:
//
//   - unmarshals the JSON into a map
//   - looks for the first slice of map[string]any (breadth-first recursive)
//   - uses that as the row data
type GenericJsonParser struct {
	BaseAdapter
}

func NewGenericJsonParser(ctx data.GremelContext) data.Parser {
	p := &GenericJsonParser{
		BaseAdapter: *NewBaseAdapter("json", ctx),
	}
	return p
}

func (p *GenericJsonParser) GetHeadings(rows []map[string]any) []string {
	headings := []string{}

	// Check the context
	if p.Ctx != nil {
		if selectValue := p.Ctx.Values().GetString("select"); selectValue != "" && selectValue != "*" {
			headings = []string{}
			for _, columnName := range strings.Split(selectValue, ",") {
				headings = append(headings, strings.TrimSpace(columnName))
			}
		}
	}

	// Nothing in the context.
	// We need to grab tem from the map.
	// Sadly, keys are not ordered so we're going to get a pseudo-random order
	if len(headings) == 0 {
		for k := range rows[0] {
			headings = append(headings, k)
		}
	}

	return headings
}

func (p *GenericJsonParser) Parse(input io.Reader) ([]map[string]any, []string, error) {

	// Step 2: Check the context, see if a root for the data has been specified
	//         via the 'root' value
	if p.Ctx != nil {
		if dataLocation := p.Ctx.Values().GetString("data"); dataLocation != "" {
			rows, err := p.getJsonObjectList(input, dataLocation)
			if err != nil {
				return nil, nil, fmt.Errorf("%s.Parse(): %s", p.Name, err.Error())
			}

			return rows, p.GetHeadings(rows), nil
		}
	}

	// Step 1: Easy mode: try unmarshalling into a []map[string]any
	jsonBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, nil, fmt.Errorf("%s.Parse(): %s", p.Name, err.Error())
	}

	var sliceOfMap []map[string]any
	err = json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&sliceOfMap)
	if err == nil {
		// It is already a []map[string]any
		return sliceOfMap, p.GetHeadings(sliceOfMap), nil
	}

	// Step 3: Try unmarshalling the JSON into a map[string]any
	var unmarshalled map[string]any
	err = json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&unmarshalled)
	if err != nil {
		return nil, nil, fmt.Errorf("%s.Parse(): %s", p.Name, err.Error())
	}

	// Step 4: Do a breadth-first recursive check for the first instance of
	//         []map[string]any, and we will use this as our rows
	rows, err := p.findJsonObjectList(unmarshalled)
	if err != nil {
		return nil, nil, fmt.Errorf("%s.Parse(): %s", p.Name, err.Error())
	}

	return rows, p.GetHeadings(rows), nil
}

// We have been told (via some parameter) where the root of the []JSONObjects are.
func (p *GenericJsonParser) getJsonObjectList(input io.Reader, rowsRoot string) ([]map[string]any, error) {
	// Use viper, so we get dotted.name.notation for free
	v := viper.New()
	v.SetConfigType("json")
	err := v.ReadConfig(input)
	if err != nil {
		return nil, fmt.Errorf("%s.getJsonObjectList(%s): %s", p.Name, rowsRoot, err.Error())
	}

	if !v.IsSet(rowsRoot) {
		return nil, fmt.Errorf("%s.getJsonObjectList(%s): no []map[string]any found", p.Name, rowsRoot)
	}

	if val, isSliceOfMap := v.Get(rowsRoot).([]map[string]any); isSliceOfMap {
		return val, nil
	}
	if val, isSliceOfAny := v.Get(rowsRoot).([]any); isSliceOfAny {
		if len(val) > 0 {
			for i := range val {
				if _, isMap := val[i].(map[string]any); isMap {
					// Need to convert the []any to []map[string]any
					rows := make([]map[string]any, len(val))
					for j := range val {
						rows[j] = val[j].(map[string]any)
					}
					return rows, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("%s.getJsonObjectList(%s): no []map[string]any found", p.Name, rowsRoot)
}

// depth-first recursive check for the first instance of []map[string]any
func (p *GenericJsonParser) findJsonObjectList(root map[string]any) ([]map[string]any, error) {
	// Get the top-level keys
	for k := range root {
		if ovenReady, isList := root[k].([]map[string]any); isList {
			return ovenReady, nil
		}
		if childMap, isMap := root[k].(map[string]any); isMap {
			rows, err := p.findJsonObjectList(childMap)
			if err == nil && len(rows) > 0 {
				return rows, err
			}
		}
		if childList, isList := root[k].([]any); isList {
			if len(childList) > 0 {
				for i := range childList {
					if _, isMap := childList[i].(map[string]any); isMap {
						// Need to convert the []any to []map[string]any
						rows := make([]map[string]any, len(childList))
						for j := range childList {
							rows[j] = childList[j].(map[string]any)
						}
						return rows, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("%s.findJsonObjectList(): no []map[string]any found", p.Name)
}

func getKeys(root map[string]any) []string {
	keys := []string{}
	for k := range root {
		keys = append(keys, k)
	}

	return keys
}
