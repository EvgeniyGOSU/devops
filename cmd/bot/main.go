package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"weather-scanner-bot/internal/bot"
	"weather-scanner-bot/internal/config"
	"weather-scanner-bot/internal/metrics"
	"weather-scanner-bot/internal/weather"
)

func main() {
	cfg := config.Load()

	// Запускаем сервер метрик в отдельной горутине
	go metrics.StartMetricsServer()

	weatherClient := weather.NewClient(cfg.WeatherAPIKey, cfg.WeatherAPIURL)

	weatherBot, err := bot.New(cfg.TelegramToken, weatherClient)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("✅ Бот запущен и слушает команды...")
	weatherBot.Start(ctx)
}
