package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sar-michal/dictionary-app/pkg/config"
	"github.com/sar-michal/dictionary-app/pkg/models"
	"github.com/sar-michal/dictionary-app/pkg/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {
	os.Setenv("GO_ENV", "development")
	config, err := config.LoadConfig()
	if err != nil {
		log.Println("No .env file found")
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
