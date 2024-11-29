package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create a Chat function that will handle the chat requests
func Chat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
