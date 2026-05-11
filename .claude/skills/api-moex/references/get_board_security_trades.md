---
endpoint: /iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/trades.json
block: trades
---

# `get_board_security_trades`

**Назначение:** последние сделки по бумаге в конкретном режиме.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities/{secid}/trades>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/trades.json
```

**Форма ответа:** блок `trades`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_market_trades`.

## Edge cases

- В TQBR без аукционов закрытия — чистый поток сделок основной сессии.
