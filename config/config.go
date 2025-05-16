package config

import (
	"os"

	"app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
}

func NewConfig() (*Config, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=postgres password=secret dbname=userapp port=5432 sslmode=disable"
	}
	return &Config{DatabaseURL: dsn}, nil
}

func NewDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	return db, nil
}
