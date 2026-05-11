---
name: api-financemarker
description: |
  Справочник по REST API [FinanceMarker.ru](https://financemarker.ru/api/swagger-ui/) — 12 эндпоинтов `/api/fm/v2/*`: глубокая карточка эмитента (мультипликаторы, отчёты, дивиденды, акционеры, инсайдеры, операционные KPI), а также рыночные ленты идей, дивидендов, инсайдерских сделок, корпоративного календаря и раскрытия. SKILL.md — роутер по эндпоинтам, детали — в `references/<endpoint>.md` (progressive disclosure).
  TRIGGER when: нужны платные/глубокие данные по российскому эмитенту MOEX (исторические P/E, P/B, EV/EBITDA, ROE, квартальные отчёты, дивидендная история и прогноз, инвест-идеи аналитиков, сделки инсайдеров, состав акционеров, операционные метрики); упомянуты financemarker / FM / fm/v2.
  SKIP when: достаточно публичных рыночных данных MOEX (свечи, история торгов, состав индекса, базовый профиль) — иди в `api-moex`; нужны live-котировки в реальном времени; нужны зарубежные тикеры (NASDAQ/NYSE/LSE/HKEX) — текущая подписка FM покрывает только MOEX (см. `references/exchanges.md`).
user-invocable: false
disable-model-invocation: false
---

# FinanceMarker.ru — REST API

Справочник по REST API сервиса [FinanceMarker](https://financemarker.ru). Скилл — **только описание эндпоинтов**: URL, параметры, поля JSON-ответа, edge cases. Вызовы выполняет Go-runtime проекта; никакого Python CLI и локального кеша в скилле нет.

> **Прогрессивная загрузка.** Этот файл — каталог эндпоинтов + краткое описание. В `references/<endpoint>.md` для каждого эндпоинта лежит сигнатура (path/query-параметры) и поля JSON-ответа. **Читай reference только когда нужны детали** — для типового использования достаточно таблицы ниже.

---

## Авторизация

Все эндпоинты требуют API-токен из профиля <https://financemarker.ru/profile>.

- Токен передаётся **query-параметром `api_token`** (например `?api_token=$FINANCE_MARKER_TOKEN`). Альтернативные варианты (заголовки `Authorization`, `Token`, `X-Token`, query `?token=` / `?auth_token=`) сервер отвергает с `{"code":400,"message":"token_not_found"}`.
- Go-runtime читает токен из переменной окружения `FINANCE_MARKER_TOKEN` и подставляет в query при сборке URL.
- В репозиторий токен не коммитим — `.env` должен быть в `.gitignore`.
- Каждый успешный запрос-ответ списывает 1 единицу `day_limit` (см. [`token_info`](references/token_info.md)).

В примерах reference-файлов токен обозначен как `$FINANCE_MARKER_TOKEN` — подставляй фактическое значение из окружения.

---

## Базовый URL и формат

- Все пути — относительно `https://financemarker.ru/api/fm/v2/`.
- Ответ — JSON. Списочные эндпоинты возвращают массив объектов; одиночные (`token_info`, `idea`, `stocks/{exchange}:{code}`) — объект.
- Семантика HTTP: `401`/`403` — нет подписки или исчерпан `day_limit`; `400`/`422` — невалидные параметры; `404` — ресурс не найден; `5xx` — серверный сбой FM.

---

## Query-параметры пагинации

Все коллекционные эндпоинты (`stocks`, `dividends`, `ideas`, `experts`, `insider_transactions`, `disclosure`, `calendar`, `operation_metrics`) принимают одинаковый набор query-параметров — в reference-файлах он не дублируется:

| Параметр | Тип | Дефолт | Описание |
|---|---|---|---|
| `limit` | int | сервер (≤100) | Размер страницы |
| `offset` | int | 0 | Сдвиг |
| `sort_by` | str | сервер | Поле сортировки |
| `sort_order` | `asc`/`desc` | сервер | Порядок |
| `updated_in_days` | int | — | Только записи, обновлённые за последние N дней |

Эндпоинт-специфичные параметры (например `mode=upcoming` у `dividends`, `include=...` у `stocks/{exchange}:{code}`) описаны в соответствующем reference.

---

## Каталог эндпоинтов

Каталог ниже — **роутер**: по описанию решай, нужен ли эндпоинт, и только тогда открывай его reference.

### Базовые / справочники

#### `token_info`
- **Reference:** [token_info.md](references/token_info.md)
- **URL:** `GET /api/fm/v2/token_info`
- **Назначение:** остаток `day_limit` и срок подписки. Используй перед серией тяжёлых запросов.
- **TRIGGER when:** диагностика 403 / планирование лимита.
- **SKIP when:** разовый запрос — не трать единицу на проверку лимита.

#### `exchanges`
- **Reference:** [exchanges.md](references/exchanges.md)
- **URL:** `GET /api/fm/v2/exchanges`
- **Назначение:** список поддерживаемых бирж. **На текущей подписке — только MOEX.**
- **TRIGGER when:** проверить, доступна ли биржа.
- **SKIP when:** работаешь только с MOEX.

#### `operation_metrics`
- **Reference:** [operation_metrics.md](references/operation_metrics.md)
- **URL:** `GET /api/fm/v2/operation_metrics`
- **Назначение:** справочник операционных KPI (id, название, единица, множитель). Используется как словарь для расшифровки `operation_metric_id` из раздела `operations` карточки эмитента.
- **TRIGGER when:** надо понять смысл `operation_metric_id` или какие KPI вообще доступны.
- **SKIP when:** значения метрик уже получены — расшифровывать не требуется.

### Эмитент

#### `stocks`
- **Reference:** [stocks.md](references/stocks.md)
- **URL:** `GET /api/fm/v2/stocks`
- **Назначение:** список компаний (`StockInfo[]`) с базовой карточкой (тикер, название, сектор, отрасль). Только пагинация и сортировка — серверной фильтрации по сектору нет, фильтруй на клиенте.
- **TRIGGER when:** нужен полный список бумаг площадки или скрининг с клиентской фильтрацией по сектору/отрасли/стране.
- **SKIP when:** известна одна конкретная бумага — `stocks/{exchange}:{code}`.

#### `stocks/{exchange}:{code}`
- **Reference:** [stock.md](references/stock.md)
- **URL:** `GET /api/fm/v2/stocks/{exchange}:{code}?include=<sections>`
- **Назначение:** **главный эндпоинт по эмитенту**. Возвращает агрегированный объект `Stock` с разделами: `info`, `summary`, `ratios`, `reports`, `dividends`, `ideas`, `insiderTransactions`, `operations`, `owners`, `shares`. Какие разделы вернуть — задаётся query-параметром `include`.
- **TRIGGER when:** нужна полная карточка эмитента / мультипликаторы / отчётность / дивиденды / инсайдеры / акционеры по конкретному тикеру.
- **SKIP when:** нужна кросс-эмитентная выборка (всё рынку) — иди в коллекционные эндпоинты ниже.

### Кросс-эмитентные ленты

#### `dividends`
- **Reference:** [dividends.md](references/dividends.md)
- **URL:** `GET /api/fm/v2/dividends`
- **Назначение:** календарь дивидендов по всему рынку. Поддерживает `mode=upcoming`.
- **TRIGGER when:** ближайшие/прошедшие выплаты по рынку, стратегия дивидендного захвата.
- **SKIP when:** дивиденды по одной бумаге — раздел `dividends` карточки эмитента.

#### `ideas`
- **Reference:** [ideas.md](references/ideas.md)
- **URL:** `GET /api/fm/v2/ideas`
- **Назначение:** лента инвест-идей всех аналитиков (заголовок, target, потенциал, статус).
- **TRIGGER when:** общая лента / фильтр по `updated_in_days` / сортировка по потенциалу.
- **SKIP when:** идеи по одной бумаге — раздел `ideas` карточки эмитента; детали одной идеи — `ideas/{id}`.

#### `ideas/{id}`
- **Reference:** [idea.md](references/idea.md)
- **URL:** `GET /api/fm/v2/ideas/{id}`
- **Назначение:** полная карточка идеи по `id` с HTML-описанием тезисов автора.
- **TRIGGER when:** есть `id` идеи и нужны её детали.
- **SKIP when:** известен только тикер.

#### `experts`
- **Reference:** [experts.md](references/experts.md)
- **URL:** `GET /api/fm/v2/experts`
- **Назначение:** рейтинг аналитиков (количество идей, % успешных, средняя доходность, средний срок).
- **TRIGGER when:** ранжировать авторов идей; найти топ-аналитиков.
- **SKIP when:** интересуют не аналитики, а их идеи.

#### `insider_transactions`
- **Reference:** [insider_transactions.md](references/insider_transactions.md)
- **URL:** `GET /api/fm/v2/insider_transactions`
- **Назначение:** лента сделок инсайдеров **полным форматом** (~24 поля).
- **TRIGGER when:** мониторинг покупок/продаж инсайдерами по рынку.
- **SKIP when:** нужны инсайдеры одной бумаги в коротком формате — раздел `insiderTransactions` карточки эмитента.

#### `disclosure`
- **Reference:** [disclosure.md](references/disclosure.md)
- **URL:** `GET /api/fm/v2/disclosure`
- **Назначение:** лента **уже раскрытых** корпоративных событий (существенные факты, заседания СД и т.п.).
- **TRIGGER when:** тайм-лайн раскрытий, мониторинг свежих фактов.
- **SKIP when:** нужны будущие события — `calendar`.

#### `calendar`
- **Reference:** [calendar.md](references/calendar.md)
- **URL:** `GET /api/fm/v2/calendar`
- **Назначение:** календарь **предстоящих** корпоративных событий.
- **TRIGGER when:** план на ближайшие даты.
- **SKIP when:** интересуют уже произошедшие события — `disclosure`; только дивиденды — `dividends`.

---

## Идентификация бумаги

Бумага в FM адресуется парой `<exchange>:<code>` (например `MOEX:SBER`, `MOEX:YDEX`, `MOEX:LKOH`, `MOEX:OZON`). Полный список бирж — `exchanges`, список тикеров — `stocks`.
