---
name: api-moex
description: |
  Справочник по публичному REST API [MOEX ISS](https://iss.moex.com/iss/reference/) — 80 эндпоинтов поверх `iss.moex.com/iss/*`: справка по бумаге, дивиденды, сплиты, свечи (HLOCV), история торгов, marketdata, состав индексов с весами, доходности облигаций и ZCYC, валютные курсы, фьючерсы, обороты, новости. Без авторизации. SKILL.md — роутер по эндпоинтам, детали — в `references/<endpoint>.md` (progressive disclosure).
  TRIGGER when: нужны публичные рыночные данные MOEX по российскому тикеру (SBER, GAZP, LKOH и т.п.) или по рынку в целом; упомянуты MOEX / Мосбиржа / ISS / iss.moex.com.
  SKIP when: данные нужны из авторизованных источников (FinanceMarker, Bloomberg, Refinitiv); запрос про зарубежные тикеры (NASDAQ/NYSE/LSE) — MOEX ISS их не отдаёт; нужны платные данные — глубокая дивидендная история, МСФО, консенсус-прогнозы, рейтинги эмитентов/выпусков, стаканы (orderbook).
user-invocable: false
disable-model-invocation: false
---

# MOEX ISS — REST API

Справочник по публичному REST API сервиса [MOEX ISS](https://iss.moex.com/iss/reference/). **Авторизация не требуется.** Справочник — **только описание эндпоинтов**: URL, query-параметры, поля JSON-ответа, edge cases. HTTP-вызовы выполняет Go-runtime проекта.

> **Прогрессивная загрузка.** Этот файл — каталог эндпоинтов + краткое описание. В `references/<endpoint>.md` для каждого эндпоинта лежит URL, параметры и поля JSON-ответа. **Читай reference только когда нужны детали** — для типового использования достаточно таблицы ниже.

---

## Базовый URL и формат

- Все пути — относительно `https://iss.moex.com/iss/`.
- Ответ — JSON. Удобный «расширенный» формат: `iss.json=extended` — payload — массив `[meta, blocks]`, где `blocks` — объект `{ <block_name>: [ {field: value, ...}, ... ] }`.
- Метаданные блока (`metadata`) глушатся флагом `iss.meta=off`.
- Один эндпоинт может вернуть **несколько блоков** в одном payload (например, `securities` + `marketdata` + `marketdata_yields`).

### Общие query-параметры

В reference-файлах эти параметры **не дублируются** — ссылайся сюда:

| Параметр          | Тип        | Дефолт    | Описание                                                                                         |
| ----------------- | ---------- | --------- | ------------------------------------------------------------------------------------------------ |
| `iss.json`        | str        | `compact` | Формат payload. **Использовать `extended`** — массив `[meta, blocks]`, разбирается единообразно. |
| `iss.meta`        | `on`/`off` | `on`      | Метаданные блоков. Для агента — `off`.                                                           |
| `iss.only`        | csv        | все блоки | Список блоков (через запятую), которые нужно вернуть. Снижает payload.                           |
| `<block>.columns` | csv        | все поля  | Подмножество полей блока через запятую.                                                          |
| `start`           | int        | `0`       | Смещение для пагинации.                                                                          |
| `from` / `till`   | date       | —         | Границы диапазона (YYYY-MM-DD).                                                                  |
| `date`            | date       | —         | Срез на конкретную дату.                                                                         |
| `interval`        | int        | `24`      | Размер свечи (1=1м, 10=10м, 60=ч, 24=день, 7=неделя, 31=месяц, 4=квартал).                       |
| `lang`            | `ru`/`en`  | сервер    | Язык текстовых полей (где поддерживается).                                                       |

### Пагинация

Коллекционные эндпоинты разбиты на страницы — следующая запрашивается через `start=<offset>`. Признак новой страницы — блок `<name>.cursor` со строкой `{INDEX, PAGESIZE, TOTAL}`. Когда `INDEX + PAGESIZE >= TOTAL`, страниц больше нет. Часть эндпоинтов (`history`) использует `history.cursor` и требует включать его в `iss.only=history,history.cursor`.

### HTTP-коды ошибок

| Код   | Когда                                                                                |
| ----- | ------------------------------------------------------------------------------------ |
| `200` | Успех, включая пустой блок `data` (нет данных по тикеру/диапазону).                  |
| `400` | Невалидные query-параметры.                                                          |
| `404` | Эндпоинт/ресурс не существует (например, неизвестный engine/market/board или тикер). |
| `5xx` | Серверный сбой ISS — повторить запрос.                                               |

---

## Идентификация бумаги

Бумага в ISS адресуется тикером `SECID` (`SBER`, `GAZP`, `LKOH`, `YDEX`, `OZON`, `SU26212RMFS9` для ОФЗ, `USD000UTSTOM` для USD/RUB).

Если тикер неизвестен — `find_securities` по подстроке (поддерживает ISIN, рег.номер, фрагмент имени). Для редомицилированных/переименованных бумаг — `get_changeover`.

---

## Каталог эндпоинтов

Каталог ниже — **роутер**: по описанию решай, нужен ли эндпоинт, и только тогда открывай его reference.

### Поиск и карточка бумаги

#### `find_securities`

- **Reference:** [find_securities.md](references/find_securities.md)
- **URL:** `GET /iss/securities.json`
- **Назначение:** поиск бумаг по подстроке (тикер, ISIN, рег.номер, имя эмитента). Бери, когда SECID неизвестен и есть только фрагмент имени/ISIN/рег.номера. Если SECID уже известен — иди в `find_security_description`.

#### `find_security_description`

- **Reference:** [find_security_description.md](references/find_security_description.md)
- **URL:** `GET /iss/securities/{TICKER}.json` (блок `description`)
- **Назначение:** полная спецификация одной бумаги (тип, валюта, статус, листинг, квалификация). Бери для шапки/карточки бумаги в отчёте. Для торговых данных и мультипликаторов используй другие эндпоинты.

#### `get_security_boards`

- **Reference:** [get_security_boards.md](references/get_security_boards.md)
- **URL:** `GET /iss/securities/{TICKER}.json` (блок `boards`)
- **Назначение:** все режимы торгов бумаги, включая исторические (с `history_from/till`, `is_primary`). Бери, чтобы понять, в каких режимах бумага торгуется, найти основной режим или (для делистнутых) — где была активна. Если нужен только основной режим — `find_security_description` отдаёт его косвенно.

#### `get_security_indices`

- **Reference:** [get_security_indices.md](references/get_security_indices.md)
- **URL:** `GET /iss/securities/{TICKER}/indices.json`
- **Назначение:** индексы, в которые входит бумага (вкл. исторические). Бери для health-check ликвидности (входит в IMOEX = голубая фишка) или проверки членства в отраслевом субиндексе. Для обратного направления (полный состав индекса) — `get_index_tickers` или `get_index_analytics`.

#### `get_security_aggregates`

- **Reference:** [get_security_aggregates.md](references/get_security_aggregates.md)
- **URL:** `GET /iss/securities/{TICKER}/aggregates.json`
- **Назначение:** агрегированные итоги торгов бумаги по всем рынкам за день (shares + repo + ndm + dark) — для дневного полного оборота, включая repo. Если нужен только основной рынок — `get_market_security` (блок `marketdata`).

#### `get_security_dividends`

- **Reference:** [get_security_dividends.md](references/get_security_dividends.md)
- **URL:** `GET /iss/securities/{TICKER}/dividends.json`
- **Назначение:** история дивидендов бумаги из ISS (~5–10 последних событий) — бери для свежих выплат и кросс-проверки. Для глубокой истории (>10 лет) и прогнозов — `api-financemarker dividends`.

### Метаданные структуры биржи

#### `get_reference`

- **Reference:** [get_reference.md](references/get_reference.md)
- **URL:** `GET /iss/index.json` (мульти-блок: engines / markets / boards / boardgroups / securitytypes / securitygroups / securitycollections / durations)
- **Назначение:** глобальный справочник ISS — все коды для подстановки в зависимые эндпоинты (например, `boardid` для облигаций) одним запросом. Если работаешь со стандартными дефолтами (stock/shares/TQBR) — можно пропустить.

#### `get_engines`

- **Reference:** [get_engines.md](references/get_engines.md)
- **URL:** `GET /iss/engines.json`
- **Назначение:** все торговые системы MOEX (stock, currency, futures, state, agro, otc, money) — нужно, чтобы найти `{engine}` для не-акций (валюта, фьючерсы, госбумаги). Для акций используй дефолт `{engine}=stock` без запроса.

#### `get_engine`

- **Reference:** [get_engine.md](references/get_engine.md)
- **URL:** `GET /iss/engines/{engine}.json` (Blocks)
- **Назначение:** описание одной системы со списком её рынков и режимов одним запросом. Если нужен только один срез — `get_markets` или `get_boards`.

#### `get_markets`

- **Reference:** [get_markets.md](references/get_markets.md)
- **URL:** `GET /iss/engines/{engine}/markets.json`
- **Назначение:** список рынков торговой системы — нужно, чтобы найти `{market}` (shares / bonds / index / selt / forts / options / repo). Если market уже известен — иди в `get_boards`.

#### `get_boards`

- **Reference:** [get_boards.md](references/get_boards.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards.json`
- **Назначение:** справочник режимов рынка (`is_traded`, `is_primary`) — нужно, чтобы найти `{board}` под класс бумаги (TQBR / TQOB / TQTF / EQOB и т.п.). Если работаешь со стандартом TQBR (акции T+) — можно пропустить.

#### `get_boardgroups`

- **Reference:** [get_boardgroups.md](references/get_boardgroups.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boardgroups.json`
- **Назначение:** группы режимов (например, ID 57 = «Т+: Акции и ДР») — нужно перед `get_boardgroup_securities`, чтобы агрегировать запрос по группе режимов одним вызовом. Для одного режима — `get_boards`.

#### `get_securitygroups`

- **Reference:** [get_securitygroups.md](references/get_securitygroups.md)
- **URL:** `GET /iss/securitygroups.json`
- **Назначение:** глобальные группы инструментов (stock_shares, stock_bonds, stock_index, currency_selt, futures_forts) — нужно перед `get_collections`, чтобы узнать ID группы. Если ID уже известен — можно пропустить.

#### `get_collections`

- **Reference:** [get_collections.md](references/get_collections.md)
- **URL:** `GET /iss/securitygroups/{securitygroup}/collections.json`
- **Назначение:** коллекции внутри группы (для `stock_index`: all/sectoral/total_return) — нужно перед `get_collection_securities`, чтобы узнать ID коллекции. Если ID уже известен — можно пропустить.

#### `get_collection_securities`

- **Reference:** [get_collection_securities.md](references/get_collection_securities.md)
- **URL:** `GET /iss/securitygroups/{securitygroup}/collections/{collection}/securities.json` (paginated)
- **Назначение:** все бумаги коллекции с пагинацией — для запросов вида «дай все индексные бумаги», «все ОФЗ», «все привилегированные акции». Для полного каталога фондового рынка — `get_all_stock_securities`.

### Каталог и листинг

#### `get_board_securities`

- **Reference:** [get_board_securities.md](references/get_board_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities.json`
- **Назначение:** все бумаги одного режима (статика + блок `marketdata`) — для скрининга по режиму TQBR (все акции T+) или TQOB (ОФЗ). Для среза по всему рынку — `get_market_securities`.

#### `get_all_stock_securities`

- **Reference:** [get_all_stock_securities.md](references/get_all_stock_securities.md)
- **URL:** `GET /iss/referencedata/engines/stock/markets/all/securities.json` (paginated)
- **Назначение:** полный каталог фондового рынка (~5400) с ISIN, INN, ISSUESIZE, LISTLEVEL — для скрининга по уровню листинга/размеру выпуска и для получения ISIN/INN пакета бумаг. Для одного режима — `get_board_securities`.

#### `get_quoted_securities`

- **Reference:** [get_quoted_securities.md](references/get_quoted_securities.md)
- **URL:** `GET /iss/statistics/engines/stock/quotedsecurities.json`
- **Назначение:** все котируемые бумаги с `IS_QUOTED` и `MAINBOARDID` — быстрая проверка «торгуется ли бумага сегодня» + её основной режим. Для полной карточки — `find_security_description`.

#### `get_securities_listing`

- **Reference:** [get_securities_listing.md](references/get_securities_listing.md)
- **URL:** `GET /iss/referencedata/engines/{engine}/markets/all/securitieslisting.json` (Blocks)
- **Назначение:** справочник торгуемости в формате referencedata — нужен, чтобы дополнить статус листинга мета-полями referencedata. Если достаточно `IS_QUOTED` — `get_quoted_securities`.

#### `get_market_listing`

- **Reference:** [get_market_listing.md](references/get_market_listing.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/listing.json` (paginated)
- **Назначение:** все когда-либо торговавшиеся инструменты рынка с history_from/till — для скрининга по делистнутым и полного архива SECID. Для только активных бумаг — `get_board_securities`.

#### `get_board_listing`

- **Reference:** [get_board_listing.md](references/get_board_listing.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/listing.json` (paginated)
- **Назначение:** все инструменты одного режима с history_from/till — архив SECID конкретного режима (например, все когда-либо торговавшиеся ETF в TQTF). Для целого рынка — `get_market_listing`.

### Текущие котировки (marketdata)

#### `get_market_security`

- **Reference:** [get_market_security.md](references/get_market_security.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}.json` (Blocks)
- **Назначение:** все режимы бумаги с marketdata одним вызовом — основной источник **текущей цены** для вердикта, статуса сессии и ликвидности (spread, depth). Для исторической цены — `get_board_history`.

#### `get_market_securities`

- **Reference:** [get_market_securities.md](references/get_market_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities.json` (Blocks)
- **Назначение:** все бумаги рынка с marketdata — для скрининга рынка с current price (например, все акции TQBR одним вызовом). Для одной бумаги — `get_market_security`.

#### `get_market_secstats`

- **Reference:** [get_market_secstats.md](references/get_market_secstats.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/secstats.json`
- **Назначение:** промежуточные итоги дня по всем бумагам в разрезе сессий — для intraday-обзора рынка (открытие/середина/закрытие). Для финальных дневных итогов — `get_market_history_all`.

#### `get_boardgroup_securities`

- **Reference:** [get_boardgroup_securities.md](references/get_boardgroup_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boardgroups/{boardgroup}/securities.json` (Blocks)
- **Назначение:** бумаги группы режимов с marketdata — для среза по группе родственных режимов одним запросом (например, все T+). Для одного режима — `get_board_securities`.

### Каталоги дат и сессии

#### `get_board_dates`

- **Reference:** [get_board_dates.md](references/get_board_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/dates.json`
- **Назначение:** диапазон дат истории режима (from..till) — нужен перед длинным запросом истории, чтобы узнать границы данных режима. Для диапазона по конкретной бумаге — `get_board_history_dates`.

#### `get_market_dates`

- **Reference:** [get_market_dates.md](references/get_market_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/dates.json`
- **Назначение:** диапазон дат истории всего рынка — нужен перед `get_market_history_all`, чтобы узнать, за какие даты есть данные. Для границ конкретного режима — `get_board_dates`.

#### `get_market_sessions`

- **Reference:** [get_market_sessions.md](references/get_market_sessions.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/sessions.json`
- **Назначение:** список торговых сессий (morning/main/evening/total/weekend) — нужен для различения сессионных цен (утренняя/основная/вечерняя) при анализе intraday. Для дневной агрегации не нужен.

#### `get_market_history_dates`

- **Reference:** [get_market_history_dates.md](references/get_market_history_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}/dates.json`
- **Назначение:** диапазон дат истории конкретной бумаги по всем режимам рынка — нужен перед скачиванием истории, чтобы узнать первую и последнюю торгуемую дату. Для одного режима — `get_board_history_dates`.

#### `get_board_history_dates`

- **Reference:** [get_board_history_dates.md](references/get_board_history_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/dates.json`
- **Назначение:** диапазон дат истории бумаги в одном режиме — нужен, чтобы ограничить запрос истории конкретным режимом и узнать его границы. Для всех режимов — `get_market_history_dates`.

### Свечи (HLOCV)

#### `get_market_candles`

- **Reference:** [get_market_candles.md](references/get_market_candles.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/candles.json` (paginated)
- **Назначение:** HLOCV-свечи бумаги по всем режимам рынка (интервалы: 1м, 10м, 1ч, день, неделя, месяц, квартал) — для графика, теханализа и расчёта доходности. Для свечей строго одного режима (без межрежимных дублей) — `get_board_candles`.

#### `get_board_candles`

- **Reference:** [get_board_candles.md](references/get_board_candles.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candles.json` (paginated)
- **Назначение:** HLOCV-свечи в одном режиме (без межрежимных дублей) — свечи строго основного режима (TQBR) для корректной adjusted-серии. Если допустимы свечи со всех режимов — `get_market_candles`.

#### `get_market_candle_borders`

- **Reference:** [get_market_candle_borders.md](references/get_market_candle_borders.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/candleborders.json`
- **Назначение:** доступные интервалы свечей по бумаге на уровне рынка — нужен перед массовой загрузкой свечей, чтобы проверить, какие интервалы поддержаны. Для стандартных интервалов (день/час) можно пропустить. Для уровня режима — `get_board_candle_borders`.

#### `get_board_candle_borders`

- **Reference:** [get_board_candle_borders.md](references/get_board_candle_borders.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candleborders.json`
- **Назначение:** доступные интервалы свечей бумаги в одном режиме — для проверки интервалов под конкретный режим. Для уровня рынка — `get_market_candle_borders`.

### История торгов

#### `get_market_history`

- **Reference:** [get_market_history.md](references/get_market_history.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}.json` (paginated, `history.cursor`)
- **Назначение:** дневная история торгов по всем режимам — для цен закрытия из всех режимов, в т.ч. вне TQBR. Для adjusted-серии используй только основной режим — `get_board_history`.

#### `get_board_history`

- **Reference:** [get_board_history.md](references/get_board_history.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}.json` (paginated)
- **Назначение:** дневная история по одному режиму (TQBR для акций по умолчанию) — для расчёта доходности по основному режиму и adjusted-ряда (вместе с `get_splits_by_security` + дивидендами). Для истории по всем режимам — `get_market_history`.

#### `get_market_history_all`

- **Reference:** [get_market_history_all.md](references/get_market_history_all.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities.json` (paginated, `date` обязателен)
- **Назначение:** срез всех бумаг рынка за один день — для запроса «дай весь рынок на дату» (закрытия по всем тикерам сразу). Для истории одной бумаги — `get_market_history` или `get_board_history`.

### Сделки

#### `get_market_trades`

- **Reference:** [get_market_trades.md](references/get_market_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/trades.json`
- **Назначение:** последние ~50 сделок всего рынка — быстрая проверка, идут ли сделки на рынке прямо сейчас. Для сделок по конкретной бумаге — `get_security_trades`.

#### `get_security_trades`

- **Reference:** [get_security_trades.md](references/get_security_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/trades.json`
- **Назначение:** последние сделки по одной бумаге — для диагностики ликвидности конкретной бумаги (частота тиков, объём). Чтобы ограничиться одним режимом — `get_board_security_trades`.

#### `get_board_security_trades`

- **Reference:** [get_board_security_trades.md](references/get_board_security_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/trades.json`
- **Назначение:** последние сделки по бумаге в одном режиме — сделки строго основного режима, без шума repo/dark. Если допустимы все режимы — `get_security_trades`.

> **Стаканы (orderbook) — платные** в ISS и в этот справочник не входят.

### Обороты

#### `get_all_turnovers`

- **Reference:** [get_all_turnovers.md](references/get_all_turnovers.md)
- **URL:** `GET /iss/turnovers.json` (Blocks: `turnovers` сегодня + `turnoversprevdate` вчера)
- **Назначение:** сводные обороты всех рынков MOEX — макрооценка состояния биржи и сравнение оборотов «сегодня vs. вчера». Для среза по одной системе — `get_engine_turnovers`.

#### `get_engine_turnovers`

- **Reference:** [get_engine_turnovers.md](references/get_engine_turnovers.md)
- **URL:** `GET /iss/engines/{engine}/turnovers.json` (Blocks)
- **Назначение:** обороты одной торговой системы по её рынкам — для сравнения оборотов внутри stock / currency / futures. Для одного рынка — `get_market_turnovers`.

#### `get_market_turnovers`

- **Reference:** [get_market_turnovers.md](references/get_market_turnovers.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/turnovers.json`
- **Назначение:** оборот одного рынка — точечная цифра оборота рынка (например, акции, облигации) за день. Для среза по нескольким системам — `get_all_turnovers`.

### Индексы

#### `get_index_tickers`

- **Reference:** [get_index_tickers.md](references/get_index_tickers.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers.json`
- **Назначение:** состав индекса (только тикеры, без весов) на текущую/указанную дату — ответ на вопрос «кто входит в IMOEX/MOEXBC/RTSI». Если нужны **веса** — `get_index_analytics`.

#### `get_all_indices`

- **Reference:** [get_all_indices.md](references/get_all_indices.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics.json`
- **Назначение:** список всех индексов MOEX (~290 включая отраслевые субиндексы) — нужно, чтобы найти ID отраслевого/специального индекса (MOEXIT, MOEXBC, MOEXOG и т.п.). Если конкретный индекс уже известен — можно пропустить.

#### `get_index_analytics`

- **Reference:** [get_index_analytics.md](references/get_index_analytics.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}.json` (paginated)
- **Назначение:** состав индекса **с весами и факторами** (waprice, cap_index, ff_factor) — нужны фактические веса в IMOEX (для benchmark / peer-группировки). Если достаточно только тикеров — `get_index_tickers`.

#### `get_index_ticker_info`

- **Reference:** [get_index_ticker_info.md](references/get_index_ticker_info.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers/{TICKER}.json` (paginated)
- **Назначение:** аналитика по тикеру в индексе — история включений с весами для отслеживания динамики веса бумаги в индексе по датам ребалансировки. Для текущего состава — `get_index_analytics`.

### Корпоративные действия

#### `get_splits_by_security`

- **Reference:** [get_splits_by_security.md](references/get_splits_by_security.md)
- **URL:** `GET /iss/statistics/engines/stock/splits/{TICKER}.json`
- **Назначение:** все сплиты по бумаге — обязательная корректировка для adjusted-ряда цен (иначе SMA200 / momentum-индикаторы сломаны), проверка истории VTBR/GMKN. Для дивидендов — `get_security_dividends` или `api-financemarker`.

#### `get_changeover`

- **Reference:** [get_changeover.md](references/get_changeover.md)
- **URL:** `GET /iss/history/engines/stock/markets/shares/securities/changeover.json`
- **Назначение:** история смены торговых кодов (YNDX → YDEX, TCSG → T) — для редомицилированных/реорганизованных эмитентов, чтобы найти историю под старым SECID. Если уверен, что тикер не менялся — можно пропустить.

### Валютный рынок

#### `get_currency_securities`

- **Reference:** [get_currency_securities.md](references/get_currency_securities.md)
- **URL:** `GET /iss/engines/currency/markets/selt/securities.json` (Blocks)
- **Назначение:** все валютные пары SELT с marketdata — срез по всем валютным парам Мосбиржи с текущими курсами. Для истории одной пары — `get_currency_history`.

#### `get_currency_history`

- **Reference:** [get_currency_history.md](references/get_currency_history.md)
- **URL:** `GET /iss/history/engines/currency/markets/selt/securities/{TICKER}.json` (paginated)
- **Назначение:** дневная история валютной пары — ряд курса USD/RUB или CNY/RUB для пересчёта валютной выручки. Для фиксингов/курсов ЦБ — `get_fixing_by_security` или `get_currency_rates`.

#### `get_currency_rates`

- **Reference:** [get_currency_rates.md](references/get_currency_rates.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/selt/rates.json` (Blocks: `cbrf` + `wap_rates`)
- **Назначение:** курсы ЦБ + WAP из торгов — для пересчёта валютной выручки экспортёров (LKOH, GMKN, NLMK, PHOR). Для исторического ряда — `get_currency_history`.

#### `get_fixing`

- **Reference:** [get_fixing.md](references/get_fixing.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/fixing.json` (ISS игнорирует `from`/`till`, возвращает текущий день).
- **Назначение:** биржевые валютные фиксинги на сегодня — для запроса «дай официальные фиксинги Мосбиржи на сегодня». Для истории фиксинга — `get_fixing_by_security`.

#### `get_fixing_by_security`

- **Reference:** [get_fixing_by_security.md](references/get_fixing_by_security.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/fixing/{TICKER}.json` (paginated)
- **Назначение:** история фиксинга конкретной валютной пары — ряд биржевых фиксингов USD/CNY за диапазон. Если достаточно текущего дня — `get_fixing`.

### Облигации

#### `get_bonds_securities`

- **Reference:** [get_bonds_securities.md](references/get_bonds_securities.md)
- **URL:** `GET /iss/engines/stock/markets/bonds/securities.json` (Blocks)
- **Назначение:** все облигации с marketdata + marketdata_yields — для скрининга рынка облигаций с текущими доходностями. Для одной бумаги — `get_bond_yields` или `get_market_yields`.

#### `get_bond_yields`

- **Reference:** [get_bond_yields.md](references/get_bond_yields.md)
- **URL:** `GET /iss/history/engines/stock/markets/bonds/boards/{board}/yields/{TICKER}.json` (paginated)
- **Назначение:** история доходностей облигации в одном режиме — доходность ОФЗ для расчёта риск-фри. Для всех режимов — `get_market_yields`.

#### `get_market_yields`

- **Reference:** [get_market_yields.md](references/get_market_yields.md)
- **URL:** `GET /iss/history/engines/stock/markets/bonds/yields/{TICKER}.json` (paginated)
- **Назначение:** история доходностей облигации по всем режимам — для доходностей из всех режимов (например, TQOB + EQOB). Для одного режима — `get_bond_yields`.

#### `get_bonds_aggregates`

- **Reference:** [get_bonds_aggregates.md](references/get_bonds_aggregates.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/bonds/aggregates.json`
- **Назначение:** агрегаты рынка облигаций по типам (ОФЗ / корп / субфед) — обзор рынка облигаций «крупными мазками», обороты по сегментам. Для конкретных бумаг — `get_bonds_securities`.

#### `get_zcyc`

- **Reference:** [get_zcyc.md](references/get_zcyc.md)
- **URL:** `GET /iss/engines/stock/zcyc.json` (Blocks: `yearyields` + `params` + `securities`)
- **Назначение:** кривая бескупонной доходности ОФЗ (NSS) — **основной источник риск-фри ставки для DCF** на горизонте 1–3 года, срез кривой на дату. Для intraday-динамики параметров кривой — `get_zcyc_history`.

#### `get_zcyc_history`

- **Reference:** [get_zcyc_history.md](references/get_zcyc_history.md)
- **URL:** `GET /iss/history/engines/stock/zcyc.json` (intraday-параметры NSS-кривой, ISS игнорирует `from`/`till`).
- **Назначение:** intraday-параметры NSS-кривой за текущий день — для отслеживания колебаний параметров кривой в течение дня. Для среза кривой на дату — `get_zcyc`.

> **Купоны, оферты, амортизации** — нет в публичном ISS, бери из FinanceMarker или e-disclosure.

### Денежные ставки

#### `get_cboper_rates`

- **Reference:** [get_cboper_rates.md](references/get_cboper_rates.md)
- **URL:** `GET /iss/statistics/engines/state/markets/repo/cboper.json` (Blocks)
- **Назначение:** WAP ставки операций ЦБ по тенорам — proxy денежной ставки (REPO/Lombard/Deposit) для краткосрочных моделей. Для индикативной ставки денежного рынка — `get_rusfar`.

#### `get_rusfar`

- **Reference:** [get_rusfar.md](references/get_rusfar.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/rusfar.json`
- **Назначение:** RUSFAR — индикатор ставок денежного рынка по срокам (overnight, 1W, 1M, 3M, CNY); альтернатива ключевой ставке ЦБ для коротких горизонтов. Для конкретных операций ЦБ — `get_cboper_rates`.

> **Ключевая ставка ЦБ** в ISS не публикуется — она только на cbr.ru.

### Капитализация и итоги

#### `get_capitalization`

- **Reference:** [get_capitalization.md](references/get_capitalization.md)
- **URL:** `GET /iss/statistics/engines/stock/capitalization.json` (Blocks)
- **Назначение:** суммарная капитализация фондового рынка — макроцифра «рыночная капа MOEX на дату». Для капитализации конкретной бумаги — поле `ISSUECAPITALIZATION` в `get_market_security`.

#### `get_totals`

- **Reference:** [get_totals.md](references/get_totals.md)
- **URL:** `GET /iss/history/engines/stock/totals/securities.json` (paginated)
- **Назначение:** итоги по всем выпускам с DAILY/MONTHLY-CAPITALIZATION — месячная динамика капитализации по эмитентам (long history). Для точки на сегодня — `get_capitalization`.

### Статистика

#### `get_stock_correlations`

- **Reference:** [get_stock_correlations.md](references/get_stock_correlations.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/shares/correlations.json`
- **Назначение:** готовые beta/корреляции акций — для поиска beta пары акций или построения корреляционной матрицы peer-группы. Если считаешь beta сам из ряда доходностей — `get_board_candles` + `get_market_candles` индекса.

#### `get_deviationcoeffs`

- **Reference:** [get_deviationcoeffs.md](references/get_deviationcoeffs.md)
- **URL:** `GET /iss/statistics/engines/stock/deviationcoeffs.json`
- **Назначение:** коэффициенты отклонения (sigma, beta, f_plus, f_minus, spread) — ready-to-use sigma/beta из MOEX без собственного расчёта. Для корреляции пары — `get_stock_correlations`.

#### `get_currentprices`

- **Reference:** [get_currentprices.md](references/get_currentprices.md)
- **URL:** `GET /iss/statistics/engines/stock/currentprices.json`
- **Назначение:** история текущих цен в разрезе сессий (CURPRICE/LASTPRICE/LEGALCLOSE) — для разделения «цена закрытия» по сессиям (morning/main/evening). Если достаточно итогового close — `get_market_history` или `get_board_history`.

### Срочный рынок (FORTS)

#### `get_futures_securities`

- **Reference:** [get_futures_securities.md](references/get_futures_securities.md)
- **URL:** `GET /iss/engines/futures/markets/forts/securities.json` (Blocks)
- **Назначение:** активные фьючерсы FORTS с marketdata — срез активных фьючерсов с текущими ценами и базисами. Для истёкших серий — `get_futures_series`.

#### `get_futures_series`

- **Reference:** [get_futures_series.md](references/get_futures_series.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/forts/series.json`
- **Назначение:** все серии фьючерсов, включая истёкшие — для разбора истории экспирированных контрактов (Si, Br, RTS). Для только активных — `get_futures_securities`.

#### `get_options_securities`

- **Reference:** [get_options_securities.md](references/get_options_securities.md)
- **URL:** `GET /iss/engines/futures/markets/options/securities.json` (Blocks)
- **Назначение:** активные опционы FORTS — срез опционов по страйкам и сериям одним вызовом. Для серий по активу (а не по страйкам) — `get_options_assets`.

#### `get_futures_history`

- **Reference:** [get_futures_history.md](references/get_futures_history.md)
- **URL:** `GET /iss/history/engines/futures/markets/forts/securities/{TICKER}.json` (paginated)
- **Назначение:** история торгов одного фьючерса — ряд закрытий фьючерса (например, Si-3.26) за диапазон. Для marketdata по активным фьючерсам — `get_futures_securities`.

#### `get_options_assets`

- **Reference:** [get_options_assets.md](references/get_options_assets.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/options/assets.json` (Blocks)
- **Назначение:** опционные серии по всем активам с открытыми позициями — обзор открытого интереса по активам опционов. Для конкретных страйков — `get_options_securities`.

#### `get_indicative_rates`

- **Reference:** [get_indicative_rates.md](references/get_indicative_rates.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/indicativerates/securities.json` (Blocks)
- **Назначение:** индикативные курсы валют срочного рынка (CAD/RUB и др.) — для валют, которых нет на SELT. Если валюта торгуется на SELT — `get_currency_rates`.

### Рейтинги (CCI бесплатные)

#### `get_rating_books`

- **Reference:** [get_rating_books.md](references/get_rating_books.md)
- **URL:** `GET /iss/cci/reference/rating-books.json`
- **Назначение:** справочник рейтинговых шкал (АКРА, Эксперт РА, НКР, НРА, ЦБ) — для сопоставления агентства и его шкалы рейтингов. Для уровней внутри шкалы — `get_rating_levels`.

#### `get_rating_levels`

- **Reference:** [get_rating_levels.md](references/get_rating_levels.md)
- **URL:** `GET /iss/cci/reference/rating-levels.json` (paginated)
- **Назначение:** уровни рейтингов с описаниями (~800 записей) — расшифровка уровней AAA/AA/A/BBB по агентствам. Сам рейтинг конкретного эмитента в публичном ISS не отдаётся — см. предупреждение ниже.

> **Рейтинги конкретных эмитентов / выпусков** (`/iss/cci/rating/companies`, `/cci/rating/issues`) — платные, в этот справочник не входят.

### Своп-кривые

#### `get_sdfi_curves`

- **Reference:** [get_sdfi_curves.md](references/get_sdfi_curves.md)
- **URL:** `GET /iss/sdfi/curves.json`
- **Назначение:** справочник своп-кривых (~70 кривых RUB/USD/EUR/CNY) — нужно, чтобы найти ID своп-кривой по валюте/типу. Сами точки на кривой (`/iss/sdfi/curves/{id}`) платные.

> **Детали конкретной кривой** (`/iss/sdfi/curves/{id}`) — платные.

### Новости и события

#### `get_sitenews`

- **Reference:** [get_sitenews.md](references/get_sitenews.md)
- **URL:** `GET /iss/sitenews.json` (~50 свежих).
- **Назначение:** свежие новости биржи — обзор последних объявлений MOEX (правки расписания, листинги, новые инструменты). Для полного текста конкретной новости — `get_news_item`.

#### `get_news_item`

- **Reference:** [get_news_item.md](references/get_news_item.md)
- **URL:** `GET /iss/sitenews/{news_id}.json` (Blocks)
- **Назначение:** полный текст новости по ID (HTML) — когда есть `news_id` из `get_sitenews` и нужно тело новости. Если достаточно заголовков — `get_sitenews`.

#### `get_events`

- **Reference:** [get_events.md](references/get_events.md)
- **URL:** `GET /iss/events.json` (Blocks)
- **Назначение:** активные мероприятия биржи — обзор предстоящих конференций/вебинаров MOEX. Для текста конкретного события — `get_event`.

#### `get_event`

- **Reference:** [get_event.md](references/get_event.md)
- **URL:** `GET /iss/events/{event_id}.json` (Blocks)
- **Назначение:** детали мероприятия по ID — программа/контент, когда есть `event_id` из `get_events`. Если достаточно списка — `get_events`.

---

## Что НЕ покрыто (платные эндпоинты ISS)

| Эндпоинт                                           | Причина                                                |
| -------------------------------------------------- | ------------------------------------------------------ |
| `/iss/cci/corp-actions/dividends`                  | Платно (глубокая история — через `api-financemarker`). |
| `/iss/cci/corp-actions/coupons`                    | Платно.                                                |
| `/iss/cci/accounting/msfo-full/...`                | МСФО — платно (через `api-financemarker`).             |
| `/iss/cci/consensus/shares-price`                  | Консенсус-прогнозы — платно.                           |
| `/iss/cci/rating/companies` / `/cci/rating/issues` | Рейтинги эмитентов/выпусков — платно.                  |
| `/iss/cci/info/companies`                          | Карточка компании — пустой ответ без подписки.         |
| `/iss/sdfi/curves/{id}`                            | Детали SDFI-кривой — закрыто.                          |
| `/iss/.../orderbook`                               | Стаканы заявок — платно.                               |
| `/iss/analyticalproducts/netflow2/futoi`           | Free-tier с задержкой 14 дней (де-факто платно).       |
