package config

import (
	"errors"
	"fmt"
	"os"
)

const PROCESS_TIMEOUT = 30
const MAX_MEMORY = 10 * 1024 * 1024 // 10MB
const MAX_PROCESSES = 5

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
