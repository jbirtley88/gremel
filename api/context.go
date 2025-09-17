package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jbirtley88/gremel/data"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO(john): get existing context or inject one if none exists
		ctx := data.NewGremelContext(context.Background())
		c.Set("gremelcontext", ctx)
		c.Next()
	}
}
