package apiimpl

import (
	"fmt"

	"github.com/jbirtley88/gremel/adapter"
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/jbirtley88/gremel/util"
)

func MountFile(ctx data.GremelContext, name, path string) error {
	_, ext, err := util.SplitFilename(path)
	if err != nil {
		return fmt.Errorf("MountFile(%s): %w", path, err)
	}
	database := db.GetGremelDB()
	err = adapter.CreateTableFromFile(ctx, database, name, ext, path)
	if err != nil {
		return fmt.Errorf("MountFile(%s): %w", path, err)
	}
	err = database.Mount(name, path)
	if err != nil {
		return fmt.Errorf("MountFile(%s): %w", path, err)
	}
	return nil
}

func GetMount(ctx data.GremelContext, tableName string) (data.Row, error) {
	database := db.GetGremelDB()
	mountInfo, err := database.GetMount(tableName)
	if err != nil {
		return data.Row{}, fmt.Errorf("GetMount(%s): %w", tableName, err)
	}
	return mountInfo, nil
}
