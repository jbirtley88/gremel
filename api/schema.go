package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbirtley88/gremel/apiimpl"
	"github.com/jbirtley88/gremel/data"
)

// GET /api/v1/schema ? table=xxx
func Schema(c *gin.Context) {
	ctx := data.NewGremelContext(context.Background())
	if gremelContext, _ := c.Get("gremelcontext"); gremelContext != nil {
		ctx = gremelContext.(data.GremelContext)
	}

	// If table is empty, return an error
	if c.Request.URL.Query().Get("table") == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("table query parameter is required"))
		return
	}

	schema, err := apiimpl.GetSchema(ctx, c.Request.URL.Query().Get("table"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schema)
}
