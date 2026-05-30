package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 20
)

// Pagination содержит параметры постраничной выборки.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Offset возвращает смещение для запроса в БД.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PageMeta описывает метаданные постраничного ответа.
type PageMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ParsePagination читает query-параметры page и limit, ограничивая limit MaxPageSize.
func ParsePagination(c *gin.Context) Pagination {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	return Pagination{Page: page, Limit: limit}
}

// NewPageMeta вычисляет метаданные на основе общего числа объектов.
func NewPageMeta(p Pagination, total int64) PageMeta {
	totalPages := int((total + int64(p.Limit) - 1) / int64(p.Limit))
	return PageMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
