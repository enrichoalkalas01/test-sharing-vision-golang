package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ServerModule = fx.Module("server",
	fx.Provide(NewFiberApp),
	fx.Invoke(registerServerLifecycle),
)

func NewFiberApp(config *viper.Viper) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:               config.GetString("APP_NAME"),
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		DisableStartupMessage: false,
	})
}

func registerServerLifecycle(lc fx.Lifecycle, app *fiber.App, config *viper.Viper, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			host := config.GetString("HOST")
			if host == "" {
				host = "0.0.0.0"
			}

			port := config.GetString("PORT")
			if port == "" {
				port = "8855"
			}

			address := fmt.Sprintf("%s:%s", host, port)

			log.Info("Starting fiber server",
				zap.String("framework", "fiber"),
				zap.String("host", host),
				zap.String("port", port),
				zap.String("address", address),
				zap.String("version", config.GetString("APP_VERSION")),
			)

			go func() {
				if err := app.Listen(address); err != nil {
					log.Error("Fiber server error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down Fiber server...")
			return app.ShutdownWithContext(ctx)
		},
	})
}
