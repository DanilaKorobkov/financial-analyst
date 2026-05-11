---
endpoint: /iss/engines/{engine}/markets/{market}/securities/{TICKER}/trades.json
block: trades
---

# `get_security_trades`

**Назначение:** последние сделки по одной бумаге.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{secid}/trades>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{TICKER}/trades.json
```

**Форма ответа:** блок `trades`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_market_trades`, фильтр по `SECID`.

## Edge cases

- Запрос свежий только в торговое время.
