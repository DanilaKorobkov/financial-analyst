---
endpoint: /iss/history/engines/stock/totals/securities.json
block: totals
paginated: true
---

# `get_totals`

**Назначение:** итоги по всем выпускам (с пагинацией) — DAILYCAPITALIZATION / MONTHLYCAPITALIZATION.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/stock/totals/securities>

## URL

```http
GET https://iss.moex.com/iss/history/engines/stock/totals/securities.json
```

**Форма ответа:** блок `totals` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                    | Тип   | Смысл                  |
| ----------------------- | ----- | ---------------------- |
| `SECID`                 | str   | Тикер                  |
| `OPEN/HIGH/LOW/CLOSE`   | float | OHLC                   |
| `DAILYCAPITALIZATION`   | float | Капитализация за день  |
| `MONTHLYCAPITALIZATION` | float | Капитализация за месяц |

## Edge cases

- Большой набор данных, пагинация автоматическая.
