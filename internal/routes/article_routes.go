package routes

import (
	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/handler"
	repository "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/repository/article"
	usecase "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/usecase/article"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ArticleRoutes(
	router fiber.Router,
	DB *gorm.DB,
	log *zap.Logger,
) {
	// Depedency Injection
	articleRepo := repository.NewArticleRepository(DB, log)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, log)
	articleHandler := handler.NewArticleHandler(articleUsecase, log)

	// Routes
	articles := router.Group("/article")

	articles.Get("/", articleHandler.GetList)
	articles.Post("/", articleHandler.Create)
	articles.Get("/:article_id", articleHandler.GetDetailByID)
	articles.Put("/:article_id", articleHandler.UpdateByID)
	articles.Delete("/:article_id", articleHandler.DeleteByID)

}
