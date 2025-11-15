# Базовый образ для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/nanny-backend cmd/api/main.go

# Финальный образ (минимальный)
FROM alpine:latest

# Добавляем ca-сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/nanny-backend .

# Копируем статические файлы (если есть)
COPY --from=builder /app/static ./static

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["./nanny-backend"]
