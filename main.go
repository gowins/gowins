package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	logger.SetupAccLogger(r)

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

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	fmt.Println(cfg.DownTime)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DownTime)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
