# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Этап запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/server .
# Копируем конфигурацию
COPY config.yaml .

# Открываем порт для gRPC
EXPOSE 50051

# Запускаем сервер
CMD ["./server"] 