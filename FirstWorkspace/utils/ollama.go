package utils

import (
	"FirstWorkspace/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func AskOllama(prompt string) (string, error) {
	reqBody := dto.OllamaRequest{
		Model:  os.Getenv("OLLAMA_MODEL"),
		Prompt: prompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("Ошибка сериализации %w", err)
	}

	resp, err := http.Post(
		os.Getenv("OLLAMA_URL")+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", fmt.Errorf("Ошибка запроса к Ollama %w", err)
	}
	defer resp.Body.Close()
	var ollamaResp dto.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("Ошибка парсинга ответа %w", err)
	}

	return ollamaResp.Response, nil
}

func ExtractJSON(response string) string {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 {
		return response
	}
	return response[start : end+1]
}
