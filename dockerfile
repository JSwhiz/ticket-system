# Используем официальный образ Go
FROM golang:1.24

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь код приложения
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/api

# Устанавливаем точку входа
ENTRYPOINT ["./main"]