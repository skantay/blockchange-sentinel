package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	APIKeyEnv   = "API_KEY"
	HTTPportEnv = "HTTP_PORT"
	envFile     = ".env"
)

var ErrEmptyEnv = errors.New("empty env")

type Config struct {
	APIKey   string
	HTTPport string
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func Load() (*Config, error) {
	if err := loadEnvFile(envFile); err != nil {
		return nil, err
	}

	apikey := os.Getenv(APIKeyEnv)
	if apikey == "" {
		return nil, ErrEmptyEnv
	}

	HTTPport := os.Getenv(HTTPportEnv)
	if HTTPport == "" {
		HTTPport = "8080"
	}

	return &Config{
		APIKey:   apikey,
		HTTPport: HTTPport,
	}, nil
}
