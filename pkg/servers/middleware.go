package servers

import (
	"time"

	middleware "github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Global Middleware
func (s *FiberServer) SetupMiddlewares() {
	// Recovery Middleware
	s.fiber.Use(middleware.NewRecoveryMiddleware(s.log))

	// Cors
	s.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Request ID - TAMBAHKAN .Handle()
	s.fiber.Use(middleware.NewRequestIDMiddleware().Handle())

	// Fiber native logger middleware
	s.fiber.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))

	// Custom Zap logger middleware
	s.fiber.Use(s.loggerMiddleware())
}

func (s *FiberServer) loggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// Zap structured logging
		s.log.Info("HTTP Request",
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
