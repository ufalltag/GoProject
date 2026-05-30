package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
}
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, gin.H{"error": message})
}
