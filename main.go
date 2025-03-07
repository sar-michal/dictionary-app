package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sar-michal/dictionary-app/models"
	"github.com/sar-michal/dictionary-app/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {
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
	models.Migrate(db)
	fmt.Println("Successfully migrated the database")
}
