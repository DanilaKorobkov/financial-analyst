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

| Параметр          | Тип        | Дефолт    | Описание                                                                                      |
| ----------------- | ---------- | --------- | --------------------------------------------------------------------------------------------- |
| `iss.json`        | str        | `compact` | Формат payload. **Использовать `extended`** — массив `[meta, blocks]`, парсится единообразно. |
| `iss.meta`        | `on`/`off` | `on`      | Метаданные блоков. Для агента — `off`.                                                        |
| `iss.only`        | csv        | все блоки | Список блоков (через запятую), которые нужно вернуть. Снижает payload.                        |
| `<block>.columns` | csv        | все поля  | Подмножество полей блока через запятую.                                                       |
| `start`           | int        | `0`       | Смещение для пагинации.                                                                       |
| `from` / `till`   | date       | —         | Границы диапазона (YYYY-MM-DD).                                                               |
| `date`            | date       | —         | Срез на конкретную дату.                                                                      |
| `interval`        | int        | `24`      | Размер свечи (1=1м, 10=10м, 60=ч, 24=день, 7=неделя, 31=месяц, 4=квартал).                    |
| `lang`            | `ru`/`en`  | сервер    | Язык текстовых полей (где поддерживается).                                                    |

### Пагинация

Коллекционные эндпоинты пагинируются через `start=<offset>`. Признак новой страницы — блок `<name>.cursor` со строкой `{INDEX, PAGESIZE, TOTAL}`. Когда `INDEX + PAGESIZE >= TOTAL`, страниц больше нет. Часть эндпоинтов (`history`) использует `history.cursor` и требует включать его в `iss.only=history,history.cursor`.

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
- **Назначение:** поиск бумаг по подстроке (тикер, ISIN, рег.номер, имя эмитента).
- **TRIGGER when:** SECID неизвестен; есть только ISIN/рег.номер/фрагмент имени.

#### `find_security_description`

- **Reference:** [find_security_description.md](references/find_security_description.md)
- **URL:** `GET /iss/securities/{TICKER}.json` (блок `description`)
- **Назначение:** полная спецификация одной бумаги (тип, валюта, статус, листинг, квалификация).
- **TRIGGER when:** нужна шапка/карточка бумаги для отчёта.

#### `get_security_boards`

- **Reference:** [get_security_boards.md](references/get_security_boards.md)
- **URL:** `GET /iss/securities/{TICKER}.json` (блок `boards`)
- **Назначение:** все режимы торгов бумаги, включая исторические (с `history_from/till`, `is_primary`).

#### `get_security_indices`

- **Reference:** [get_security_indices.md](references/get_security_indices.md)
- **URL:** `GET /iss/securities/{TICKER}/indices.json`
- **Назначение:** индексы, в которые входит бумага (вкл. исторические).

#### `get_security_aggregates`

- **Reference:** [get_security_aggregates.md](references/get_security_aggregates.md)
- **URL:** `GET /iss/securities/{TICKER}/aggregates.json`
- **Назначение:** агрегированные итоги торгов бумаги по всем рынкам за день (shares + repo + ndm + dark).

#### `get_security_dividends`

- **Reference:** [get_security_dividends.md](references/get_security_dividends.md)
- **URL:** `GET /iss/securities/{TICKER}/dividends.json`
- **Назначение:** история дивидендов бумаги из ISS (~5–10 последних событий).
- **SKIP when:** нужна глубокая история (>10 лет) — `api-financemarker dividends`.

### Метаданные структуры биржи

#### `get_reference`

- **Reference:** [get_reference.md](references/get_reference.md)
- **URL:** `GET /iss/index.json` (мульти-блок: engines / markets / boards / boardgroups / securitytypes / securitygroups / securitycollections / durations)
- **Назначение:** глобальный справочник ISS.

#### `get_engines`

