package bot

import (
    "context"
    "fmt"
    "log"
    "regexp"

    "weather-scanner-bot/internal/weather"
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

type WeatherBot struct {
    bot          *bot.Bot
    weatherClient *weather.Client
}

func New(token string, weatherClient *weather.Client) (*WeatherBot, error) {
    b, err := bot.New(token)
    if err != nil {
        return nil, err
    }

    wb := &WeatherBot{
        bot:          b,
        weatherClient: weatherClient,
    }

    // Регистрируем обработчик команд
    b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wb.startHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/weather", bot.MatchTypePrefix, wb.weatherHandler)
    
    return wb, nil
}

func (wb *WeatherBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   "🌤️ Привет! Я бот погоды.\nОтправь /weather Москва, чтобы узнать погоду.",
    })
}

func (wb *WeatherBot) weatherHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    // Извлекаем название города из команды
    re := regexp.MustCompile(`^/weather\s+(.+)`)
    matches := re.FindStringSubmatch(update.Message.Text)
    
    if len(matches) < 2 {
        b.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "❌ Укажите город: /weather Москва",
        })
        return
    }

    city := matches[1]
    weather, err := wb.weatherClient.GetWeather(city)
    if err != nil {
        b.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   fmt.Sprintf("❌ Ошибка: %s", err),
        })
        return
    }

    message := fmt.Sprintf(
        "🌍 %s, %s\n🌡️ Температура: %.1f°C\n💨 Ветер: %.1f км/ч\n💧 Влажность: %d%%\n☁️ %s",
        weather.Location.Name,
        weather.Location.Country,
        weather.Current.TempC,
        weather.Current.WindKph,
        weather.Current.Humidity,
        weather.Current.Condition.Text,
    )

    b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   message,
    })
}

func (wb *WeatherBot) Start(ctx context.Context) {
    log.Println("Бот запущен...")
    wb.bot.Start(ctx)
}