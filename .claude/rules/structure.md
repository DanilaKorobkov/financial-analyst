# Project structure

Карта репозитория для быстрой навигации. Слои и правила изоляции — `golang.md`.

Репозиторий организован как модульный монолит: серверная часть на Go живёт в
`server/`, фронт — в `web/`. Toolchain **изолирован между модулями**: у каждого
модуля свой `mise.toml`, свой `Taskfile.yml`, своя копия
`taskfile-snippets.yml` и свой словарь cspell. Корневой слой репозитория
держит только линтеры «ничьей земли» — `.github/workflows/`, корневой
`README.md`, `.claude/**`, корневые конфиги.

## Дерево

```text
.
├── server/                                # Go-модуль приложения (модуль server)
│   ├── api/
│   │   └── proto/company/v1/              # gRPC-контракт (source of truth). См. grpc.md
│   │       ├── company.proto              # service CompanyService + агрегат Company
│   │       ├── stock.proto
│   │       ├── stock_dividend.proto
│   │       ├── stock_idea.proto
│   │       ├── stock_info.proto
│   │       ├── stock_insider_transaction.proto
│   │       ├── stock_operation.proto
│   │       ├── stock_owner.proto
│   │       ├── stock_period.proto
│   │       ├── stock_ratio.proto
│   │       ├── stock_report.proto
│   │       ├── stock_share.proto
│   │       ├── stock_summary.proto
│   │       ├── security_description.proto
│   │       └── enums.proto
│   │
│   ├── gen/                               # Сгенерированный из api/proto код (buf). Руками не править
│   │   └── company/v1/
│   │       ├── *.pb.go                    # messages
│   │       └── companyv1connect/          # Connect-клиент/сервер
│   │
│   ├── app/                               # Composition root (знает обо всех слоях)
│   │   ├── app.go                         # Wiring зависимостей
│   │   └── config.go                      # Конфиг (caarlos0/env). См. golang.md → "Конфигурация"
│   │
│   ├── cmd/server/
│   │   └── main.go                        # Точка входа бинарника, только bootstrap
│   │
│   ├── internal/
│   │   ├── domain/                        # Ядро, не знает о транспорте и провайдерах
│   │   ├── infra/                         # Адаптеры внешних систем (реализации портов)
│   │   └── presentation/                  # Транспортный слой
│   │
│   ├── mocks/                             # mockery, зеркалит internal/. Руками не править
│   │
│   ├── .cspell/project.txt                # Локальный словарь cspell scope server
│   ├── cspell.config.yaml                 # Конфиг cspell scope server
│   ├── .yamllint.yml                      # Конфиг yamllint scope server
│   ├── mise.toml                          # Toolchain модуля (go, golangci-lint, buf, mockery, линтеры)
│   ├── taskfile-snippets.yml              # Переиспользуемые snippets (копия web-версии)
│   ├── .env.example                       # Шаблон env (FINANCEMARKER_TOKEN, MOEX_BASE_URL, …)
│   ├── go.mod                             # module github.com/DanilaKorobkov/financial-analyst
│   ├── go.sum
│   ├── Taskfile.yml                       # task server:<name> (lint/test/run/docker-build/…)
│   ├── buf.yaml, buf.gen.yaml             # Конфиг buf (proto lint + кодген)
│   ├── .go-arch-lint.yml                  # Проверка изоляции слоёв
│   ├── .golangci.yml                      # Линтеры Go
│   ├── .mockery.yaml                      # Конфиг mockery
│   ├── .testcoverage.yml                  # Пороги покрытия
│   └── Dockerfile                         # multi-stage сборка сервера
│
├── web/                                   # Vite + React фронт
│   ├── src/, index.html, …
│   ├── .cspell/project.txt                # Локальный словарь cspell scope web
│   ├── cspell.config.yaml                 # Конфиг cspell scope web
│   ├── .yamllint.yml                      # Конфиг yamllint scope web
│   ├── mise.toml                          # Toolchain модуля (node, buf, protoc-gen-es, линтеры)
│   ├── taskfile-snippets.yml              # Переиспользуемые snippets (копия server-версии)
│   ├── Taskfile.yml                       # task web:<name> (lint/build/run/docker-build/…)
│   ├── buf.gen.yaml                       # Кодген TS-клиента из server/api/proto
│   └── Dockerfile                         # multi-stage сборка фронта под Caddy
│
├── .claude/
│   ├── CLAUDE.md                          # Концепция системы
│   ├── rules/                             # Конвенции проекта
│   │   ├── global.md                      # Язык, общее
│   │   ├── golang.md                      # Go: слои, имена, mockery, resty, conc, lo, jsoniter, тесты
│   │   ├── grpc.md                        # Конвенции .proto
│   │   ├── git.md                         # Branches, commits, PR
│   │   ├── tooling.md                     # mise — единственная точка правды для версий
│   │   ├── cspell.md                      # Словари и борьба с англицизмами
│   │   ├── web/code.md                    # Конвенции frontend-модуля
│   │   └── structure.md                   # ← этот файл
│   └── skills/                            # Локальные skills
│
├── .github/                               # GitHub Actions workflows
│                                          # repo-checks / server-checks / web-checks
├── .cspell/project.txt                    # Словарь cspell «ничьей земли»
├── cspell.config.yaml                     # Конфиг cspell «ничьей земли» (ignorePaths server/** web/**)
├── .yamllint.yml                          # Конфиг yamllint «ничьей земли» (server/, web/ в `ignore:`)
├── .markdownlint-cli2.jsonc               # Конфиг markdownlint (общий)
├── .prettierrc.yaml                       # Конфиг prettier (общий)
├── lychee.toml                            # Конфиг lychee (общий)
├── mise.toml                              # Toolchain «ничьей земли» (yamllint, markdownlint, cspell, lychee, actionlint)
├── Taskfile.yml                           # repo-checks + агрегатор checks/fmt/run/docker-build
├── LICENSE
└── README.md
```

## Команды

Из корня:

- `task checks` — `server:checks` + `web:checks` + `repo-checks` параллельно.
- `task fmt` — fmt во всех модулях + «ничья земля».
- `task run` — поднять server и web параллельно.
- `task docker-build` — собрать оба Docker-образа.
- `task repo-checks` — линтеры «ничьей земли».

Точечно: `task server:<name>` (например `task server:test`,
`task server:lint-golangci-lint`), `task web:<name>` (например
`task web:lint-generated-proto`, `task web:build`).

Изнутри модуля работают «короткие» имена без префикса: `cd server && task test`,
`cd web && task build`.

`task server:run` читает `server/.env` (см. `server/.env.example`).
