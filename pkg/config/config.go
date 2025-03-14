package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Returns Config based on the GO_ENV environment variable.
// By default, it returns the standard database config.
// If GO_ENV is set to "test", it loads the test database config
func LoadConfig() (*Config, error) {
	env := os.Getenv("GO_ENV")
	var err error

	if env == "test" {
		err = godotenv.Load(".env.test")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	return cfg, nil
}
