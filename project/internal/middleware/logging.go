package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		log.Printf("REQ %s %s Authorization=%s RemoteIP=%s", c.Request.Method, c.Request.URL.Path, auth, c.ClientIP())
		c.Next()
	}
}
