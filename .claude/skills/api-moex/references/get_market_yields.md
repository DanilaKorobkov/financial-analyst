---
endpoint: /iss/history/engines/stock/markets/bonds/yields/{TICKER}.json
block: history_yields
paginated: true
---

# `get_market_yields`

**Назначение:** доходности облигации **по всем режимам сразу**.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/stock/markets/bonds/yields/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/history/engines/stock/markets/bonds/yields/{TICKER}.json
```

**Форма ответа:** блок `history_yields` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_bond_yields` — те же поля, но без фильтра по режиму.

## Edge cases

- Бумага может торговаться в нескольких режимах (TQOB, EQOB, …).
