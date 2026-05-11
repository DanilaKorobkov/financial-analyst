---
endpoint: /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}.json
block: history
paginated: true
---

# `get_board_history`

**Назначение:** дневная история торгов по бумаге **в конкретном режиме торгов** — без дублей по другим режимам (в отличие от [`get_market_history`](get_market_history.md)).

**Reference ISS:** <https://iss.moex.com/iss/reference/65>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип        | Обязательно | Дефолт | Описание             |
| ----------------- | ---------- | ----------- | ------ | -------------------- |
| `{TICKER}` (path) | str        | да          | —      | Тикер бумаги         |
| `from`            | YYYY-MM-DD | нет         | —      | Начало диапазона     |
| `till`            | YYYY-MM-DD | нет         | —      | Конец диапазона      |
| `{board}` (path)  | str        | нет         | `TQBR` | Режим торгов         |
| `<block>.columns` | csv        | нет         | все    | Подмножество колонок |

## Пример запроса

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}.json?iss.json=extended&iss.meta=off
```

## Поля JSON-ответа

Поля **идентичны** [`get_market_history`](get_market_history.md):

`BOARDID`, `TRADEDATE`, `SHORTNAME`, `SECID`, `NUMTRADES`, `VALUE`, `OPEN`, `LOW`, `HIGH`, `LEGALCLOSEPRICE`, `WAPRICE`, `CLOSE`, `VOLUME`, `MARKETPRICE2`, `MARKETPRICE3`, `ADMITTEDQUOTE`, `MP2VALTRD`, `MARKETPRICE3TRADESVALUE`, `ADMITTEDVALUE`, `WAVAL`, `TRADINGSESSION`, `CURRENCYID`, `TRENDCLSPR`, `TRADE_SESSION_DATE`.

Колонка «Обязательно» — см. [`get_market_history.md`](get_market_history.md).

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
