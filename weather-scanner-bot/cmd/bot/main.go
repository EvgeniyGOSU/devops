package main

import (
    "context"
    "log"
    "os"
    "os/signal"

    "weather-scanner-bot/internal/bot"
    "weather-scanner-bot/internal/config"
    "weather-scanner-bot/internal/weather"
)

func main() {
    cfg := config.Load()
    
    weatherClient := weather.NewClient(cfg.WeatherAPIKey, cfg.WeatherAPIURL)
    
    weatherBot, err := bot.New(cfg.TelegramToken, weatherClient)
    if err != nil {
        log.Fatal("Ошибка создания бота:", err)
    }
    
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()
    
    weatherBot.Start(ctx)
}