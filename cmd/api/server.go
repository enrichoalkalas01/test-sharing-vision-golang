package main

import (
	"os"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/configs"
	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/database"
	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/logger"
	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/servers"
	"go.uber.org/zap"
)

func main() {
	// Step 1
	log, err := logger.NewLogger(logger.Config{
		Environtment: getEnv("APP_ENV", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		OutputPath:   "stdout",
	})

	if err != nil {
		panic("Failed to init logger: " + err.Error())
	}

	defer log.Sync()

	log.Info("Starting application...")

	// Step 2
	configEnv, err := configs.NewViper(".env", "env", ".", "../../")
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
		panic(err)
	}

	log.Info("ENV configuration loaded successfully",
		zap.String("app_name", configEnv.GetString("APP_NAME")),
		zap.String("app_env", configEnv.GetString("APP_ENV")),
		zap.String("app_version", configEnv.GetString("APP_VERSION")),
	)

	// Step 3 Init Mysql
	mysql, err := database.NewMySQL(configEnv, log.Logger)
	if err != nil {
		log.Fatal("failed to init mysql database", zap.Error(err))
		panic(err)
	}

	log.Info("mysql connected successfully")

	// Init Echo Server
	server := servers.NewFiberServer(configEnv, log.Logger, mysql)
	server.SetupMiddlewares()
	server.SetupRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

// Helper Get Env
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
