---
endpoint: /iss/engines/{engine}/markets/{market}/securities.json
---

# `get_market_securities`

**Назначение:** все бумаги рынка одним вызовом — справочник + текущие котировки (Blocks: securities + marketdata).

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

CSV со строками из обоих блоков (колонка `_block`):

- `securities` — статика (LOTSIZE, FACEVALUE, ISSUESIZE, ISIN, …);
- `marketdata` — котировки (LAST, BID, OFFER, WAPRICE, VOLTODAY, …).

## Edge cases

- Полезно для скрининга «дай все акции с current price».
- Полный ответ ~250–300 строк × 50+ колонок.
