package config

import (
    "os"
    "strconv"
)

type Config struct {
    TelegramToken string
    WeatherAPIKey string
    WeatherAPIURL string
    Port          int
}

func Load() *Config {
    return &Config{
        TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
        WeatherAPIKey: getEnv("WEATHER_API_KEY", ""),
        WeatherAPIURL: getEnv("WEATHER_API_URL", "https://api.weatherapi.com/v1"),
        Port:          getEnvAsInt("PORT", 8080),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}