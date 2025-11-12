package servers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ===== Setup Servers =====
type FiberServer struct {
	fiber  *fiber.App
	config *viper.Viper
	log    *zap.Logger

	DB *gorm.DB
}

func NewFiberServer(config *viper.Viper, log *zap.Logger, DB *gorm.DB) *FiberServer {
	app := fiber.New(fiber.Config{
		AppName:               config.GetString("APP_NAME"),
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		DisableStartupMessage: false,
	})

	return &FiberServer{
		fiber:  app,
		config: config,
		log:    log,
		DB:     DB,
	}
}

func (s *FiberServer) Start() error {
	// Get host and port with fallback
	host := s.config.GetString("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := s.config.GetString("PORT")
	if port == "" {
		port = "8855"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	s.log.Info("Starting fiber server",
		zap.String("framework", "fiber"),
		zap.String("host", host),
		zap.String("port", port),
		zap.String("address", address),
		zap.String("version", s.config.GetString("APP_VERSION")),
	)

	// Fiber native log
	log.Printf("⚡ Fiber server starting on http://%s", address)

	// Manage server lifecycle
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server in go routine
	go func() {
		if err := s.fiber.Listen(address); err != nil {
			s.log.Error("Fiber server startup error", zap.Error(err))
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	s.log.Info("Fiber server started successfully",
		zap.String("address", address),
		zap.String("url", fmt.Sprintf("http://%s:%s", host, port)),
	)

	// Wait for shutdown signal
	<-ctx.Done()

	// Gracefull shutdown with configrable timeout
	shutdownTimeout := s.config.GetDuration("TIMEOUT_GRACEFUL_SHUTDOWN")
	if shutdownTimeout == 0 {
		shutdownTimeout = 10 // default 10 seconds
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	return s.Shutdown(shutdownCtx)
}

func (s *FiberServer) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down Fiber server...")
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.fiber.ShutdownWithContext(shutdownCtx); err != nil {
		s.log.Error("Fiber server forced to shutdown", zap.Error(err))
		log.Fatalf("Server shutdown error: %v", err)
		return err
	}

	s.log.Info("Fiber server exited gracefully")
	log.Println("Server exited gracefully")
	return nil
}
