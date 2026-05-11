# Golang

Конвенции проекта для Go-кода. Дополняют общие правила `~/.claude/rules/golang/*`.

## Имена файлов

Если в названии больше одного слова — разделять `_`.

- ✅ `company_info.go`, `find_security.go`
- ❌ `companyinfo.go`, `findsecurity.go`

## Имена полей-репозиториев

Поля, хранящие ссылку на репозиторий, именуются как entity **во множественном числе** — репозиторий рассматривается как коллекция (DDD).

```go
type CompanyInfo struct {
    companies entities.CompanyRepository  // ✅
}

type CompanyInfo struct {
    repo entities.CompanyRepository       // ❌
    companyRepo entities.CompanyRepository // ❌
}
```

Аналогично: `users` для `UserRepository`, `orders` для `OrderRepository`.

## Mocks

Используется **mockery** (последняя стабильная версия). Конфиг — `.mockery.yaml` в корне проекта.

Hand-rolled fakes (`entities/fakes/...`, ручные структуры с заглушками) — **запрещены**.

Шаблон `.mockery.yaml`:

```yaml
all: true
keeptree: true
case: snake
with-expecter: true
disable-version-string: true
```

## HTTP-клиент

Для внешних HTTP-вызовов используется `github.com/go-resty/resty/v2`, не голый `net/http`. Исключение — `cmd/server/main.go`, где поднимается собственный `http.Server`.

## Нормализация входных данных

Тикеры (и другие идентификаторы внешних источников) передаются **как есть** — без `strings.ToUpper`, без `strings.TrimSpace`. Источники (MOEX ISS, FinanceMarker) регистронезависимы, а нормализация на нашей стороне создаёт ложное ощущение, что без неё API не работает.

Валидация формата (regex на структуру тикера и т.п.) — тоже не нужна; ошибки формата отдаст источник.

## Doc-комментарии интерфейсов

Описание возвращаемых значений и ошибок располагается на **методе** интерфейса, а не на типе:

```go
// CompanyRepository — порт доступа к справочнику компаний.
type CompanyRepository interface {
    // FindByTicker возвращает entities.ErrCompanyNotFound, если бумаги нет.
    FindByTicker(ctx context.Context, ticker string) (Company, error)
}
```

На самом типе интерфейса — только общее назначение порта.

## Конфигурация

Конфиг приложения — `github.com/caarlos0/env/v11`. Источник — **только переменные окружения**. Никаких дефолтов в коде, никаких файлов, никаких flag'ов. Все поля помечаются `env:"NAME,required"`.

Тесты, которым нужен `Config`, выставляют env через `t.Setenv(...)`.

## Тесты

Стиль — **testify suites** (`suite.Suite`). Не голый `*testing.T`. Подробнее — в memory.
