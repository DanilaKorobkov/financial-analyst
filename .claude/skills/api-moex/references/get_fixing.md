---
endpoint: /iss/statistics/engines/currency/markets/fixing.json
block: history
---

# `get_fixing`

**Назначение:** биржевые валютные фиксинги (текущий день).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/currency/markets/fixing>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/currency/markets/fixing.json
```

**Форма ответа:** блок `history`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле        | Тип   | Смысл                          |
| ----------- | ----- | ------------------------------ |
| `tradedate` | date  | Дата                           |
| `secid`     | str   | Код фиксинга (USD000UTSFIX, …) |
| `rate`      | float | Значение фиксинга              |

## Edge cases

- ISS **игнорирует** `from`/`till` — отдаёт только текущий день. Для истории — `get_fixing_by_security`.
