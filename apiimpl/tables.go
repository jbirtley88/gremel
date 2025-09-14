package apiimpl

import (
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
)

func GetTables(ctx data.GremelContext) ([]string, error) {
	database := db.GetGremelDB()
	tables, err := database.GetTables()
	if err != nil {
		return nil, err
	}
	return tables, nil
}
