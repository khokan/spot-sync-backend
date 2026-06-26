package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {
	// Try loading .env from current directory, fallback to parent directory if not found
	if err := godotenv.Load(); err != nil {
		if err = godotenv.Load("../.env"); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	cfg := &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       os.Getenv("DSN"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.Dsn == "" || cfg.JwtSecret == "" {
		log.Fatal("Error: DSN or JWT_SECRET environment variables are not set")
	}

	return cfg
}
