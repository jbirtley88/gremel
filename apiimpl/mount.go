package apiimpl

import (
	"fmt"

	"github.com/jbirtley88/gremel/adapter"
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/jbirtley88/gremel/util"
)

func MountFile(ctx data.GremelContext, path string) error {
	name, ext, err := util.SplitFilename(path)
	if err != nil {
		return fmt.Errorf("MountFile(%s): %w", path, err)
	}
	database := db.GetGremelDB()
	err = adapter.CreateTableFromFile(ctx, database, name, ext, path)
	return err
}
