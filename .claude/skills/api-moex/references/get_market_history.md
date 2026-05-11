---
endpoint: /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}.json
block: history
paginated: true
---

# `get_market_history`

**Назначение:** дневная история торгов по бумаге **по всем режимам рынка** за интервал дат. На одну дату может приходиться несколько строк (по одной на режим). Если нужен один режим — [`get_board_history`](get_board_history.md).

**Reference ISS:** <https://iss.moex.com/iss/reference/63>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/securities/{TICKER}.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип        | Обязательно | Дефолт | Описание                                   |
| ----------------- | ---------- | ----------- | ------ | ------------------------------------------ |
| `{TICKER}` (path) | str        | да          | —      | Тикер бумаги                               |
| `from`            | YYYY-MM-DD | нет         | —      | Начало диапазона; пусто = с начала истории |
| `till`            | YYYY-MM-DD | нет         | —      | Конец диапазона; пусто = до конца          |
| `<block>.columns` | csv        | нет         | все    | Подмножество колонок                       |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов. Полный набор полей:

| Поле                      | Тип        | Обязательно | Смысл                                    |
| ------------------------- | ---------- | ----------- | ---------------------------------------- |
| `BOARDID`                 | str        | да          | Режим торгов                             |
| `TRADEDATE`               | date       | да          | Дата сессии                              |
| `SHORTNAME`               | str        | да          | Краткое название                         |
| `SECID`                   | str        | да          | Тикер                                    |
| `NUMTRADES`               | int        | да          | Число сделок                             |
| `VALUE`                   | float      | да          | Оборот, ₽                                |
| `OPEN`                    | float      | нет         | Цена открытия                            |
| `LOW`                     | float      | нет         | Минимум                                  |
| `HIGH`                    | float      | нет         | Максимум                                 |
| `LEGALCLOSEPRICE`         | float      | нет         | Юридическая цена закрытия                |
| `WAPRICE`                 | float      | нет         | Средневзвешенная цена                    |
| `CLOSE`                   | float      | нет         | Цена закрытия                            |
| `VOLUME`                  | int        | да          | Объём, шт                                |
| `MARKETPRICE2`            | float      | нет         | Рыночная цена 2                          |
| `MARKETPRICE3`            | float      | нет         | Рыночная цена 3                          |
| `ADMITTEDQUOTE`           | float/null | нет         | Признанная котировка                     |
| `MP2VALTRD`               | float      | нет         | Оборот сделок MARKETPRICE2               |
| `MARKETPRICE3TRADESVALUE` | float      | нет         | Оборот сделок MARKETPRICE3               |
| `ADMITTEDVALUE`           | float/null | нет         | Оборот по признанным сделкам             |
| `WAVAL`                   | float      | нет         | Объём по средневзвешенной цене           |
| `TRADINGSESSION`          | int        | нет         | Код сессии (1 утро, 2 вечер, 3 основная) |
| `CURRENCYID`              | str        | нет         | Валюта (`SUR` = ₽)                       |
| `TRENDCLSPR`              | float      | нет         | Δ цены закрытия к предыдущему дню        |
| `TRADE_SESSION_DATE`      | date       | нет         | Дата торговой сессии                     |

## Пример ответа

```json
[
  {
    "BOARDID": "TQBR",
    "TRADEDATE": "2026-04-28",
    "SECID": "YDEX",
    "CLOSE": 4056.5,
    "VOLUME": 1547709,
    "VALUE": 6357898828
  },
  {
    "BOARDID": "TQBR",
    "TRADEDATE": "2026-04-29",
    "SECID": "YDEX",
    "CLOSE": 4013.5,
    "VOLUME": 1362966,
    "VALUE": 5476523943
  },
  {
    "BOARDID": "TQBR",
    "TRADEDATE": "2026-04-30",
    "SECID": "YDEX",
    "CLOSE": 4060,
    "VOLUME": 655045,
    "VALUE": 2643169407.5
  }
]
```

## Edge cases

- На одну `TRADEDATE` придёт несколько строк (по `BOARDID`) — фильтруй сам или используй [`get_board_history`](get_board_history.md).
- Доступный диапазон дат — через [`get_board_dates`](get_board_dates.md).
