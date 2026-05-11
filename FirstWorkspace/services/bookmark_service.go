package services

import (
	"FirstWorkspace/database"
	"FirstWorkspace/models"
)

func GetBookmarksByFolder(folderID string, userID uint) []models.Bookmark {
	var bookmarks []models.Bookmark
	database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).Find(&bookmarks)
	return bookmarks
}

func CreateBookmark(url, title string, folderID, userID uint) (models.Bookmark, error) {
	bookmark := models.Bookmark{
		URL:      url,
		Title:    title,
		FolderID: folderID,
		UserID:   userID,
	}
	result := database.DB.Create(&bookmark)
	return bookmark, result.Error
}

func UpdateBookmark(id string, userID uint, title string, folderID uint) bool {
	result := database.DB.
		Model(&models.Bookmark{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"title":     title,
			"folder_id": folderID,
		})
	return result.RowsAffected > 0
}

func DeleteBookmark(id string, userID uint) bool {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Bookmark{})
	return result.RowsAffected > 0
}

func FolderBelongsToUser(folderID, userID uint) bool {
	var folder models.Folder
	result := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder)
	return result.Error == nil
}
