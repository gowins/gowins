package main

import (
	"log"

	"gowins/conf"
	"gowins/logger"
	"gowins/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string `json:"address" yaml:"address"`
	Mode    string `json:"mode" yaml:"mode"`
}

func main() {
	// Load configuration
	conf.LoadConfig()

	cfg := &ServerConfig{}
	if err := viper.Sub("server").Unmarshal(cfg); err != nil {
		panic("Failed to parse config: " + err.Error())
	}

	// Set Gin mode
	gin.SetMode(cfg.Mode)

	// Initialize Gin engine
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	logger.SetupAccLogger(r)

	// Start server
	if err := r.Run(cfg.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
