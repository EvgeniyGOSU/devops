package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"weather-scanner-bot/internal/bot"
	"weather-scanner-bot/internal/config"
	"weather-scanner-bot/internal/weather"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Метрики Prometheus
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_bot_requests_total",
			Help: "Total number of requests by command",
		},
		[]string{"command"},
	)

	errorRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_bot_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"error_type"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "weather_bot_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"command"},
	)

	activeUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "weather_bot_active_users",
			Help: "Currently active users in the last 5 minutes",
		},
	)
)

func init() {
	// Регистрируем метрики в Prometheus
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(errorRequests)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(activeUsers)
}

// Экспортируем функции для инкремента метрик (будут вызываться из bot.go)
func IncrementRequests(command string) {
	totalRequests.WithLabelValues(command).Inc()
}

func IncrementErrors(errorType string) {
	errorRequests.WithLabelValues(errorType).Inc()
}

func ObserveDuration(command string, duration float64) {
	requestDuration.WithLabelValues(command).Observe(duration)
}

func SetActiveUsers(count float64) {
	activeUsers.Set(count)
}

func startMetricsServer() {
	// Эндпоинт для метрик Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Ready check
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	log.Println("📊 Metrics server listening on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Printf("Metrics server error: %v", err)
	}
}

func main() {
	cfg := config.Load()

	// Запускаем сервер метрик в отдельной горутине
	go startMetricsServer()

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
