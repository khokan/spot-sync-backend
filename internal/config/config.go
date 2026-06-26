package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	Dsn            string
	JwtSecret      string
	JWTExpiryHours int
	BcryptCost     int
}

func LoadEnv() *Config {
	// Try loading .env from current directory and fall back to parent for cmd execution.
	if err := godotenv.Load(); err != nil {
		if err = godotenv.Load("../.env"); err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	jwtExpiryHours := getIntEnv("JWT_EXPIRY_HOURS", 24)
	bcryptCost := getIntEnv("BCRYPT_COST", 10)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:           port,
		Dsn:            os.Getenv("DSN"),
		JwtSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: jwtExpiryHours,
		BcryptCost:     bcryptCost,
	}

	if cfg.Dsn == "" || cfg.JwtSecret == "" {
		log.Fatal("Error: DSN or JWT_SECRET environment variables are not set")
	}

	if cfg.BcryptCost < 10 || cfg.BcryptCost > 12 {
		log.Fatal("Error: BCRYPT_COST must be between 10 and 12")
	}

	if cfg.JWTExpiryHours <= 0 {
		log.Fatal("Error: JWT_EXPIRY_HOURS must be greater than 0")
	}

	return cfg
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
