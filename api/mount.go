package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbirtley88/gremel/apiimpl"
	"github.com/jbirtley88/gremel/data"
)

// PUT /api/v1/mount ? name=xxx & source=yyy
func MountTable(c *gin.Context) {
	table := c.Request.URL.Query().Get("table")
	source := c.Request.URL.Query().Get("source")
	if table == "" || source == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("table and source query parameters are required"))
		return
	}

	ctx := data.NewGremelContext(context.Background())
	if gremelContext, _ := c.Get("gremelcontext"); gremelContext != nil {
		ctx = gremelContext.(data.GremelContext)
	}
	err := apiimpl.Mount(ctx, table, source)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error mounting table: %v", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("mounted '%s' as '%s'", table, source)})
}

// GET /api/v1/mount [? table=xxx]
func GetMount(c *gin.Context) {
	ctx := data.NewGremelContext(context.Background())
	if gremelContext, _ := c.Get("gremelcontext"); gremelContext != nil {
		ctx = gremelContext.(data.GremelContext)
	}

	// If table is empty, return all mounts
	table := c.Request.URL.Query().Get("table")
	tables := []string{}
	if table == "" {
		var err error
		tables, err = apiimpl.GetTables(ctx)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting tables: %v", err))
			return
		}
	} else {
		tables = append(tables, table)
	}

	mounts := make([]data.Row, 0)
	for _, table := range tables {
		mountInfo, err := apiimpl.GetMount(ctx, table)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("error getting mount info: %s", err.Error())})
			return
		}
		mounts = append(mounts, mountInfo)
	}
	c.JSON(http.StatusOK, mounts)
}
