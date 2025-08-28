package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MySQLUsername string
	MySQLPassword string
	MySQLAddress  string // host:port
	MySQLDatabase string
	TableName     string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist, just log it
		fmt.Printf("Warning: .env file not found or could not be loaded: %v\n", err)
	}

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		MySQLUsername: os.Getenv("MYSQL_USERNAME"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLAddress:  getEnv("MYSQL_ADDRESS", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", ""),
		TableName:     getEnv("TABLE_NAME", "counters"),
	}

	// 不再强制要求数据库配置，允许应用在没有数据库配置的情况下启动
	if cfg.MySQLUsername == "" || cfg.MySQLAddress == "" || cfg.MySQLDatabase == "" {
		fmt.Printf("Warning: Missing database configuration. Application will start with limited functionality.\n")
		fmt.Printf("Please check .env file for MYSQL_USERNAME, MYSQL_ADDRESS, MYSQL_DATABASE configuration.\n")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
