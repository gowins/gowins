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
	Address        string `json:"address" yaml:"address" mapstructure:"address"`
	Mode           string `json:"mode" yaml:"mode" mapstructure:"mode"`
	DownTime       int64  `json:"down_time" yaml:"down_time" mapstructure:"down_time"`
	MaxBodySize    int64  `json:"maxBodySize" yaml:"maxBodySize" mapstructure:"max_body_size"`
	RequestTimeout int64  `json:"requestTimeout" yaml:"requestTimeout" mapstructure:"request_timeout"`
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

	// Setup loggers
	logger.SetupAccLogger(r)

	// Register middleware
	// Set max body size middleware
	if cfg.MaxBodySize > 0 {
		r.Use(middlewares.SetupMaxBodySizeMiddleware(cfg.MaxBodySize))
	}

	r.Use(middlewares.TimeoutMiddleware(time.Duration(cfg.RequestTimeout) * time.Second))

	// Setup routes
	routes.SetupRoutes(r)
	// Set timeout middleware
	//r.Use(middlewares.TimeoutMiddleware(time.Duration(cfg.RequestTimeout) * time.Second))
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
		//ReadTimeout:  time.Duration(cfg.RequestTimeout) * time.Second,
		//WriteTimeout: time.Duration(cfg.RequestTimeout) * time.Second,
	}

	RegisterGraceful(srv, cfg.DownTime)

}
