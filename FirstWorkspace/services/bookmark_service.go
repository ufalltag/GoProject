package services

import (
	"FirstWorkspace/database"
	"FirstWorkspace/models"
	"strings"
)

func GetBookmarksByFolder(folderID string, userID uint, limit, offset int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	query := database.DB.Model(&models.Bookmark{}).Where("folder_id = ? AND user_id = ?", folderID, userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookmarks)
	return bookmarks, total, result.Error
}

// GetRecentBookmarks возвращает последние добавленные закладки пользователя из всех папок.
func GetRecentBookmarks(userID uint, limit int) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	result := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&bookmarks)
	return bookmarks, result.Error
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

func UpdateBookmark(id string, userID uint, title *string, folderID *uint) bool {
	updates := map[string]interface{}{}
	if title != nil {
		updates["title"] = *title
	}
	if folderID != nil {
		updates["folder_id"] = *folderID
	}
	if len(updates) == 0 {
		return false
	}

	result := database.DB.
		Model(&models.Bookmark{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)
	return result.RowsAffected > 0
}

func DeleteBookmark(id string, userID uint) bool {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Bookmark{})
	return result.RowsAffected > 0
}

func SearchBookmarks(userID uint, query string, limit, offset int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	q := database.DB.Model(&models.Bookmark{}).
		Where("user_id = ? AND lower(title) LIKE ?", userID, strings.ToLower(query)+"%")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&bookmarks)
	return bookmarks, total, result.Error
}

func FolderBelongsToUser(folderID, userID uint) bool {
	var folder models.Folder
	result := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder)
	return result.Error == nil
}
