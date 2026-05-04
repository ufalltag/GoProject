package handlers

import (
	"FirstWorkspace/models"
	"FirstWorkspace/services"
	"FirstWorkspace/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFolders(c *gin.Context) {
	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	folders, err := services.GetFoldersByUser(user.ID)

	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	var response []models.FolderResponse
	for _, folder := range folders {
		response = append(response, models.FolderResponse{
			ID:   folder.ID,
			Name: folder.Name,
		})
	}
	c.JSON(http.StatusOK, gin.H{"folders": response})
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

	folder := services.CreateFolder(input.Name, user.ID)
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

	if !services.UpdateFolder(c.Param("id"), user.ID, input.Name) {
		utils.NotFound(c, "Папка не найдена")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Папка обновлена"})
}
