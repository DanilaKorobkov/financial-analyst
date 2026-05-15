# Project structure

Карта репозитория для быстрой навигации. Слои и правила изоляции — `golang.md`.

## Дерево

```text
.
├── api/
│   └── proto/company/v1/             # gRPC-контракт (source of truth). См. grpc.md
│       ├── company.proto             # service CompanyService + агрегат Company
│       ├── stock.proto
│       ├── stock_dividend.proto
│       ├── stock_idea.proto
│       ├── stock_info.proto
│       ├── stock_insider_transaction.proto
│       ├── stock_operation.proto
│       ├── stock_owner.proto
│       ├── stock_period.proto
│       ├── stock_ratio.proto
│       ├── stock_report.proto
│       ├── stock_share.proto
│       ├── stock_summary.proto
│       ├── security_description.proto
│       └── enums.proto
│
├── gen/                              # Сгенерированный из api/proto код (buf). Руками не править
│   └── company/v1/
│       ├── *.pb.go                   # messages
│       └── companyv1connect/         # Connect-клиент/сервер
│
├── app/                              # Composition root (знает обо всех слоях)
│   ├── app.go                        # Wiring зависимостей
│   └── config.go                     # Конфиг (caarlos0/env). См. golang.md → "Конфигурация"
│
├── cmd/server/
│   └── main.go                       # Точка входа бинарника, только bootstrap
│
├── internal/
│   ├── domain/                       # Ядро, не знает о транспорте и провайдерах
│   │   ├── aggregates/
│   │   │   └── company/              # Агрегат Company: корень + секции + порты + ошибки
│   │   │       ├── company.go        # Корень агрегата
│   │   │       ├── securitydescription.go
│   │   │       ├── stockdividend.go
│   │   │       ├── stockidea.go
│   │   │       ├── stockinfo.go
│   │   │       ├── stockinsidertransaction.go
│   │   │       ├── stockoperation.go
│   │   │       ├── stockowner.go
│   │   │       ├── stockratio.go
│   │   │       ├── stockreport.go
│   │   │       ├── stockshare.go
│   │   │       ├── stocksummary.go
│   │   │       ├── stock_period.go
│   │   │       ├── stock_source.go   # Порт источника stock-секций
│   │   │       ├── repository.go     # Порт CompanyRepository
│   │   │       ├── enums.go
│   │   │       └── errors.go         # Sentinel-ошибки (ErrCompanyNotFound, …)
│   │   └── services/
│   │       └── company_service.go    # Use-case'ы поверх портов агрегатов
│   │
│   ├── infra/                        # Адаптеры внешних систем (реализации портов)
│   │   ├── company/
│   │   │   └── repository.go         # Composite-репозиторий: собирает Company из источников
│   │   ├── moex/                     # MOEX ISS
│   │   │   ├── client/client.go      # resty-клиент
│   │   │   └── securitydescription/
│   │   │       ├── source.go         # Реализация порта
│   │   │       ├── parser.go         # raw → domain
│   │   │       └── testdata/         # Фикстуры живых ответов ISS
│   │   ├── financemarker/            # FinanceMarker
│   │   │   ├── client/client.go
│   │   │   └── stock/
│   │   │       ├── source.go
│   │   │       ├── parser.go
│   │   │       ├── query.go
│   │   │       └── testdata/
│   │   └── cache/
│   │       ├── filecache/            # Файловый кэш HTTP-ответов
│   │       └── httpcache/            # HTTP middleware поверх кэша
│   │
│   └── presentation/                 # Транспортный слой
│       └── connect/
│           ├── server.go             # Connect-handlers
│           └── mapper.go             # domain ↔ proto
│
├── mocks/                            # mockery, зеркалит internal/. Руками не править
│   └── internal_/...
│
├── .claude/
│   ├── CLAUDE.md                     # Концепция системы
│   ├── rules/                        # Конвенции проекта
│   │   ├── global.md                 # Язык, общее
│   │   ├── golang.md                 # Go: слои, имена, mockery, resty, conc, lo, jsoniter, тесты
│   │   ├── grpc.md                   # Конвенции .proto
│   │   ├── git.md                    # Branches, commits, PR
│   │   ├── tooling.md                # mise — единственная точка правды для версий
│   │   ├── cspell.md                 # Словари и борьба с англицизмами
│   │   └── structure.md              # ← этот файл
│   └── skills/                       # Локальные skills
│
├── mise.toml                         # Версии инструментов. См. tooling.md
├── Taskfile.yml                      # task checks, lint-*, кодген
├── buf.yaml, buf.gen.yaml            # Конфиг buf (proto lint + кодген)
├── .go-arch-lint.yml                 # Проверка изоляции слоёв
├── .golangci.yml                     # Линтеры Go
├── .mockery.yaml                     # Конфиг mockery
├── .testcoverage.yml                 # Пороги покрытия
├── .cspell/, cspell.config.yaml      # Словари cspell
├── .env.example                      # Шаблон env (FINANCEMARKER_TOKEN, MOEX_BASE_URL, …)
├── go.mod, go.sum
└── README.md
```
