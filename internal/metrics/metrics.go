package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

// Метрики Prometheus
var (
    TotalRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "weather_bot_requests_total",
            Help: "Total number of requests by command",
        },
        []string{"command"},
    )
    
    ErrorRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "weather_bot_errors_total",
            Help: "Total number of errors by type",
        },
        []string{"error_type"},
    )
    
    RequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "weather_bot_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
        },
        []string{"command"},
    )
    
    ActiveUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "weather_bot_active_users",
            Help: "Currently active users in the last 5 minutes",
        },
    )
)

// Функции для удобного использования
func IncrementRequests(command string) {
    TotalRequests.WithLabelValues(command).Inc()
}

func IncrementErrors(errorType string) {
    ErrorRequests.WithLabelValues(errorType).Inc()
}

func ObserveDuration(command string, duration float64) {
    RequestDuration.WithLabelValues(command).Observe(duration)
}

func SetActiveUsers(count float64) {
    ActiveUsers.Set(count)
}

// Запуск HTTP сервера для метрик
func StartMetricsServer() {
    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Ready"))
    })
}

func init() {
    // Регистрируем метрики
    prometheus.MustRegister(TotalRequests)
    prometheus.MustRegister(ErrorRequests)
    prometheus.MustRegister(RequestDuration)
    prometheus.MustRegister(ActiveUsers)
}
