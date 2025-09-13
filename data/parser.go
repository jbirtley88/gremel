package data

import (
	"io"
)

// Parser is used to turn arbitray input data info the format required for extracting
// data - which is a map[string]any
//
// This is necessary because not all of the input data we deal with (a lot of the JSON for instance)
// is in a format suitable for extracting data from it. For instance:
//
//		{
//		  "data":
//		     [
//			   "key1": "value1",
//			   "key2": "value2"
//	      ],
//		     [
//			   "key1": "value1",
//			   "key2": "value2"
//	      ]
//		  }
//		}
//
// does not lend itself to being parsed into a map[string]any, so we need to parse it
type Parser interface {
	Parse(input io.Reader) (*RowList, error)
	GetHeadings(rows []Row) []string
	GetName() string
}
