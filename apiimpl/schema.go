package apiimpl

import (
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
)

func GetSchema(ctx data.GremelContext, tableName string) (data.Row, error) {
	database := db.GetGremelDB()
	schema, err := database.GetSchema(tableName)
	if err != nil {
		return data.Row{}, err
	}
	return schema, nil
}
