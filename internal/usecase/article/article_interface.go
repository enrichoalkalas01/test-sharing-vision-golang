package usecase

import (
	"context"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
	dto "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/dto/article"
)

type ArticleUsecase interface {
	Create(ctx context.Context, article *domain.Article) (*domain.Article, error)
	GetList(ctx context.Context, filter *dto.ArticleFilter) ([]domain.Article, int64, error)
	GetDetailByID(ctx context.Context, id uint) (*domain.Article, error)
	UpdateByID(ctx context.Context, id uint, updateReq *dto.UpdateArticleRequest) (*domain.Article, error)
	DeleteByID(ctx context.Context, id uint) error
}
