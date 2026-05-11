---
endpoint: /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candles.json
block: candles
paginated: true
---

# `get_board_candles`

**Назначение:** свечи HLOCV по бумаге **в конкретном режиме торгов** — без дублей по разным режимам (в отличие от [`get_market_candles`](get_market_candles.md)).

**Reference ISS:** <https://iss.moex.com/iss/reference/46>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candles.json
```

**Форма ответа:** блок `candles` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип        | Обязательно | Дефолт | Описание                                                            |
| ----------------- | ---------- | ----------- | ------ | ------------------------------------------------------------------- |
| `{TICKER}` (path) | str        | да          | —      | Тикер бумаги                                                        |
| `interval`        | int        | нет         | `24`   | Размер свечи (см. [`get_market_candles.md`](get_market_candles.md)) |
| `from`            | YYYY-MM-DD | нет         | —      | Начало диапазона                                                    |
| `till`            | YYYY-MM-DD | нет         | —      | Конец диапазона                                                     |
| `{board}` (path)  | str        | нет         | `TQBR` | Режим торгов                                                        |

## Пример запроса

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candles.json?iss.json=extended&iss.meta=off
```

## Поля JSON-ответа

Массив объектов с теми же полями, что у [`get_market_candles`](get_market_candles.md):

| Поле     | Тип      | Обязательно | Смысл            |
| -------- | -------- | ----------- | ---------------- |
| `open`   | float    | да          | Цена открытия    |
| `close`  | float    | да          | Цена закрытия    |
| `high`   | float    | да          | Максимум         |
| `low`    | float    | да          | Минимум          |
| `value`  | float    | да          | Оборот, ₽        |
| `volume` | int      | нет         | Объём, шт        |
| `begin`  | datetime | да          | Начало интервала |
| `end`    | datetime | нет         | Конец интервала  |

## Пример ответа

```json
[
  {
    "begin": "2026-04-30 09:00:00",
    "open": 4025,
    "close": 4032,
    "high": 4045,
    "low": 4019,
    "value": 412567890,
    "volume": 102345
  },
  {
    "begin": "2026-04-30 10:00:00",
    "open": 4032,
    "close": 4051.5,
    "high": 4055,
    "low": 4030,
    "value": 587234112,
    "volume": 145678
  }
]
```
