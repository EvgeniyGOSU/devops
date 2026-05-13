# 🌤️ Weather Scanner Bot
[![GitHub](https://img.shields.io/badge/Telegram-Bot-blue?logo=telegram)](https://t.me/AirQuality174Bot)

**Telegram-бот для погоды с полным CI/CD пайплайном, деплоем в Kubernetes и мониторингом Prometheus + Grafana.**

Этот проект — пример production-инфраструктуры, построенной по DevOps-практикам. Всё работает на локальной инфраструктуре (VM + K3s + GitLab), но готово к масштабированию.

## 🌟 Ключевые возможности

### 🚀 Разработка и доставка
*   **Go-бот**: Реагирует на команды `/start` и `/weather <город>`, получает данные от WeatherAPI.
*   **GitLab CI/CD**: Автоматические тесты (`test`), сборка Docker-образа (`build`) и деплой в Kubernetes (`deploy`) при каждом `git push`.
*   **GitOps**: Весь код и конфигурация приложения и инфраструктуры версионируются в Git.

### ☸️ Оркестрация и инфраструктура
*   **Kubernetes (K3s)**: Приложение работает в кластере, управляется через Deployment, Service, ConfigMap и Secrets.
*   **Sidecar-прокси**: Контейнер с **Hysteria2** внутри одного Pod`а обеспечивает обход блокировок Telegram через VPS в Нидерландах (SOCKS5).
*   **GitLab Runner**: Запускает пайплайны прямо внутри Kubernetes, не требуя отдельного сервера.

### 📊 Мониторинг и наблюдаемость
*   **Prometheus**: Собирает метрики приложения (`/metrics`), которые встроены в Go-код.
*   **Grafana**: Визуализирует количество запросов, ошибок, время ответа. Дашборды доступны через NodePort.
*   **kube-prometheus-stack**: Полный стек мониторинга для Kubernetes.
