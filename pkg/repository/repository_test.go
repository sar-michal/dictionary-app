package repository_test

import (
	"log"
	"os"
	"testing"

	"github.com/sar-michal/dictionary-app/pkg/config"
	"github.com/sar-michal/dictionary-app/pkg/models"
	"github.com/sar-michal/dictionary-app/pkg/repository"
	"github.com/sar-michal/dictionary-app/pkg/storage"
)

var repo *repository.GormRepository

func TestMain(m *testing.M) {

	os.Setenv("GO_ENV", "test")
	config, err := config.LoadConfig()
	if err != nil {
		log.Println("No .env file found")
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer storage.CloseDB(db)

	if err := models.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	repo.DB = db
	code := m.Run()
	os.Exit(code)
}
