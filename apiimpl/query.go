package apiimpl

import (
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
)

// sqlQuery MUST be sanitized before calling this function
// cf. Bobby Tables: https://xkcd.com/327/
func Query(ctx data.GremelContext, sqlQuery string) ([]data.Row, []string, error) {
	database := db.GetGremelDB()
	rows, columns, err := database.Query(sqlQuery)
	if err != nil {
		return nil, nil, err
	}
	return rows, columns, nil
}
