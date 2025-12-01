package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultGroup = "default"
)

// Config holds the application configuration loaded from environment variables
type Config struct {
	URL        string
	Passphrase string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	url := os.Getenv("COFFER_URL")
	passphrase := os.Getenv("COFFER_PASSPHRASE")

	// Remove trailing slash from URL if present
	url = strings.TrimSuffix(url, "/")

	return &Config{
		URL:        url,
		Passphrase: passphrase,
	}, nil
}

// Validate checks that required configuration is present for worker operations
func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("COFFER_URL not set\n  Set: export COFFER_URL=\"https://your-worker.workers.dev\"")
	}
	if c.Passphrase == "" {
		return fmt.Errorf("COFFER_PASSPHRASE not set\n  Set: export COFFER_PASSPHRASE=\"your-passphrase\"")
	}
	return nil
}
