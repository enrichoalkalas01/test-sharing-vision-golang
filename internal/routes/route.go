package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Router struct {
	app    *fiber.App
	config *viper.Viper
	log    *zap.Logger
	DB     *gorm.DB
}

func NewRouter(
	app *fiber.App,
	config *viper.Viper,
	log *zap.Logger,
	DB *gorm.DB,
) *Router {
	return &Router{
		app:    app,
		config: config,
		log:    log,
		DB:     DB,
	}
}

// Simple Routes
func (r *Router) SetupRoutes(api fiber.Router, config *viper.Viper, log *zap.Logger, DB *gorm.DB) {
	// Root Endpoint
	api.Get("/", r.RootHandler)

	// Api V1
	apiV1 := api.Group("/v1")
	apiV1.Get("/", r.RootHandler)

	ArticleRoutes(apiV1, r.DB, r.log)
}

func (r *Router) RootHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Server is running...",
		"version": r.config.GetString("APP_VERSION"),
		"time":    time.Now().Format(time.RFC3339),
	})
}
