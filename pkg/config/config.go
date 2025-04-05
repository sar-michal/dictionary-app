package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

// Returns Config based on the GO_ENV environment variable.
// By default, it returns the standard database config.
// If GO_ENV is set to "test", it loads the test database config.
// Path of .env files is the project root.
func LoadConfig() (*Config, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to determine current file location")
	}
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	env := os.Getenv("GO_ENV")
	var envFile string
	if env == "test" {
		envFile = filepath.Join(projectRoot, ".env.test")
	} else {
		envFile = filepath.Join(projectRoot, ".env")
	}

	err := godotenv.Load(envFile)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	return cfg, nil
}
