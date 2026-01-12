# ---------- builder ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники проекта
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---------- runtime ----------
FROM alpine:3.19

WORKDIR /app

# Устанавливаем сертификаты (для https / postgres)
RUN apk add --no-cache ca-certificates

# Копируем бинарник из builder
COPY --from=builder /app/server .

# Копируем папку миграций
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
