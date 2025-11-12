package servers

import (
	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/routes"
)

func (s *FiberServer) SetupRoutes() {
	// Setup Routes group
	apiGroup := s.fiber.Group("/api")

	// Initialize router
	router := routes.NewRouter(s.fiber, s.config, s.log, s.DB)

	// Setup all routes
	router.SetupRoutes(apiGroup, s.config, s.log, s.DB)

	s.log.Info("Routes configured successfully")
}
