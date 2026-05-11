---
endpoint: /iss/history/engines/{engine}/markets/{market}/boards/{board}/listing.json
block: securities
paginated: true
---

# `get_board_listing`

**Назначение:** инструменты конкретного режима с history_from/till.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/listing>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/listing.json
```

**Форма ответа:** блок `securities` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_market_listing` — те же поля, но только для одного режима.

## Edge cases

- Для TQBR: ~3000+ инструментов с момента запуска режима.
