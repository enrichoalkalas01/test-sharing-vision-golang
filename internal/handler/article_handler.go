package handler

import (
	"strconv"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/domain"
	dto "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/dto/article"
	usecase "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/usecase/article"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/common/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	articleUsecase usecase.ArticleUsecase
	log            *zap.Logger
}

func NewArticleHandler(
	articleUsecase usecase.ArticleUsecase,
	log *zap.Logger,
) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: articleUsecase,
		log:            log,
	}
}

func (h *ArticleHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateArticleRequest

	// Parse Request Body
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Error("failed to parse create article request", zap.Error(err))
		errResponse := response.NewErrorResponseWithPath(
			"Invalid request body",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)

		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	// Validate Request
	if err := req.Validate(); err != nil {
		h.log.Warn("validation error on create article", zap.Error(err))
		errResponse := response.NewErrorResponseWithDetails(
			err.Error(),
			string(dto.MapErrorToCode(err)),
			map[string]any{
				"field": "create_article",
			},
		)
		errResponse.Path = ctx.Path()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	// Convert to entity / domain
	article := toDomainArticle(&req)

	// Create article
	createdArticle, err := h.articleUsecase.Create(ctx.Context(), article)
	if err != nil {
		h.log.Error("failed to create article", zap.Error(err))
		errResp := response.NewErrorResponseWithPath(
			"Failed to create article",
			string(dto.MapErrorToCode(err)),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	// Convert response
	articleResponse := dto.ToArticleResponse(createdArticle)

	resp := response.NewSuccessResponseWithPath(
		articleResponse,
		"Article created successfully",
		ctx.Path(),
	)
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *ArticleHandler) GetList(ctx *fiber.Ctx) error {
	// Parse filter from query param
	articleFilter := dto.NewArticleFilter()
	if err := ctx.QueryParser(articleFilter); err != nil {
		h.log.Error("faield to parse query param", zap.Error(err))
		errResponse := response.NewErrorResponseWithPath(
			"Invalid query parameters",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	category := ctx.Query("category")
	status := ctx.Query("status")

	if category != "" {
		articleFilter.Filters.Category = category
	}

	if status != "" {
		articleFilter.Filters.Status = status
	}

	h.log.Info("Parsed filter values",
		zap.String("category", articleFilter.GetCategory()),
		zap.String("status", articleFilter.GetStatus()),
	)

	// Validate filter
	if err := articleFilter.Validate(); err != nil {
		h.log.Warn("validation error on filter", zap.Error(err))
		errResponse := response.NewErrorResponseWithPath(
			err.Error(),
			string(dto.MapErrorToCode(err)),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	// Get article data
	articles, total, err := h.articleUsecase.GetList(ctx.Context(), articleFilter)
	if err != nil {
		h.log.Error("faield to get articles", zap.Error(err))
		errResponse := response.NewErrorResponseWithPath(
			"failed to get articles",
			string(dto.ErrCodeDBError),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	// Convert response
	articleResponses := dto.ToArticleResponseList(articles)

	// calculate pagination metadata
	paginationMeta := response.CalculatePaginationMeta(
		articleFilter.GetDefaultPage(),
		articleFilter.GetDefaultLimit(),
		total,
	)

	// Return data passing
	responseHandler := response.NewPaginatedResponseWithPath(
		articleResponses,
		"Articles retrieved successfully",
		ctx.Path(),
		paginationMeta,
	)

	return ctx.Status(fiber.StatusOK).JSON(responseHandler)
}

func (h *ArticleHandler) GetDetailByID(ctx *fiber.Ctx) error {
	// Parse string article_id into uint id db
	articleIDStr := ctx.Params("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		h.log.Warn("invalid article id format", zap.Error(err), zap.String("id", articleIDStr))
		errResp := response.NewErrorResponseWithPath(
			"Invalid article ID format",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// Get by id
	article, err := h.articleUsecase.GetDetailByID(ctx.Context(), uint(articleID))
	if err != nil {
		h.log.Error("failed to get article detail", zap.Error(err), zap.Uint("id", uint(articleID)))
		errResp := response.NewErrorResponseWithPath(
			"Article not found",
			string(dto.ErrCodeNotFound),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusNotFound).JSON(errResp)
	}

	// Convert response
	articleResponse := dto.ToArticleResponse(article)

	resp := response.NewSuccessResponseWithPath(
		articleResponse,
		"Article retrieved successfully",
		ctx.Path(),
	)

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h *ArticleHandler) UpdateByID(ctx *fiber.Ctx) error {
	// Parse string article_id into uint id db
	articleIDStr := ctx.Params("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		h.log.Warn("invalid article id format", zap.Error(err), zap.String("id", articleIDStr))
		errResp := response.NewErrorResponseWithPath(
			"Invalid article ID format",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// Parse request body
	var req dto.UpdateArticleRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Error("failed to parse update article request", zap.Error(err))
		errResp := response.NewErrorResponseWithPath(
			"Invalid request body",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.log.Warn("validation error on update article", zap.Error(err))
		errResp := response.NewErrorResponseWithDetails(
			err.Error(),
			string(dto.MapErrorToCode(err)),
			map[string]any{
				"field": "update_article",
			},
		)
		errResp.Path = ctx.Path()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// Update article
	updatedArticle, err := h.articleUsecase.UpdateByID(ctx.Context(), uint(articleID), &req)
	if err != nil {
		h.log.Error("failed to update article", zap.Error(err), zap.Uint("id", uint(articleID)))
		errResp := response.NewErrorResponseWithPath(
			"Failed to update article",
			string(dto.MapErrorToCode(err)),
			ctx.Path(),
		)
		statusCode := fiber.StatusInternalServerError
		if err == dto.ErrArticleNotFound {
			statusCode = fiber.StatusNotFound
		}
		return ctx.Status(statusCode).JSON(errResp)
	}

	// Convert to response
	articleResponse := dto.ToArticleResponse(updatedArticle)

	resp := response.NewSuccessResponseWithPath(
		articleResponse,
		"Article updated successfully",
		ctx.Path(),
	)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h *ArticleHandler) DeleteByID(ctx *fiber.Ctx) error {
	// Parse string article_id into uint id db
	articleIDStr := ctx.Params("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		h.log.Warn("invalid article id format", zap.Error(err), zap.String("id", articleIDStr))
		errResp := response.NewErrorResponseWithPath(
			"Invalid article ID format",
			string(dto.ErrCodeValidation),
			ctx.Path(),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// Delete by id
	if err := h.articleUsecase.DeleteByID(ctx.Context(), uint(articleID)); err != nil {
		h.log.Error("failed to delete article", zap.Error(err), zap.Uint("id", uint(articleID)))
		errResp := response.NewErrorResponseWithPath(
			"Failed to delete article",
			string(dto.MapErrorToCode(err)),
			ctx.Path(),
		)
		statusCode := fiber.StatusInternalServerError
		if err == dto.ErrArticleNotFound {
			statusCode = fiber.StatusNotFound
		}
		return ctx.Status(statusCode).JSON(errResp)
	}

	resp := response.NewSuccessResponseWithPath(
		"",
		"Article deleted successfully",
		ctx.Path(),
	)

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// === Helper Handler ===
func toDomainArticle(req *dto.CreateArticleRequest) *domain.Article {
	return &domain.Article{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}
}
