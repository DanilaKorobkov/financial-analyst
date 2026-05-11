---
endpoint: /iss/engines/{engine}/markets/{market}/securities/{TICKER}/candleborders.json
block: borders
---

# `get_market_candle_borders`

**Назначение:** узнать, **за какие диапазоны дат и для каких размеров свечей** доступны данные по бумаге **на всём рынке** (по всем режимам).

**Reference ISS:** <https://iss.moex.com/iss/reference/156>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{TICKER}/candleborders.json
```

**Форма ответа:** блок `borders`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт | Описание     |
| ----------------- | --- | ----------- | ------ | ------------ |
| `{TICKER}` (path) | str | да          | —      | Тикер бумаги |

## Пример запроса

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{TICKER}/candleborders.json?iss.json=extended&iss.meta=off
```

## Поля JSON-ответа

Массив объектов, по одному на доступный интервал:

| Поле             | Тип      | Обязательно | Смысл                                                                           |
| ---------------- | -------- | ----------- | ------------------------------------------------------------------------------- |
| `begin`          | datetime | да          | Начало доступной истории для интервала                                          |
| `end`            | datetime | да          | Конец доступной истории                                                         |
| `interval`       | int      | да          | Размер свечи (см. кодировку в [`get_market_candles.md`](get_market_candles.md)) |
| `board_group_id` | int      | да          | Группа режимов, к которой относится диапазон                                    |

## Пример ответа

```json
[
  {
    "begin": "2024-07-24 09:59:00",
    "end": "2026-05-03 18:59:59",
    "interval": 1,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-24 09:50:00",
    "end": "2026-05-03 18:59:59",
    "interval": 10,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-24 09:00:00",
    "end": "2026-05-03 18:59:59",
    "interval": 60,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-24 00:00:00",
    "end": "2026-05-03 21:04:41",
    "interval": 24,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-22 00:00:00",
    "end": "2026-05-03 00:00:00",
    "interval": 7,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-01 00:00:00",
    "end": "2026-05-31 00:00:00",
    "interval": 31,
    "board_group_id": 57
  },
  {
    "begin": "2024-07-01 00:00:00",
    "end": "2026-06-30 00:00:00",
    "interval": 4,
    "board_group_id": 57
  }
]
```

## Edge cases

- Если бумага торгуется в нескольких группах режимов — на один `interval` будет несколько строк (по одной на группу).
- Для одного режима используй [`get_board_candle_borders`](get_board_candle_borders.md).
