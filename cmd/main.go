package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sar-michal/dictionary-app/pkg/models"
	"github.com/sar-michal/dictionary-app/pkg/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file", err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	defer storage.CloseDB(db)

	err = models.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}
	fmt.Println("Successfully migrated the database")
}
