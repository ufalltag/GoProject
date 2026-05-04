package services

import (
	"FirstWorkspace/database"
	"FirstWorkspace/models"
)

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).Find(&user)
	return user, result.Error
}
