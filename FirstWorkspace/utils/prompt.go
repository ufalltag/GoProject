package utils

import (
	"fmt"
	"strings"
)

func BuildAnalyzePrompt(title, description, body string, folders []string) string {

	var contentParts []string
	if title != "" {
		contentParts = append(contentParts, "Title: "+title)
	}
	if description != "" {
		contentParts = append(contentParts, "Description: "+description)
	}
	if body != "" {
		contentParts = append(contentParts, "Body: "+body)
	}
	content := strings.Join(contentParts, "\n")

	folderList := ""
	if len(folders) > 0 {
		folderList = "Existing folders: " + strings.Join(folders, ", ")
	} else {
		folderList = "No existing folders yet. You MUST set is_new_folder to true."
	}
	return fmt.Sprintf(analyzePrompt, content, folderList)
}
