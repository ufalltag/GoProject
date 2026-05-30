package handlers

import (
	"FirstWorkspace/models"
	"FirstWorkspace/services"
	"FirstWorkspace/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetFolders(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	pagination := utils.ParsePagination(c)

	folders, total, err := services.GetFoldersByUser(user.ID, pagination.Limit, pagination.Offset())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	response := make([]models.FolderResponse, 0, len(folders))
	for _, folder := range folders {
		response = append(response, models.FolderResponse{
			ID:             folder.ID,
			Name:           folder.Name,
			BookmarksCount: folder.BookmarksCount,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"folders": response,
		"meta":    utils.NewPageMeta(pagination, total),
	})
}

func CreateFolder(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Неверные данные")
		return
	}

	folder, err := services.CreateFolder(input.Name, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.Conflict(c, "Папка с таким именем уже существует")
			return
		}
		utils.InternalError(c, "Ошибка создания папки")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"folder": models.FolderResponse{ID: folder.ID, Name: folder.Name}})
}

func DeleteFolder(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	if !services.DeleteFolder(c.Param("id"), user.ID) {
		utils.NotFound(c, "Папка не найдена")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Папка удалена"})
}

func UpdateFolder(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Неверные данные")
		return
	}

	if err := services.UpdateFolder(c.Param("id"), user.ID, input.Name); err != nil {
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			utils.Conflict(c, "Папка с таким именем уже существует")
		case errors.Is(err, services.ErrFolderNotFound):
			utils.NotFound(c, "Папка не найдена")
		default:
			utils.InternalError(c, "Ошибка обновления папки")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Папка обновлена"})
}
