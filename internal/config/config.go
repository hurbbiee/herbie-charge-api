package config

import (
	"errors"
	"os"
)

type Config struct {
	Port                   string
	LineLoginChannelID     string
	LineChannelAccessToken string
}

func LoadConfig() Config {
	return Config{
		Port:                   getEnv("PORT", "8080"),
		LineLoginChannelID:     os.Getenv("LINE_LOGIN_CHANNEL_ID"),
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	}
}

func (c Config) Validate() error {
	if c.LineLoginChannelID == "" {
		return errors.New(
			"LINE_LOGIN_CHANNEL_ID is required",
		)
	}

	if c.LineChannelAccessToken == "" {
		return errors.New(
			"LINE_CHANNEL_ACCESS_TOKEN is required",
		)
	}

	return nil
}

func getEnv(
	key string,
	defaultValue string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}