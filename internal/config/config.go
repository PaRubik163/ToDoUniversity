// Package config loads and validates process configuration.
package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config contains the settings required to run the bot and its HTTP API.
type Config struct {
	TelegramBotToken  string         `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	DatabaseURL       string         `envconfig:"DATABASE_URL" required:"true"`
	TimezoneName      string         `envconfig:"BOT_TIMEZONE" default:"Europe/Moscow"`
	Timezone          *time.Location `ignored:"true"`
	MiniAppURL        string         `envconfig:"MINI_APP_URL"`
	Port              string         `envconfig:"PORT" default:"8080"`
	JWTSecret         string         `envconfig:"JWT_SECRET" required:"true"`
	RedisURL          string         `envconfig:"REDIS_URL" default:"redis://localhost:6379/0"`
	TrustProxyHeaders bool           `envconfig:"TRUST_PROXY_HEADERS" default:"false"`
}

// Load reads environment variables, applies defaults, and validates required
// values before any application components are started.
func Load() (Config, error) {
	var result Config
	if err := envconfig.Process("", &result); err != nil {
		return Config{}, fmt.Errorf("load environment config: %w", err)
	}

	result.TelegramBotToken = strings.TrimSpace(result.TelegramBotToken)
	result.DatabaseURL = strings.TrimSpace(result.DatabaseURL)
	result.TimezoneName = strings.TrimSpace(result.TimezoneName)
	result.MiniAppURL = strings.TrimSpace(result.MiniAppURL)
	result.Port = strings.TrimSpace(result.Port)
	result.JWTSecret = strings.TrimSpace(result.JWTSecret)
	result.RedisURL = strings.TrimSpace(result.RedisURL)
	if result.TelegramBotToken == "" {
		return Config{}, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if result.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if len(result.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	if result.Port == "" {
		result.Port = "8080"
	}
	if result.RedisURL == "" {
		result.RedisURL = "redis://localhost:6379/0"
	}
	port, err := strconv.Atoi(result.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("PORT must be a number between 1 and 65535")
	}
	if result.TimezoneName == "" {
		result.TimezoneName = "Europe/Moscow"
	}
	tz, err := time.LoadLocation(result.TimezoneName)
	if err != nil {
		return Config{}, fmt.Errorf("invalid BOT_TIMEZONE %q: %w", result.TimezoneName, err)
	}
	result.Timezone = tz
	return result, nil
}
