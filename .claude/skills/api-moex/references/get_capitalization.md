---
endpoint: /iss/statistics/engines/stock/capitalization.json
---

# `get_capitalization`

**Назначение:** суммарная капитализация фондового рынка (Blocks: capitalization — итог по всему рынку; issuecapitalization — суммарно по выпускам).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/capitalization>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/capitalization.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                  | Тип   | Смысл                     |
| --------------------- | ----- | ------------------------- |
| `CAPITALIZATION`      | float | Капитализация рынка, ₽    |
| `TRADEDATE`           | date  | Дата                      |
| `ISSUECAPITALIZATION` | float | Капитализация по выпускам |
| `UPDATETIME`          | str   | Время обновления          |
