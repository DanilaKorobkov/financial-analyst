---
endpoint: /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candleborders.json
block: borders
---

# `get_board_candle_borders`

**Назначение:** узнать, **за какие диапазоны дат и для каких размеров свечей** доступны данные по бумаге **в конкретном режиме торгов** (без дублей по другим режимам, в отличие от [`get_market_candle_borders`](get_market_candle_borders.md)).

**Reference ISS:** <https://iss.moex.com/iss/reference/48>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/candleborders.json
```

**Форма ответа:** блок `borders`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт | Описание     |
| ----------------- | --- | ----------- | ------ | ------------ |
| `{TICKER}` (path) | str | да          | —      | Тикер бумаги |
| `{board}` (path)  | str | нет         | `TQBR` | Режим торгов |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов, по одному на доступный интервал:

| Поле       | Тип      | Обязательно | Смысл                                                               |
| ---------- | -------- | ----------- | ------------------------------------------------------------------- |
| `begin`    | datetime | да          | Начало доступной истории                                            |
| `end`      | datetime | да          | Конец доступной истории                                             |
| `interval` | int      | да          | Размер свечи (см. [`get_market_candles.md`](get_market_candles.md)) |

## Пример ответа

```json
[
  {
    "begin": "2024-07-24 09:59:00",
    "end": "2026-05-03 18:59:59",
    "interval": 1
  },
  {
    "begin": "2024-07-24 09:50:00",
    "end": "2026-05-03 18:59:59",
    "interval": 10
  },
  {
    "begin": "2024-07-24 09:00:00",
    "end": "2026-05-03 18:59:59",
    "interval": 60
  },
  {
    "begin": "2024-07-24 00:00:00",
    "end": "2026-05-03 21:04:46",
    "interval": 24
  },
  {
    "begin": "2024-07-22 00:00:00",
    "end": "2026-05-03 00:00:00",
    "interval": 7
  },
  {
    "begin": "2024-07-01 00:00:00",
    "end": "2026-05-31 00:00:00",
    "interval": 31
  },
  {
    "begin": "2024-07-01 00:00:00",
    "end": "2026-06-30 00:00:00",
    "interval": 4
  }
]
```
