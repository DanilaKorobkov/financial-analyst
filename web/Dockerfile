# syntax=docker/dockerfile:1.7

# Сборка SPA-фронта в two-stage образ.
# Builder ставит npm-зависимости и собирает production-сборку;
# runtime — Caddy, раздаёт статику и делает reverse proxy Connect-вызовов на server.

# Версия Node — LTS-линейка, зеркалит mise.toml (node = 22.x).
FROM node:22-alpine AS builder
WORKDIR /src

# Кеширование зависимостей: слой `npm ci` пересобирается только при
# изменении package-lock.json/package.json, а не на каждое правление src/.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Caddy — простой статик-сервер с встроенным reverse proxy.
# Версия 2.x — последняя стабильная линейка.
FROM caddy:2-alpine
COPY --from=builder /src/dist /srv
COPY Caddyfile /etc/caddy/Caddyfile
EXPOSE 80
