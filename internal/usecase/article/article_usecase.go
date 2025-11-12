package usecase

import (
	"context"
	"time"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
	dto "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/dto/article"
	repository "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/repository/article"
	"go.uber.org/zap"
)

type articleUsecase struct {
	repoArticle repository.ArticleRepository
	log         *zap.Logger
}

func NewArticleUsecase(repoArticle repository.ArticleRepository, log *zap.Logger) ArticleUsecase {
	return &articleUsecase{
		repoArticle: repoArticle,
		log:         log,
	}
}

func (u *articleUsecase) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	u.log.Info("creating new article", zap.String("title", article.Title))

	// Validate request
	if err := article.Validate(); err != nil {
		u.log.Warn("article validation failed", zap.Error(err))
		return nil, err
	}

	// Check existing data by title
	existing, err := u.repoArticle.GetByTitle(ctx, article.Title)
	if err != nil && err != dto.ErrArticleNotFound {
		u.log.Error("failed to check existing article", zap.Error(err))
		return nil, dto.ErrFailedCreateArticle
	}

	if existing != nil {
		u.log.Warn("article with same title already exists", zap.String("title", article.Title))
		return nil, dto.ErrArticleExists
	}

	// Set timestamps
	article.CreatedDate = time.Now()
	article.UpdatedDate = time.Now()

	if article.Status == "" {
		article.Status = "Draft"
	}

	if err := u.repoArticle.Create(ctx, article); err != nil {
		u.log.Error("failed to save article to database", zap.Error(err))
		return nil, dto.ErrFailedCreateArticle
	}

	u.log.Info("article created successfully", zap.Uint("id", article.ID), zap.String("title", article.Title))

	return article, nil
}

func (u *articleUsecase) GetList(ctx context.Context, filter *dto.ArticleFilter) ([]domain.Article, int64, error) {
	u.log.Info("getting article list",
		zap.Int("page", filter.GetDefaultPage()),
		zap.Int("limit", filter.GetDefaultLimit()),
		zap.String("category", filter.GetCategory()),
		zap.String("status", filter.GetStatus()),
	)

	// Validate filter
	if err := filter.Validate(); err != nil {
		u.log.Warn("filter validation failed", zap.Error(err))
		return nil, 0, err
	}

	// Get data from repository
	articles, total, err := u.repoArticle.GetList(ctx, filter)
	if err != nil {
		u.log.Error("failed to get articles from repository", zap.Error(err))
		return nil, 0, dto.ErrFailedCreateArticle
	}

	u.log.Info("articles retrieved successfully", zap.Int("count", len(articles)), zap.Int64("total", total))

	return articles, total, nil
}

func (u *articleUsecase) GetDetailByID(ctx context.Context, id uint) (*domain.Article, error) {
	u.log.Info("getting article detail", zap.Uint("id", id))

	// Get article by id
	article, err := u.repoArticle.GetDetailByID(ctx, id)
	if err != nil {
		u.log.Warn("article not found", zap.Uint("id", id), zap.Error(err))
		return nil, dto.ErrArticleNotFound
	}

	u.log.Info("article detail retrieved successfully", zap.Uint("id", id))
	return article, nil
}

func (u *articleUsecase) UpdateByID(ctx context.Context, id uint, updateReq *dto.UpdateArticleRequest) (*domain.Article, error) {
	u.log.Info("updating article", zap.Uint("id", id))

	// Validate request
	if err := updateReq.Validate(); err != nil {
		u.log.Warn("update request validation failed", zap.Error(err))
		return nil, err
	}

	// Get existing article data
	article, err := u.repoArticle.GetDetailByID(ctx, id)
	if err != nil {
		u.log.Warn("article not found for update", zap.Uint("id", id), zap.Error(err))
		return nil, dto.ErrArticleNotFound
	}

	if updateReq.Title != "" {
		// Check if new title already exists (but not in this article)
		existing, err := u.repoArticle.GetByTitle(ctx, updateReq.Title)
		if err != nil && err != dto.ErrArticleNotFound {
			u.log.Error("failed to check title availability", zap.Error(err))
			return nil, dto.ErrFailedUpdateArticle
		}
		if existing != nil && existing.ID != id {
			u.log.Warn("title already exists", zap.String("title", updateReq.Title))
			return nil, dto.ErrArticleExists
		}
		article.Title = updateReq.Title
	}

	if updateReq.Content != "" {
		article.Content = updateReq.Content
	}

	if updateReq.Category != "" {
		article.Category = updateReq.Category
	}

	if updateReq.Status != "" {
		article.Status = updateReq.Status
	}

	article.UpdatedDate = time.Now()

	// Update article by id
	if err := u.repoArticle.UpdateByID(ctx, id, article); err != nil {
		u.log.Error("failed to update article", zap.Uint("id", id), zap.Error(err))
		return nil, dto.ErrFailedUpdateArticle
	}

	u.log.Info("article updated successfully", zap.Uint("id", id))
	return article, nil
}

func (u *articleUsecase) DeleteByID(ctx context.Context, id uint) error {
	u.log.Info("deleting article", zap.Uint("id", id))

	// Check exists article
	article, err := u.repoArticle.GetDetailByID(ctx, id)
	if err != nil {
		u.log.Warn("article not found for delete", zap.Uint("id", id), zap.Error(err))
		return dto.ErrArticleNotFound
	}

	// Delete if exist
	if err := u.repoArticle.DeleteByID(ctx, id); err != nil {
		u.log.Error("failed to delete article", zap.Uint("id", id), zap.Error(err))
		return dto.ErrFailedDeleteArticle
	}

	u.log.Info("article deleted successfully", zap.Uint("id", id), zap.String("title", article.Title))
	return nil
}
