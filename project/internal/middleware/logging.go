package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs method, path and Authorization header for incoming requests.
// It intentionally avoids logging request bodies to not leak passwords.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		log.Printf("REQ %s %s Authorization=%s RemoteIP=%s", c.Request.Method, c.Request.URL.Path, auth, c.ClientIP())
		c.Next()
	}
}
