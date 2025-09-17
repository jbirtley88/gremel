package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbirtley88/gremel/apiimpl"
	"github.com/jbirtley88/gremel/data"
)

// GET /api/v1/tables
func Tables(c *gin.Context) {
	ctx := data.NewGremelContext(context.Background())
	if gremelContext, _ := c.Get("gremelcontext"); gremelContext != nil {
		ctx = gremelContext.(data.GremelContext)
	}
	tables, err := apiimpl.GetTables(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tables)
}
