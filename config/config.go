package config

import (
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
}

type Config struct {
	Server ServerConfig `json:"server"`
}

func LoadConfig() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Address: ":8080",
			Mode:    gin.DebugMode,
		},
	}

	// Load configuration from file or environment variables
	// Example: utils.LoadConfigFromFile("config.json", cfg)
	return cfg
}
