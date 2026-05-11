package handlers

import (
	"FirstWorkspace/dto"
	"FirstWorkspace/services"
	"FirstWorkspace/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Analyze(c *gin.Context) {

	user, err := services.GetUserByEmail(c.MustGet("email").(string))
	if err != nil {
		utils.NotFound(c, "Пользователь не найден")
		return
	}

	var input dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Неверные данные")
		return
	}

	page, err := utils.ScrapePage(input.URL)
	if err != nil {
		utils.BadRequest(c, "Не удалось загрузить страницу")
		return
	}

	folders, err := services.GetFoldersByUser(user.ID)
	if err != nil {
		utils.InternalError(c, "Ошибка получения папок")
		return
	}

	var folderNames []string
	for _, folder := range folders {
		folderNames = append(folderNames, folder.Name)
	}

	prompt := utils.BuildAnalyzePrompt(page.Title, page.Description, page.Body, folderNames)
	response, err := utils.AskOllama(prompt)
	if err != nil {
		utils.InternalError(c, "Ошибка анализа страницы")
		return
	}

	cleanResponse := utils.ExtractJSON(response)

	var result dto.AnalyzeResult
	if err := json.Unmarshal([]byte(cleanResponse), &result); err != nil {
		utils.InternalError(c, "Ошибка парсинга ответа нейронки")
		return
	}

	isNewFolder := true
	for _, f := range folders {
		if strings.EqualFold(f.Name, result.Folder) {
			isNewFolder = false
			break
		}
	}

	if result.Confidence < 0.7 {
		c.JSON(http.StatusOK, dto.AnalyzeResponse{
			Title:           result.Title,
			URL:             input.URL,
			SuggestedFolder: nil,
			IsNewFolder:     false,
			Confidence:      result.Confidence,
			Message:         "Не удалось определить подходящую папку, выберите вручную",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AnalyzeResponse{
		Title:           result.Title,
		URL:             input.URL,
		SuggestedFolder: &result.Folder,
		IsNewFolder:     isNewFolder,
		Confidence:      result.Confidence,
	})
}
