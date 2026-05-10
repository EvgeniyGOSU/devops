# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /weather-bot ./cmd/bot

# Этап запуска
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /weather-bot /app/weather-bot

EXPOSE 8080

CMD ["/app/weather-bot"]