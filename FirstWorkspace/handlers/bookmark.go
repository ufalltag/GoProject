package handlers

import (
	"FirstWorkspace/models"
	"FirstWorkspace/services"
	"FirstWorkspace/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBookmarks(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	pagination := utils.ParsePagination(c)

	bookmarks, total, err := services.GetBookmarksByFolder(c.Param("id"), user.ID, pagination.Limit, pagination.Offset())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	response := make([]models.BookmarkResponse, 0, len(bookmarks))
	for _, b := range bookmarks {
		response = append(response, models.BookmarkResponse{ID: b.ID, URL: b.URL, Title: b.Title, FolderID: b.FolderID})
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": response,
		"meta":      utils.NewPageMeta(pagination, total),
	})
}

const recentBookmarksLimit = 4

func GetRecentBookmarks(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	bookmarks, err := services.GetRecentBookmarks(user.ID, recentBookmarksLimit)
	if err != nil {
		utils.InternalError(c, "Ошибка получения закладок")
		return
	}

	response := make([]models.BookmarkResponse, 0, len(bookmarks))
	for _, b := range bookmarks {
		response = append(response, models.BookmarkResponse{ID: b.ID, URL: b.URL, Title: b.Title, FolderID: b.FolderID})
	}

	c.JSON(http.StatusOK, gin.H{"bookmarks": response})
}

func CreateBookmark(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	var input struct {
		URL      string `json:"url"`
		Title    string `json:"title"`
		FolderID uint   `json:"folder_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Неверные данные")
		return
	}

	if !services.FolderBelongsToUser(input.FolderID, user.ID) {
		utils.NotFound(c, "Папка не найдена")
		return
	}

	bookmark, err := services.CreateBookmark(input.URL, input.Title, input.FolderID, user.ID)
	if err != nil {
		utils.InternalError(c, "Ошибка создания закладки")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"bookmark": models.BookmarkResponse{ID: bookmark.ID, URL: bookmark.URL, Title: bookmark.Title, FolderID: bookmark.FolderID}})
}

func UpdateBookmark(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	var input struct {
		Title    *string `json:"title"`
		FolderID *uint   `json:"folder_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Неверные данные")
		return
	}

	if input.FolderID != nil && !services.FolderBelongsToUser(*input.FolderID, user.ID) {
		utils.NotFound(c, "Папка не найдена")
		return
	}

	if !services.UpdateBookmark(c.Param("id"), user.ID, input.Title, input.FolderID) {
		utils.NotFound(c, "Закладка не найдена")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Закладка обновлена"})
}

func DeleteBookmark(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	if !services.DeleteBookmark(c.Param("id"), user.ID) {
		utils.NotFound(c, "Закладка не найдена")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Закладка удалена"})
}
