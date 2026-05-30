package services

import (
	"errors"

	"FirstWorkspace/database"
	"FirstWorkspace/models"
)

// ErrFolderNotFound возвращается, когда папка не найдена или не принадлежит пользователю.
var ErrFolderNotFound = errors.New("папка не найдена")

// FolderWithCount — папка вместе с числом закладок в ней.
type FolderWithCount struct {
	ID             uint
	Name           string
	BookmarksCount int64
}

func GetFoldersByUser(userID uint, limit, offset int) ([]FolderWithCount, int64, error) {
	var total int64
	if err := database.DB.Model(&models.Folder{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var folders []FolderWithCount
	result := database.DB.
		Model(&models.Folder{}).
		Select("folders.id, folders.name, COUNT(bookmarks.id) AS bookmarks_count").
		Joins("LEFT JOIN bookmarks ON bookmarks.folder_id = folders.id AND bookmarks.deleted_at IS NULL").
		Where("folders.user_id = ?", userID).
		Group("folders.id").
		Order("folders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&folders)
	return folders, total, result.Error
}

// GetAllFoldersByUser возвращает все папки пользователя без постраничной выборки.
func GetAllFoldersByUser(userID uint) ([]models.Folder, error) {
	var folders []models.Folder
	result := database.DB.Where("user_id = ?", userID).Find(&folders)
	return folders, result.Error
}

func CreateFolder(name string, userID uint) (models.Folder, error) {
	folder := models.Folder{Name: name, UserID: userID}
	result := database.DB.Create(&folder)
	return folder, result.Error
}

func UpdateFolder(id string, userID uint, name string) error {
	result := database.DB.
		Model(&models.Folder{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", name)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrFolderNotFound
	}
	return nil
}

func DeleteFolder(id string, userID uint) bool {
	result := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Folder{})
	return result.RowsAffected > 0
}
