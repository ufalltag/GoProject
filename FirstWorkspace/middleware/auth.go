package middleware

import (
	"FirstWorkspace/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
		c.Abort()
		return
	}

	token, err := utils.ValidateToken(parts[1], os.Getenv("ACCESS_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен недействителен"})
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Set("email", claims["email"])

	c.Next()
}
