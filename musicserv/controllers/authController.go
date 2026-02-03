package controllers

import "github.com/gin-gonic/gin"

func AmIAuthorized(c *gin.Context) {
	c.JSON(200, gin.H{"message": "You are authorized"})
}
