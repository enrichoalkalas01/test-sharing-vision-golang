package modules

import (
	"time"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/handler"
	middleware "github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var RouteModule = fx.Module("route",
	fx.Invoke(registerMiddlewares),
	fx.Invoke(registerRoutes),
)

func registerMiddlewares(app *fiber.App, log *zap.Logger) {
	// Recovery Middleware
	app.Use(middleware.NewRecoveryMiddleware(log))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Request ID
	app.Use(middleware.NewRequestIDMiddleware().Handle())

	// Fiber native logger middleware
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))

	// Custom Zap logger middleware
	app.Use(zapLoggerMiddleware(log))
}

func zapLoggerMiddleware(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		log.Info("HTTP Request",
			zap.String("framework", "Fiber"),
			zap.String("request_id", c.GetRespHeader("X-Request-ID")),
			zap.String("method", c.Method()),
			zap.String("uri", c.OriginalURL()),
			zap.String("remote_ip", c.IP()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("latency_ms", time.Since(start).Milliseconds()),
		)

		return err
	}
}

func registerRoutes(app *fiber.App, articleHandler *handler.ArticleHandler, config *viper.Viper, log *zap.Logger) {
	apiGroup := app.Group("/api")

	// Root handler
	apiGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running...",
			"version": config.GetString("APP_VERSION"),
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API v1
	apiV1 := apiGroup.Group("/v1")
	apiV1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running...",
			"version": config.GetString("APP_VERSION"),
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Article routes
	articles := apiV1.Group("/article")
	articles.Get("/", articleHandler.GetList)
	articles.Post("/", articleHandler.Create)
	articles.Get("/:article_id", articleHandler.GetDetailByID)
	articles.Put("/:article_id", articleHandler.UpdateByID)
	articles.Delete("/:article_id", articleHandler.DeleteByID)

	log.Info("Routes configured successfully")
}
