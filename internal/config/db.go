package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	dsn := cfg.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	} else {
		println("Database connection successful")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to resolve sql.DB from gorm: %v", err))
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db
}
