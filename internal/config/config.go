// Package config loads and validates process configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains the settings required to run the bot and its HTTP API.
type Config struct {
	TelegramBotToken string
	DatabaseURL      string
	Timezone         *time.Location
	MiniAppURL       string
	Port             string
}

// Load reads environment variables, applies defaults, and validates required
// values before any application components are started.
func Load() (Config, error) {
	result := Config{
		TelegramBotToken: strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		DatabaseURL:      strings.TrimSpace(os.Getenv("DATABASE_URL")),
		MiniAppURL:       strings.TrimSpace(os.Getenv("MINI_APP_URL")),
		Port:             strings.TrimSpace(os.Getenv("PORT")),
	}
	if result.TelegramBotToken == "" {
		return Config{}, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if result.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if result.Port == "" {
		result.Port = "8080"
	}
	port, err := strconv.Atoi(result.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("PORT must be a number between 1 and 65535")
	}
	timezoneName := strings.TrimSpace(os.Getenv("BOT_TIMEZONE"))
	if timezoneName == "" {
		timezoneName = "Europe/Moscow"
	}
	tz, err := time.LoadLocation(timezoneName)
	if err != nil {
		return Config{}, fmt.Errorf("invalid BOT_TIMEZONE %q: %w", timezoneName, err)
	}
	result.Timezone = tz
	return result, nil
}
