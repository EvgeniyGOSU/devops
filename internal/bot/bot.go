package bot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"weather-scanner-bot/internal/weather"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type WeatherBot struct {
	bot           *bot.Bot
	weatherClient *weather.Client
}

func New(token string, weatherClient *weather.Client) (*WeatherBot, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}

	wb := &WeatherBot{
		bot:           b,
		weatherClient: weatherClient,
	}

	// Регистрируем обработчики команд
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wb.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/weather", bot.MatchTypePrefix, wb.weatherHandler)

	return wb, nil
}

func (wb *WeatherBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Увеличиваем счётчик запросов /start
	IncrementRequests("start")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "🌤️ Привет! Я бот погоды.\nОтправь /weather Москва, чтобы узнать погоду.",
	})
}

func (wb *WeatherBot) weatherHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	startTime := time.Now()
	command := "weather"

	// Увеличиваем счётчик запросов /weather
	IncrementRequests(command)

	// Извлекаем название города из команды
	re := regexp.MustCompile(`^/weather\s+(.+)`)
	matches := re.FindStringSubmatch(update.Message.Text)

	if len(matches) < 2 {
		// Увеличиваем счётчик ошибок (неверный формат)
		IncrementErrors("invalid_format")

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Укажите город: /weather Москва",
		})
		ObserveDuration(command, time.Since(startTime).Seconds())
		return
	}

	city := matches[1]
	weatherData, err := wb.weatherClient.GetWeather(city)
	if err != nil {
		// Увеличиваем счётчик ошибок API
		IncrementErrors("api_error")

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ Ошибка: %s", err),
		})
		ObserveDuration(command, time.Since(startTime).Seconds())
		return
	}

	message := fmt.Sprintf(
		"🌍 %s, %s\n🌡️ Температура: %.1f°C\n💨 Ветер: %.1f км/ч\n💧 Влажность: %d%%\n☁️ %s",
		weatherData.Location.Name,
		weatherData.Location.Country,
		weatherData.Current.TempC,
		weatherData.Current.WindKph,
		weatherData.Current.Humidity,
		weatherData.Current.Condition.Text,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})

	// Записываем время выполнения успешного запроса
	ObserveDuration(command, time.Since(startTime).Seconds())
}

func (wb *WeatherBot) Start(ctx context.Context) {
	log.Println("✅ Бот запущен и слушает команды...")
	wb.bot.Start(ctx)
}
