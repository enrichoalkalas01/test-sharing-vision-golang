package repository

import (
	"context"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
	dto "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/dto/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) error
	GetList(ctx context.Context, articleFilter *dto.ArticleFilter) ([]domain.Article, int64, error)
	GetByTitle(ctx context.Context, title string) (*domain.Article, error)
	GetDetailByID(ctx context.Context, id uint) (*domain.Article, error)
	UpdateByID(ctx context.Context, id uint, article *domain.Article) error
	DeleteByID(ctx context.Context, id uint) error
}
