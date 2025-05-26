package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const PROCESS_TIMEOUT = 30 * time.Second
const MAX_MEMORY = 10 * 1024 * 1024 // 10MB
const MAX_PROCESSES = 5
const CHECK_IMAGES_INTERVAL = 10 * time.Minute
const RATE_LIMIT_REQ_S = 1
const RATE_LIMIT_BURST_S = 2

type Config struct {
	PORT   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func GetConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("configuration error: PORT environment variable is missing")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	missing := []string{}
	if dbHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if dbPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if dbUser == "" {
		missing = append(missing, "DB_USER")
	}
	if dbPass == "" {
		missing = append(missing, "DB_PASS")
	}
	if dbName == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return nil, errors.New("configuration error: missing env vars: " + fmt.Sprint(missing))
	}

	return &Config{
		PORT:   fmt.Sprintf(":%s", port),
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
		DBPass: dbPass,
		DBName: dbName,
	}, nil
}