- **Reference:** [get_engines.md](references/get_engines.md)
- **URL:** `GET /iss/engines.json`
- **Назначение:** все торговые системы MOEX (stock, currency, futures, state, agro, otc, money).

#### `get_engine`

- **Reference:** [get_engine.md](references/get_engine.md)
- **URL:** `GET /iss/engines/{engine}.json` (Blocks)
- **Назначение:** описание одной системы со списком её рынков и режимов.

#### `get_markets`

- **Reference:** [get_markets.md](references/get_markets.md)
- **URL:** `GET /iss/engines/{engine}/markets.json`

#### `get_boards`

- **Reference:** [get_boards.md](references/get_boards.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards.json`

#### `get_boardgroups`

- **Reference:** [get_boardgroups.md](references/get_boardgroups.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boardgroups.json`

#### `get_securitygroups`

- **Reference:** [get_securitygroups.md](references/get_securitygroups.md)
- **URL:** `GET /iss/securitygroups.json`

#### `get_collections`

- **Reference:** [get_collections.md](references/get_collections.md)
- **URL:** `GET /iss/securitygroups/{securitygroup}/collections.json`

#### `get_collection_securities`

- **Reference:** [get_collection_securities.md](references/get_collection_securities.md)
- **URL:** `GET /iss/securitygroups/{securitygroup}/collections/{collection}/securities.json` (paginated)

### Каталог и листинг

#### `get_board_securities`

- **Reference:** [get_board_securities.md](references/get_board_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities.json`
- **Назначение:** все бумаги одного режима (статика + блок `marketdata`).

#### `get_all_stock_securities`

- **Reference:** [get_all_stock_securities.md](references/get_all_stock_securities.md)
- **URL:** `GET /iss/referencedata/engines/stock/markets/all/securities.json` (paginated)
- **Назначение:** полный каталог фондового рынка (~5400) с ISIN, INN, ISSUESIZE, LISTLEVEL.

#### `get_quoted_securities`

- **Reference:** [get_quoted_securities.md](references/get_quoted_securities.md)
- **URL:** `GET /iss/statistics/engines/stock/quotedsecurities.json`
- **Назначение:** все котируемые бумаги с `IS_QUOTED` и `MAINBOARDID`.

#### `get_securities_listing`

- **Reference:** [get_securities_listing.md](references/get_securities_listing.md)
- **URL:** `GET /iss/referencedata/engines/{engine}/markets/all/securitieslisting.json` (Blocks)

#### `get_market_listing`

- **Reference:** [get_market_listing.md](references/get_market_listing.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/listing.json` (paginated)
- **Назначение:** все когда-либо торговавшиеся инструменты рынка с history_from/till.

#### `get_board_listing`

- **Reference:** [get_board_listing.md](references/get_board_listing.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/listing.json` (paginated)

### Текущие котировки (marketdata)

#### `get_market_security`

- **Reference:** [get_market_security.md](references/get_market_security.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}.json` (Blocks)
- **Назначение:** все режимы бумаги с marketdata одним вызовом.
- **TRIGGER when:** **текущая цена** для вердикта; статус сессии; ликвидность.

#### `get_market_securities`

- **Reference:** [get_market_securities.md](references/get_market_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities.json` (Blocks)
- **Назначение:** все бумаги рынка с marketdata.

#### `get_market_secstats`

- **Reference:** [get_market_secstats.md](references/get_market_secstats.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/secstats.json`
- **Назначение:** промежуточные итоги дня по всем бумагам в разрезе сессий.

#### `get_boardgroup_securities`

- **Reference:** [get_boardgroup_securities.md](references/get_boardgroup_securities.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boardgroups/{boardgroup}/securities.json` (Blocks)

### Каталоги дат и сессии

#### `get_board_dates`

- **Reference:** [get_board_dates.md](references/get_board_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/dates.json`

#### `get_market_dates`

- **Reference:** [get_market_dates.md](references/get_market_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/dates.json`

#### `get_market_sessions`

- **Reference:** [get_market_sessions.md](references/get_market_sessions.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/sessions.json`

#### `get_market_history_dates`

- **Reference:** [get_market_history_dates.md](references/get_market_history_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}/dates.json`

#### `get_board_history_dates`

- **Reference:** [get_board_history_dates.md](references/get_board_history_dates.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/dates.json`

### Свечи (HLOCV)

#### `get_market_candles`

- **Reference:** [get_market_candles.md](references/get_market_candles.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/candles.json` (paginated)
- **Назначение:** HLOCV-свечи бумаги по всем режимам рынка.

#### `get_board_candles`

- **Reference:** [get_board_candles.md](references/get_board_candles.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candles.json` (paginated)

#### `get_market_candle_borders`

- **Reference:** [get_market_candle_borders.md](references/get_market_candle_borders.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/candleborders.json`

#### `get_board_candle_borders`

- **Reference:** [get_board_candle_borders.md](references/get_board_candle_borders.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candleborders.json`

### История торгов

#### `get_market_history`

- **Reference:** [get_market_history.md](references/get_market_history.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}.json` (paginated, `history.cursor`)
- **Назначение:** дневная история торгов по всем режимам.

#### `get_board_history`

- **Reference:** [get_board_history.md](references/get_board_history.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}.json` (paginated)
- **TRIGGER when:** расчёт доходности по основному режиму; adjusted-ряд (вместе с `get_splits_by_security` + дивидендами).

#### `get_market_history_all`

- **Reference:** [get_market_history_all.md](references/get_market_history_all.md)
- **URL:** `GET /iss/history/engines/{engine}/markets/{market}/securities.json` (paginated, `date` обязателен)
- **Назначение:** срез всех бумаг рынка за один день.

### Сделки

#### `get_market_trades`

- **Reference:** [get_market_trades.md](references/get_market_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/trades.json`

#### `get_security_trades`

- **Reference:** [get_security_trades.md](references/get_security_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/securities/{TICKER}/trades.json`

#### `get_board_security_trades`

- **Reference:** [get_board_security_trades.md](references/get_board_security_trades.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/trades.json`

> **Стаканы (orderbook) — платные** в ISS и в этот скилл не входят.

### Обороты

#### `get_all_turnovers`

- **Reference:** [get_all_turnovers.md](references/get_all_turnovers.md)
- **URL:** `GET /iss/turnovers.json` (Blocks: `turnovers` сегодня + `turnoversprevdate` вчера)

#### `get_engine_turnovers`

- **Reference:** [get_engine_turnovers.md](references/get_engine_turnovers.md)
- **URL:** `GET /iss/engines/{engine}/turnovers.json` (Blocks)

#### `get_market_turnovers`

- **Reference:** [get_market_turnovers.md](references/get_market_turnovers.md)
- **URL:** `GET /iss/engines/{engine}/markets/{market}/turnovers.json`

### Индексы

#### `get_index_tickers`

- **Reference:** [get_index_tickers.md](references/get_index_tickers.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers.json`
- **Назначение:** состав индекса (тикеры, без весов) на текущую/указанную дату.

#### `get_all_indices`

- **Reference:** [get_all_indices.md](references/get_all_indices.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics.json`

#### `get_index_analytics`

- **Reference:** [get_index_analytics.md](references/get_index_analytics.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}.json` (paginated)
- **Назначение:** состав индекса **с весами и факторами** (waprice, cap_index, ff_factor).

#### `get_index_ticker_info`

- **Reference:** [get_index_ticker_info.md](references/get_index_ticker_info.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers/{TICKER}.json` (paginated)

### Корпоративные действия

#### `get_splits_by_security`

- **Reference:** [get_splits_by_security.md](references/get_splits_by_security.md)
- **URL:** `GET /iss/statistics/engines/stock/splits/{TICKER}.json`
- **TRIGGER when:** adjusted-ряд цен (иначе SMA200/моментум сломаны).

#### `get_changeover`

- **Reference:** [get_changeover.md](references/get_changeover.md)
- **URL:** `GET /iss/history/engines/stock/markets/shares/securities/changeover.json`
- **Назначение:** история смены торговых кодов (YNDX → YDEX, TCSG → T).

### Валютный рынок

#### `get_currency_securities`

- **Reference:** [get_currency_securities.md](references/get_currency_securities.md)
- **URL:** `GET /iss/engines/currency/markets/selt/securities.json` (Blocks)

#### `get_currency_history`

- **Reference:** [get_currency_history.md](references/get_currency_history.md)
- **URL:** `GET /iss/history/engines/currency/markets/selt/securities/{TICKER}.json` (paginated)

#### `get_currency_rates`

- **Reference:** [get_currency_rates.md](references/get_currency_rates.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/selt/rates.json` (Blocks: `cbrf` + `wap_rates`)
- **TRIGGER when:** пересчёт валютной выручки экспортёров (LKOH, GMKN, NLMK, PHOR).

#### `get_fixing`

- **Reference:** [get_fixing.md](references/get_fixing.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/fixing.json` (ISS игнорирует `from`/`till`, возвращает текущий день).

#### `get_fixing_by_security`

- **Reference:** [get_fixing_by_security.md](references/get_fixing_by_security.md)
- **URL:** `GET /iss/statistics/engines/currency/markets/fixing/{TICKER}.json` (paginated)

### Облигации

#### `get_bonds_securities`

- **Reference:** [get_bonds_securities.md](references/get_bonds_securities.md)
- **URL:** `GET /iss/engines/stock/markets/bonds/securities.json` (Blocks)

#### `get_bond_yields`

- **Reference:** [get_bond_yields.md](references/get_bond_yields.md)
- **URL:** `GET /iss/history/engines/stock/markets/bonds/boards/{board}/yields/{TICKER}.json` (paginated)
- **TRIGGER when:** доходность ОФЗ для расчёта риск-фри.

#### `get_market_yields`

- **Reference:** [get_market_yields.md](references/get_market_yields.md)
- **URL:** `GET /iss/history/engines/stock/markets/bonds/yields/{TICKER}.json` (paginated)

#### `get_bonds_aggregates`

- **Reference:** [get_bonds_aggregates.md](references/get_bonds_aggregates.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/bonds/aggregates.json`

#### `get_zcyc`

- **Reference:** [get_zcyc.md](references/get_zcyc.md)
- **URL:** `GET /iss/engines/stock/zcyc.json` (Blocks: `yearyields` + `params` + `securities`)
- **TRIGGER when:** **источник риск-фри ставки для DCF** на горизонте 1–3 года.

#### `get_zcyc_history`

- **Reference:** [get_zcyc_history.md](references/get_zcyc_history.md)
- **URL:** `GET /iss/history/engines/stock/zcyc.json` (intraday-параметры NSS-кривой, ISS игнорирует `from`/`till`).

> **Купоны, оферты, амортизации** — нет в публичном ISS, бери из FinanceMarker или e-disclosure.

### Денежные ставки

#### `get_cboper_rates`

- **Reference:** [get_cboper_rates.md](references/get_cboper_rates.md)
- **URL:** `GET /iss/statistics/engines/state/markets/repo/cboper.json` (Blocks)
- **Назначение:** WAP ставки операций ЦБ по тенорам.

#### `get_rusfar`

- **Reference:** [get_rusfar.md](references/get_rusfar.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/index/rusfar.json`
- **Назначение:** RUSFAR — индикатор ставок денежного рынка по срокам (overnight, 1W, 1M, 3M, CNY).

> **Ключевая ставка ЦБ** в ISS не публикуется — она только на cbr.ru.

### Капитализация и итоги

#### `get_capitalization`

- **Reference:** [get_capitalization.md](references/get_capitalization.md)
- **URL:** `GET /iss/statistics/engines/stock/capitalization.json` (Blocks)

#### `get_totals`

- **Reference:** [get_totals.md](references/get_totals.md)
- **URL:** `GET /iss/history/engines/stock/totals/securities.json` (paginated)
- **Назначение:** итоги по всем выпускам с DAILY/MONTHLY-CAPITALIZATION.

### Статистика

#### `get_stock_correlations`

- **Reference:** [get_stock_correlations.md](references/get_stock_correlations.md)
- **URL:** `GET /iss/statistics/engines/stock/markets/shares/correlations.json`
- **Назначение:** beta/корреляции акций.

#### `get_deviationcoeffs`

- **Reference:** [get_deviationcoeffs.md](references/get_deviationcoeffs.md)
- **URL:** `GET /iss/statistics/engines/stock/deviationcoeffs.json`

#### `get_currentprices`

- **Reference:** [get_currentprices.md](references/get_currentprices.md)
- **URL:** `GET /iss/statistics/engines/stock/currentprices.json`

### Срочный рынок (FORTS)

#### `get_futures_securities`

- **Reference:** [get_futures_securities.md](references/get_futures_securities.md)
- **URL:** `GET /iss/engines/futures/markets/forts/securities.json` (Blocks)

#### `get_futures_series`

- **Reference:** [get_futures_series.md](references/get_futures_series.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/forts/series.json`

#### `get_options_securities`

- **Reference:** [get_options_securities.md](references/get_options_securities.md)
- **URL:** `GET /iss/engines/futures/markets/options/securities.json` (Blocks)

#### `get_futures_history`

- **Reference:** [get_futures_history.md](references/get_futures_history.md)
- **URL:** `GET /iss/history/engines/futures/markets/forts/securities/{TICKER}.json` (paginated)

#### `get_options_assets`

- **Reference:** [get_options_assets.md](references/get_options_assets.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/options/assets.json` (Blocks)

#### `get_indicative_rates`

- **Reference:** [get_indicative_rates.md](references/get_indicative_rates.md)
- **URL:** `GET /iss/statistics/engines/futures/markets/indicativerates/securities.json` (Blocks)

### Рейтинги (CCI бесплатные)

#### `get_rating_books`

- **Reference:** [get_rating_books.md](references/get_rating_books.md)
- **URL:** `GET /iss/cci/reference/rating-books.json`

#### `get_rating_levels`

- **Reference:** [get_rating_levels.md](references/get_rating_levels.md)
- **URL:** `GET /iss/cci/reference/rating-levels.json` (paginated)

> **Рейтинги конкретных эмитентов / выпусков** (`/iss/cci/rating/companies`, `/cci/rating/issues`) — платные, в этот скилл не входят.

### Своп-кривые

#### `get_sdfi_curves`

- **Reference:** [get_sdfi_curves.md](references/get_sdfi_curves.md)
- **URL:** `GET /iss/sdfi/curves.json`
- **Назначение:** справочник своп-кривых (~70 кривых RUB/USD/EUR/CNY).

> **Детали конкретной кривой** (`/iss/sdfi/curves/{id}`) — платные.

### Новости и события

#### `get_sitenews`

- **Reference:** [get_sitenews.md](references/get_sitenews.md)
- **URL:** `GET /iss/sitenews.json` (~50 свежих).

#### `get_news_item`

- **Reference:** [get_news_item.md](references/get_news_item.md)
- **URL:** `GET /iss/sitenews/{news_id}.json` (Blocks)

#### `get_events`

- **Reference:** [get_events.md](references/get_events.md)
- **URL:** `GET /iss/events.json` (Blocks)

#### `get_event`

- **Reference:** [get_event.md](references/get_event.md)
- **URL:** `GET /iss/events/{event_id}.json` (Blocks)

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
