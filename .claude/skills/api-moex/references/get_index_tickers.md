---
endpoint: /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers.json
block: tickers
---

# `get_index_tickers`

**Назначение:** получить **состав индекса** на конкретную дату или всю историю изменений состава.

**Reference ISS:** <https://iss.moex.com/iss/reference/148>
**Список индексов:** <https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics/{index}/tickers.json
```

**Форма ответа:** блок `tickers`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип        | Обязательно | Дефолт | Описание                                                                                       |
| ----------------- | ---------- | ----------- | ------ | ---------------------------------------------------------------------------------------------- |
| `INDEX`           | str        | да          | —      | Код индекса (`IMOEX`, `RTSI`, `MOEXBC`, …)                                                     |
| `date`            | YYYY-MM-DD | нет         | —      | Если задано — только активные на эту дату; если в этот день торгов не было — **пустой массив** |
| `<block>.columns` | csv        | нет         | все    | Подмножество колонок                                                                           |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов:

| Поле             | Тип  | Обязательно | Смысл                                     |
| ---------------- | ---- | ----------- | ----------------------------------------- |
| `ticker`         | str  | да          | Тикер инструмента в составе индекса       |
| `from`           | date | да          | Дата включения в индекс                   |
| `till`           | date | да          | Дата исключения (или последняя доступная) |
| `tradingsession` | int  | да          | Код торговой сессии                       |

## Пример ответа

```json
[
  {
    "ticker": "AFKS",
    "from": "2026-04-30",
    "till": "2026-04-30",
    "tradingsession": 3
  },
  {
    "ticker": "GAZP",
    "from": "2026-04-30",
    "till": "2026-04-30",
    "tradingsession": 3
  },
  {
    "ticker": "LKOH",
    "from": "2026-04-30",
    "till": "2026-04-30",
    "tradingsession": 3
  },
  {
    "ticker": "SBER",
    "from": "2026-04-30",
    "till": "2026-04-30",
    "tradingsession": 3
  },
  {
    "ticker": "YDEX",
    "from": "2026-04-30",
    "till": "2026-04-30",
    "tradingsession": 3
  }
]
```

## Edge cases

- Пустой массив при заданном `--date` → выходной/праздник. Подбери ближайшую торговую дату через [`get_board_dates`](get_board_dates.md).
- Без `--date` вернётся **полная история** включений/исключений (длинная таблица).
