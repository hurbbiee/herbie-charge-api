package config

import (
	"os"
)

type Config struct {
	Port                   string
	LineLoginChannelID     string
	LineChannelAccessToken string
}

func LoadConfig() Config {
	return Config{
		Port:                   getEnv("PORT", "3000"),
		LineLoginChannelID:     os.Getenv("LINE_LOGIN_CHANNEL_ID"),
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	}
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