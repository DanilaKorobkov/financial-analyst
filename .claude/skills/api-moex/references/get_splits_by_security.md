---
endpoint: /iss/statistics/engines/stock/splits/{TICKER}.json
block: splits
---

# `get_splits_by_security`

**Назначение:** сплиты по конкретной бумаге (если были).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/splits/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/splits/{TICKER}.json
```

**Форма ответа:** блок `splits`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `splits`: `tradedate, secid, before, after`.

## Edge cases

- Сервер фильтрует по тикеру (быстрее, чем `get_splits_by_security`).
- Бумага без сплитов → пустой массив `data` в блоке.
