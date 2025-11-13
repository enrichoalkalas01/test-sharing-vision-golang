package repository

import (
	"context"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
	dto "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/dto/article"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type articleRepository struct {
	DB  *gorm.DB
	log *zap.Logger
}

func NewArticleRepository(DB *gorm.DB, log *zap.Logger) ArticleRepository {
	return &articleRepository{
		DB:  DB,
		log: log,
	}
}

func (r *articleRepository) Create(ctx context.Context, article *domain.Article) error {
	r.log.Debug("repository: creating article", zap.String("title", article.Title))

	if err := r.DB.WithContext(ctx).Table("posts").Create(article).Error; err != nil {
		r.log.Error("repository: faield to create article", zap.Error(err))
		return err
	}

	r.log.Debug("repository: article created successfully", zap.Uint("id", article.ID))

	return nil
}

func (r *articleRepository) GetList(ctx context.Context, articleFilter *dto.ArticleFilter) ([]domain.Article, int64, error) {
	r.log.Debug("repository: getting article list", zap.Int("page", articleFilter.Page), zap.Int("limit", articleFilter.Limit))

	var articles []domain.Article
	var total int64

	// Build query
	query := r.DB.WithContext(ctx).Table("posts")

	// Apply filters into query
	if articleFilter.HasCategory() {
		query = query.Where("category = ?", articleFilter.GetCategory())
	}

	if articleFilter.HasStatus() {
		query = query.Where("status = ?", articleFilter.GetStatus())
	}

	if articleFilter.GetSearch() != "" {
		query = query.Where("title LIKE ?", "%"+articleFilter.Search+"%")
	}

	// Get total count data
	if err := query.Model(&domain.Article{}).Count(&total).Error; err != nil {
		r.log.Error("repository: failed to count articles", zap.Error(err))
		return nil, 0, err
	}

	// Apply sorting
	sortBy := articleFilter.GetSortBy()
	if sortBy == "" {
		sortBy = "created_date"
	}

	sortOrder := articleFilter.GetSortOrder()
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	offset := articleFilter.GetOffset()
	limit := articleFilter.GetDefaultLimit()
	query = query.Offset(offset).Limit(limit)

	// Execute query
	if err := query.Find(&articles).Error; err != nil {
		r.log.Error("repository: failed to get articles", zap.Error(err))
		return nil, 0, nil
	}

	r.log.Debug("repository: articles retrieved successfully", zap.Int("count", len(articles)))

	return articles, total, nil
}

func (r *articleRepository) GetByTitle(ctx context.Context, title string) (*domain.Article, error) {
	r.log.Debug("repository: getting article by title", zap.String("title", title))

	var article domain.Article
	if err := r.DB.WithContext(ctx).Table("posts").Where("title = ?", title).First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug("repository: article not found by title", zap.String("title", title))
			return nil, dto.ErrArticleNotFound
		}
		r.log.Error("repository: failed to get article by title", zap.String("title", title), zap.Error(err))
		return nil, err
	}

	r.log.Debug("repository: article found by title", zap.String("title", title), zap.Uint("id", article.ID))
	return &article, nil
}

func (r *articleRepository) GetDetailByID(ctx context.Context, id uint) (*domain.Article, error) {
	var article domain.Article
	if err := r.DB.WithContext(ctx).Table("posts").First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warn("repository: article not found", zap.Uint("id", id))
			return nil, dto.ErrArticleNotFound
		}
		r.log.Error("repository: failed to get article detail", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	r.log.Debug("repository: article detail retrieved successfully", zap.Uint("id", id))
	return &article, nil
}

func (r *articleRepository) UpdateByID(ctx context.Context, id uint, article *domain.Article) error {
	r.log.Debug("repository: updating article", zap.Uint("id", id))

	// Update hanya field-field yang berubah
	if err := r.DB.WithContext(ctx).Table("posts").Where("id = ?", id).Model(&domain.Article{}).Updates(article).Error; err != nil {
		r.log.Error("repository: failed to update article", zap.Uint("id", id), zap.Error(err))
		return err
	}

	r.log.Debug("repository: article updated successfully", zap.Uint("id", id))
	return nil
}

func (r *articleRepository) DeleteByID(ctx context.Context, id uint) error {
	r.log.Debug("repository: deleting article", zap.Uint("id", id))

	if err := r.DB.WithContext(ctx).Table("posts").Where("id = ?", id).Delete(&domain.Article{}, id).Error; err != nil {
		r.log.Error("repository: failed to delete article", zap.Uint("id", id), zap.Error(err))
		return err
	}

	r.log.Debug("repository: article deleted successfully", zap.Uint("id", id))
	return nil
}
