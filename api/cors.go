package api

import (
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	}
}
