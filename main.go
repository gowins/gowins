package main

import (
	"net/http"
	"time"

	"gowins/conf"
	"gowins/logger"
	"gowins/middlewares"
	"gowins/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address     string `json:"address" yaml:"address" mapstructure:"address"`
	Mode        string `json:"mode" yaml:"mode" mapstructure:"mode"`
	DownTime    int64  `json:"down_time" yaml:"down_time" mapstructure:"down_time"`
	MaxBodySize int64  `json:"maxBodySize" yaml:"maxBodySize" mapstructure:"maxBodySize"`
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

	// Setup loggers
	logger.SetupAccLogger(r)

	// Register middleware
	// Set max body size middleware
	if cfg.MaxBodySize > 0 {
		r.Use(middlewares.SetupMaxBodySizeMiddleware(cfg.MaxBodySize))
	}

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	RegisterGraceful(srv)

}
