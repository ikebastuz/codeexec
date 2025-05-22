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
	PORT string
}

func GetConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("configuration error: PORT environment variable is missing")
	}

	return &Config{
		PORT: fmt.Sprintf(":%s", port),
	}, nil
}
