package services

import (
	"FirstWorkspace/database"
	"FirstWorkspace/models"
)

func GetFoldersByUser(userID uint) ([]models.Folder, error) {
	var folders []models.Folder
	result := database.DB.Where("user_id = ?", userID).Find(&folders)
	return folders, result.Error
}

func CreateFolder(name string, userID uint) models.Folder {
	folder := models.Folder{Name: name, UserID: userID}
	database.DB.Create(&folder)
	return folder
}

func UpdateFolder(id string, userID uint, name string) bool {
	result := database.DB.
		Model(&models.Folder{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", name)
	return result.RowsAffected > 0
}

func DeleteFolder(id string, userID uint) bool {
	result := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Folder{})
	return result.RowsAffected > 0
}
