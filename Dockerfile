# Используем образ Golang для сборки приложения
FROM golang:1.22.1 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы приложения
COPY . .

# Собираем приложение
RUN go build -o cmd/main ./cmd

# Используем минимальный образ alpine для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл из предыдущего образа
COPY --from=builder /app/cmd/main .

# Запускаем приложение при старте контейнера
CMD ["./main"]
