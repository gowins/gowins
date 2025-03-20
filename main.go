package main

import (
	"log"

	"gowins/config"
	"gowins/logger"
	"gowins/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize Gin engine
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	logger.SetupAccLogger(r)

	// Start server
	if err := r.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
