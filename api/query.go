package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbirtley88/gremel/apiimpl"
	"github.com/jbirtley88/gremel/data"
)

// GET /api/v1/query ? q=xxxxxx
func Query(c *gin.Context) {
	ctx := data.NewGremelContext(context.Background())
	if gremelContext, _ := c.Get("gremelcontext"); gremelContext != nil {
		ctx = gremelContext.(data.GremelContext)
	}

	// If query is empty, return an error
	if c.Request.URL.Query().Get("q") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q= query parameter is required"})
		return
	}

	rows, headings, err := apiimpl.Query(ctx, c.Request.URL.Query().Get("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error executing query: %v", err)})
		return
	}

	// Return the query results as JSON
	c.JSON(http.StatusOK, gin.H{
		"rows":     rows,
		"headings": headings,
	})
}
