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

## Параллелизм

Для оркестрации параллельных горутин используется `github.com/sourcegraph/conc`, не голый `sync.WaitGroup` / `golang.org/x/sync/errgroup` / самописные `chan error`.

- Fail-fast пул с отменой контекста по первой ошибке — `pool.New().WithErrors().WithContext(ctx)`:

  ```go
  import "github.com/sourcegraph/conc/pool"

  p := pool.New().WithErrors().WithContext(ctx)
  p.Go(func(ctx context.Context) error { return fetchA(ctx) })
  p.Go(func(ctx context.Context) error { return fetchB(ctx) })
  if err := p.Wait(); err != nil { ... }
  ```

- Параллельный обход коллекции — `iter.ForEach` / `iter.Map`.
- Без оркестрации ошибок — `conc.NewWaitGroup()` вместо `sync.WaitGroup` (recover-safe).

`errgroup` и голый `sync.WaitGroup` запрещены: `conc` даёт recover из паник и единый стиль на проекте.

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

## Конструкторы

Если конструктор принимает больше одного параметра (помимо `ctx`) —
параметры заворачиваются в struct `Config<Entity>`, объявленный в том же
пакете рядом с конструируемой сущностью. `ctx` остаётся отдельным
параметром, если он нужен.

```go
type ConfigCompanyCardRepository struct {
    BaseURL string
    Token   string
    Timeout time.Duration
}

func NewCompanyCardRepository(cfg ConfigCompanyCardRepository) *CompanyCardRepository { ... }   // ✅

func NewCompanyCardRepository(baseURL, token string, timeout time.Duration) *CompanyCardRepository { ... } // ❌
```

Один параметр без `ctx` — оставляем как есть, конфиг-структуру не вводим
ради единственного поля.

## Domain-ошибки (sentinel)

Sentinel-ошибка (`var ErrX = errors.New(...)`) вводится в
`internal/domain` только при выполнении **обоих** условий:

1. Ошибка явно обрабатывается каким-то слоем сверху — `service` делает
   fallback или `presentation` переводит её в конкретный HTTP/Connect-код,
   отличный от внутренней ошибки.
2. Ошибка не привязана к конкретной реализации порта — у неё есть смысл
   для всех потенциальных источников данного порта.

Если хотя бы одно из условий не выполняется — ошибку не выделяем как
sentinel. В `infra` оборачиваем причину обычным `fmt.Errorf("...: %w", err)`,
и она доезжает до presentation как непомеченный «внутренний сбой» (преобразуется в
`CodeInternal`).

Пример. `ErrCompanyNotFound` / `ErrNotFound` — да: presentation переводит в 404,
смысл общий для любого источника. `ErrUnauthorized` / `ErrQuotaExceeded`
(специфика FM-токена/квоты) — нет: для пользователя это всё равно «внутренний
сбой», другой источник (MOEX без токена) такого не имеет.

## Конфигурация

Конфиг приложения — `github.com/caarlos0/env/v11`. Источник — **только переменные окружения**. Никаких дефолтов в коде, никаких файлов, никаких flag'ов.

Тесты, которым нужен `Config`, выставляют env через `t.Setenv(...)`.

## Разрыв цепочки вызовов

Цепочка из нескольких вызовов через `.` разбивается на строки —
по одному звену на строку. Цель — глазами видеть последовательность
шагов и легко комментировать или удалять отдельные звенья в diff.

Касается прежде всего цепочек, которые регулярно растут вширь:

- mockery EXPECT (`.EXPECT().Method(...).Return(...).Once()`);
- resty (`client.R().SetContext(ctx).SetQueryParam(...).Get(url)`).

```go
// ✅ mocks
s.identities.EXPECT().
    FindByTicker(mock.Anything, "SBER").
    Return(identity, nil).
    Once()

// ❌ mocks
s.identities.EXPECT().FindByTicker(mock.Anything, "SBER").Return(identity, nil).Once()

// ✅ resty
resp, err := c.client.R().
    SetContext(ctx).
    SetPathParam("ticker", ticker).
    SetResult(&out).
    Get("/iss/securities/{ticker}.json")

// ❌ resty
resp, err := c.client.R().SetContext(ctx).SetPathParam("ticker", ticker).SetResult(&out).Get("/iss/securities/{ticker}.json")
```

Короткие цепочки из двух звеньев (`foo.Bar()`) держим в одной строке —
правило про разрыв включается с трёх звеньев и выше.

## Имена тикеров в тестах

Тестовые литералы тикеров отражают **семантику** проверки, а не
случайный набор букв. Это правило сквозное по проекту и касается
любых unit-тестов поверх domain/presentation (infra-тесты против
реальных фикстур MOEX/FM используют настоящие тикеры — там литерал
обязан совпасть с фикстурой).

| Семантика                                            | Литерал                           |
| ---------------------------------------------------- | --------------------------------- |
| Тикер не найден (проверяем `ErrNotFound`/404)        | `"missing"`                       |
| Тикер не важен (проверяем маппинг/ошибку downstream) | `"any"`                           |
| Тикер важен по смыслу (happy-path, специальный кейс) | `"SBER"` и т.п. — настоящий тикер |

```go
// ✅ говорит, что именно проверяем
s.identities.EXPECT().
    FindByTicker(mock.Anything, "missing").
    Return(company.Identity{}, company.ErrNotFound).
    Once()

// ❌ "ZZZZ" — случайный шум, читателю не помогает
s.identities.EXPECT().
    FindByTicker(mock.Anything, "ZZZZ").
    Return(company.Identity{}, company.ErrNotFound).
    Once()
```

То же относится к другим placeholder-значениям в тестах
(идентификаторы, имена пользователей, URL-ы) — литерал должен
называть свою роль.

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
