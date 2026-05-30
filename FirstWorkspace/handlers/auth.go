package handlers

import (
	"FirstWorkspace/database"
	"FirstWorkspace/models"
	"FirstWorkspace/utils"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}

	newUser := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Username: input.Username,
	}

	// сохраняем пользователя в бд
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь создан"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	var user models.User

	result := database.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или password"})
		return
	}

	// access token — живёт 15 минут
	accessTokenString, err := utils.GenerateAccessToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
		return
	}

	// refresh token — живёт 7 дней
	refreshTokenString, err := utils.GenerateRefreshToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания jwt токена"})
		return
	}
	database.DB.Model(&user).Update("refresh_token", refreshTokenString)
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	})
}

func ChangePassword(c *gin.Context) {
	email := c.MustGet("email").(string)

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// проверяем, что старый пароль совпадает
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный текущий пароль"})
		return
	}

	// новый пароль не должен совпадать со старым
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.NewPassword)); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Новый пароль должен отличаться от старого"})
		return
	}

	// проверяем сложность нового пароля
	if err := utils.ValidatePasswordStrength(input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}

	database.DB.Model(&user).Update("password", string(hashedPassword))

	c.JSON(http.StatusOK, gin.H{"message": "Пароль изменён"})
}

func Profile(c *gin.Context) {
	email := c.MustGet("email").(string)

	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":    user.Email,
		"username": user.Username,
	})
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Проверка refresh токена
	token, err := utils.ValidateToken(input.RefreshToken, os.Getenv("REFRESH_SECRET"))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен недействителен"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var user models.User
	result := database.DB.Where("email = ? AND refresh_token = ?", email, input.RefreshToken).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен недействителен"})
		return
	}

	// Новый access токен
	newAccessTokenString, err := utils.GenerateAccessToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в создание токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessTokenString})
}
