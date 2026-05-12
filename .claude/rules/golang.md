# Golang

Конвенции проекта для Go-кода.

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

Не используем самописные fakes, используется **mockery**.

## HTTP-клиент

Для внешних HTTP-вызовов используется `github.com/go-resty/resty/v2`, не голый `net/http`.

## JSON

Для разбора и сериализации JSON используется `github.com/json-iterator/go` в режиме `jsoniter.ConfigCompatibleWithStandardLibrary` — drop-in для `encoding/json` (sorted map keys, html-escape, `ValidateJsonRawMessage`), но 2–3× быстрее. Тип `json.RawMessage` остаётся из `encoding/json` — это просто `[]byte`.

```go
import jsoniter "github.com/json-iterator/go"

var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary
// jsonParser.Unmarshal(...) / jsonParser.Marshal(...)
```

`ConfigFastest` НЕ используем: `MarshalFloatWith6Digits=true` режет точность float-полей (FACEVALUE и т.п.).

## Изоляция слоёв в комментариях

Слои: `domain` (ядро — `entities`, `services`), `infra` (адаптеры внешних систем — `internal/infra/<provider>/`), `presentation` (транспортный контракт — `api/proto/*.proto` и handlers `internal/presentation/`).

В doc-комментариях и публичных строках **верхнего** слоя нельзя упоминать детали **нижнего** или **смежного** слоя. Правило симметрично коду: слой, который не зависит от другого, не должен знать о нём и в текстовом виде.

| Слой           | Может упоминать                                      | НЕ должен упоминать                                            |
| -------------- | ---------------------------------------------------- | -------------------------------------------------------------- |
| `domain`       | только domain-понятия                                | провайдеров, эндпоинты, имена полей внешних API, proto/HTTP    |
| `presentation` | domain-понятия, особенности своего транспорта        | конкретных провайдеров, имена полей и эндпоинты внешних систем |
| `infra`        | domain-понятия + детали своего конкретного источника | другие источники, presentation                                 |

Хорошо:

- ✅ `// ListingLevel — котировальный уровень бумаги.` (`domain`)
- ✅ `// SecurityType — тип бумаги.` (`presentation`, `.proto`)
- ✅ `// parseListingLevel переводит LISTLEVEL блока description MOEX → entities.ListingLevel.` (`infra/moex`)

Плохо:

- ❌ `// Company — карточка эмитента MOEX. Поля из блока description /iss/securities/{TICKER}.json.` в `domain`
- ❌ `// Возвращает справочную карточку эмитента MOEX по тикеру.` в `.proto`
- ❌ `// Тикер передаётся в источник как есть — все источники (MOEX ISS, FinanceMarker) регистронезависимы.` в `domain/services`

Привязка к конкретному источнику — деталь `infra`-слоя; такие комментарии живут рядом с транслятором (`internal/infra/<provider>/translator.go`), а не в `domain` или `presentation`.

**Исключение — composition root.** `app/` и `cmd/` собирают приложение и по природе знают обо всех слоях: тип `MoexConfig`, конструктор `moex.NewCompanyRepository`, env-переменная `MOEX_BASE_URL` называются на своём «настоящем» имени. Правило изоляции к ним не применяется.

## Doc-комментарии интерфейсов

Описание возвращаемых значений и ошибок располагается на **методе** интерфейса, а не на типе:

```go
// CompanyRepository — порт доступа к справочнику компаний.
type CompanyRepository interface {
    // FindByTicker - <описание что он делаем>
    // Возвращает entities.ErrCompanyNotFound, если бумаги нет.
    // ...
    FindByTicker(ctx context.Context, ticker string) (Company, error)
}
```

На самом типе интерфейса — только общее назначение интерфейса.

## Конфигурация

Конфиг приложения — `github.com/caarlos0/env/v11`. Источник — **только переменные окружения**. Никаких дефолтов в коде, никаких файлов, никаких flag'ов.

Тесты, которым нужен `Config`, выставляют env через `t.Setenv(...)`.

## Тесты

### Имена testify suite

Тип suite именуется с **маленькой буквы** — это деталь реализации тестового файла, наружу пакета не торчит.

```go
type companyInfoSuite struct { suite.Suite }   // ✅
type CompanyInfoSuite struct { suite.Suite }   // ❌
```

### Порядок объявлений в тестовом файле

1. Тип suite (`type xxxSuite struct {...}`).
2. **Сразу под ним** — функция-runner `func TestXxxSuite(t *testing.T) { suite.Run(t, new(xxxSuite)) }`. Точка входа `go test` идёт первой, чтобы по файлу было видно, что и как запускается.
3. Хуки suite: `SetupSuite` / `SetupTest` / `TearDownTest` / `TearDownSuite`.
4. Тесты `func (s *xxxSuite) TestXxx() {...}`.
5. **Приватные helpers (`readFixture`, `mustParse` и т.п.) — в самом низу файла**, после всех тестов.

```go
type repositorySuite struct { suite.Suite }

func TestRepositorySuite(t *testing.T) {
    t.Parallel()
    suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest()       { ... }

func (s *repositorySuite) TestFindByTickerHappyPath() { ... }
func (s *repositorySuite) TestFindByTickerNotFound()  { ... }

// helpers — ниже всех тестов

func (s *repositorySuite) readFixture(name string) []byte { ... }
```
