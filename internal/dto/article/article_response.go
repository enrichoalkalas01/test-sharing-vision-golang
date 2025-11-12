package dto

import (
	"time"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
)

type ArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleListResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ============ Mapper Functions ============
func ToArticleResponse(article *domain.Article) *ArticleResponse {
	if article == nil {
		return nil
	}

	return &ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Category:  article.Category,
		Status:    article.Status,
		CreatedAt: article.CreatedDate,
		UpdatedAt: article.UpdatedDate,
	}
}

func ToArticleListResponse(article *domain.Article) *ArticleListResponse {
	if article == nil {
		return nil
	}

	return &ArticleListResponse{
		ID:        article.ID,
		Title:     article.Title,
		Category:  article.Category,
		Status:    article.Status,
		CreatedAt: article.CreatedDate,
		UpdatedAt: article.UpdatedDate,
	}
}

func ToArticleResponseList(articles []domain.Article) []ArticleListResponse {
	responses := make([]ArticleListResponse, len(articles))
	for i, article := range articles {
		responses[i] = *ToArticleListResponse(&article)
	}
	return responses
}
