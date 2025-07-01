# ────────────────────────────────────────────────────
# Stage 1: build
# ────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Копируем модули и подтягиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код и собираем бинарник
COPY . .
RUN go build -o bin/ticket-system cmd/api/main.go

# ────────────────────────────────────────────────────
# Stage 2: runtime
# ────────────────────────────────────────────────────
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Копируем сам бинарь
COPY --from=builder /app/bin/ticket-system /usr/local/bin/ticket-system

# Копируем папку миграций в /migrations
COPY --from=builder /app/migrations /migrations

# По желанию: копируем .env (если используете godotenv)
COPY .env .env

WORKDIR /

ENTRYPOINT ["/usr/local/bin/ticket-system"]
