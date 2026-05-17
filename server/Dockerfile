# syntax=docker/dockerfile:1.7

# Сборка Connect-сервера в two-stage образ.
# Builder тянет зависимости, компилирует статический бинарник;
# runtime — distroless static-debian12 с непривилегированным пользователем.

# Версия Go — основной рантайм проекта; зеркалит mise.toml (go = 1.26.x).
FROM golang:1.26-alpine AS builder
WORKDIR /src

# Кеширование зависимостей: слой `go mod download` пересобирается только
# при изменении манифеста модуля, а не на каждый правленный .go-файл.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO выключен — на выходе полностью статический бинарник, который
# работает в distroless static без glibc/musl.
# -trimpath и -ldflags="-s -w" убирают абсолютные пути сборки и
# debug-таблицы, чтобы образ был меньше и воспроизводим.
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS=-trimpath \
    go build -ldflags="-s -w" -o /out/server ./cmd/server

# distroless static — минимальный образ без shell и пакетного менеджера;
# содержит только корни доверия CA и /etc/passwd с пользователем nonroot.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/server /app/server
USER nonroot:nonroot
# Порт совпадает с SERVER_PORT по умолчанию (см. .env.example).
EXPOSE 8080
ENTRYPOINT ["/app/server"]
