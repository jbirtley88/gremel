package apiimpl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jbirtley88/gremel/adapter"
	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
	"github.com/jbirtley88/gremel/helper"
	"github.com/jbirtley88/gremel/util"
)

func Mount(ctx data.GremelContext, name string, source string) error {
	// See if this is a local path
	fi, err := os.Stat(source)
	if err != nil {
		// Not a local path, try URL
		if strings.Contains(source, "://") {
			// This at least looks like a URL
			return MountUrl(ctx, name, source)
		}

		// We don't know what to do with this
		return fmt.Errorf("Mount(%s): unsupported data source", source)
	}

	if fi.IsDir() {
		return fmt.Errorf("Mount(%s): directories not supported", source)
	}
	return nil
}

// Only used for testing - allows us to mock out HTTP calls
var httpHelper = helper.NewHttpHelperBuilder().Build()

func MountUrl(ctx data.GremelContext, name string, sourceUrl string) error {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		return fmt.Errorf("MountUrl(%s): %w", sourceUrl, err)
	}
	if u.Scheme == "file" {
		return MountFile(ctx, name, u.Path)
	}

	// Check that url is http://s or https://
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("MountUrl(%s): unsupported URL scheme '%s'", sourceUrl, u.Scheme)
	}

	// Defer to the HTTP helper to fetch the content
	code, body, err := httpHelper.Get(sourceUrl)
	if err != nil {
		return fmt.Errorf("MountUrl(%s): %w", sourceUrl, err)
	}
	defer body.Close()
	if code != http.StatusOK {
		return fmt.Errorf("MountUrl(%s): GET returned HTTP %d", sourceUrl, code)
	}
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("MountUrl(%s): %w", sourceUrl, err)
	}

	// Now pass it to the JSON adapter for parsing
	// TODO(john): will we ever want to mount over HTTP when the content is not JSON?
	// If so, we need to detect the type somehow and pass it to the appropriate parser
	// For now, we just assume JSON
	err = adapter.CreateTableFromReader(ctx, db.GetGremelDB(), name, bytes.NewBuffer(bodyBytes), adapter.NewGenericJsonParser(ctx))
	if err != nil {
		return fmt.Errorf("MountUrl(%s): %w", sourceUrl, err)
	}
	err = db.GetGremelDB().Mount(name, sourceUrl)
	if err != nil {
		return fmt.Errorf("MountUrl(%s): %w", sourceUrl, err)
	}
	return nil
}

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
