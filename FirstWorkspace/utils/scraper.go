package utils

import (
	"FirstWorkspace/dto"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapePage(url string) (dto.PageContent, error) {
	resp, err := http.Get(url)
	if err != nil {
		return dto.PageContent{}, fmt.Errorf("Не удалось загрузить страницу: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return dto.PageContent{}, fmt.Errorf("Не удалось распарсить HTML: %w", err)
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	description, _ := doc.Find("meta[name='description']").Attr("content")
	description = strings.TrimSpace(description)

	if title != "" && description != "" {
		return dto.PageContent{
			Title:       title,
			Description: description,
			Body:        "",
		}, nil
	}

	var bodyParts []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if len(text) > 50 {
			bodyParts = append(bodyParts, text)
		}
	})
	body := strings.Join(bodyParts, " ")
	if len(body) > 1000 {
		body = body[:1000]
	}

	return dto.PageContent{
		Title:       title,
		Description: description,
		Body:        body,
	}, nil
}
