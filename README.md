# financial-analyst

[![checks](https://img.shields.io/github/checks-status/DanilaKorobkov/financial-analyst/main?label=checks&cacheSeconds=60)](https://github.com/DanilaKorobkov/financial-analyst/actions?query=branch%3Amain)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)

Программный финансовый аналитик: оценивает компании и индустрии несколькими независимыми методами.

## Установка окружения

Репозиторий — модульный монолит. У каждого модуля свой toolchain:

- [`server/mise.toml`](./server/mise.toml) — Go-сервер (go, golangci-lint, buf, mockery и т.п.) и линтеры scope `server/`.
- [`web/mise.toml`](./web/mise.toml) — Vite + React фронт (node, buf, protoc-gen-es) и линтеры scope `web/`.
- [`mise.toml`](./mise.toml) (корень) — линтеры «ничьей земли» репозитория (`.github/workflows/`, корневой README, `.claude/**`).

Все инструменты ставятся через [mise](https://mise.jdx.dev/) — версии зафиксированы в `mise.toml`.

### 1. Поставить mise

macOS:

```bash
brew install mise
```

Linux/macOS (универсальный установщик):

```bash
curl https://mise.run | sh
```

### 2. Активировать shell-хук

Без хука версии не переключаются автоматически при `cd` в проект. Выбрать по своему shell — см. [официальную инструкцию](https://mise.jdx.dev/getting-started.html#activate-mise).

zsh:

```bash
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc
```

bash:

```bash
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
```

После — перезапустить терминал.

### 3. Поставить инструменты проекта

mise активируется по `cd` в каталог с `mise.toml`, поэтому:

```bash
mise install                  # корневой toolchain «ничьей земли»
mise install -C server        # toolchain Go-сервера
mise install -C web           # toolchain фронта
```

### 4. Проверить

```bash
task fmt checks
```

Из корня запускаются три параллельных слоя: `server:checks` (Go-линтеры, тесты, линтеры scope `server/`), `web:checks` (генерация TS-клиента, линтеры scope `web/`) и `repo-checks` (линтеры «ничьей земли»). Если зелено — окружение готово.

Точечно: `task server:checks`, `task web:checks`, `task repo-checks`.

## Docker-образы

Production-сборка обоих сервисов оформлена как multi-stage Docker-образы:
[`server/Dockerfile`](./server/Dockerfile) — статический бинарник Go в
distroless static-debian12, [`web/Dockerfile`](./web/Dockerfile) — статика
Vite-сборки под Caddy.

Локальная сборка обоих образов:

```bash
task docker-build
```

По отдельности — `task server:docker-build` и `task web:docker-build`.

Запуск Connect-сервера. `.env` с реальными значениями (см.
[`server/.env.example`](./server/.env.example)) обязателен — конфиг читается
только из переменных окружения:

```bash
docker run --rm --env-file server/.env -p 8080:8080 fa-server:dev
```

Запуск фронта. Smoke-проверка раздачи статики; вызовы `/company.v1.*`
без живого `server:8080` в той же Docker-сети дадут 502 — это нормально:

```bash
docker run --rm -p 8080:80 fa-web:dev
```

Для повседневной локальной разработки Docker не нужен: бэк поднимается
через `task server:run` (читает `server/.env`), фронт — через `task web:run`
(Vite dev-сервер перенаправляет `/company.v1.*` на `http://localhost:8080`);
обе команды разом — `task run`. Docker-образы используются только для
production-сборки и smoke-проверки артефактов.
