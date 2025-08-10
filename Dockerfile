# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o admin ./example/with_storage/main.go

# Runtime stage - используем scratch для минимального размера
FROM scratch

# Копируем ca-certificates из builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем собранное приложение
COPY --from=builder /build/admin /admin

# Открываем порт
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/admin"]