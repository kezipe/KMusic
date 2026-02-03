package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	headerKey := "Api-Key"
	apiKey := c.GetHeader(headerKey)
	if apiKey == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "authentication method required"})
		return
	}
	if apiKey != os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}
