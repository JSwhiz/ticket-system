# ────────────────────────────────────────────────────
# Stage 1: build
# ────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/ticket-system cmd/api/main.go

# ────────────────────────────────────────────────────
# Stage 2: runtime
# ────────────────────────────────────────────────────
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/ticket-system /usr/local/bin/ticket-system

COPY --from=builder /app/migrations /migrations

COPY .env .env

WORKDIR /

ENTRYPOINT ["/usr/local/bin/ticket-system"]
