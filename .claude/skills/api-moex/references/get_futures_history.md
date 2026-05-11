---
endpoint: /iss/history/engines/futures/markets/forts/securities/{TICKER}.json
block: history
paginated: true
---

# `get_futures_history`

**Назначение:** история одного фьючерса.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/futures/markets/forts/securities/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/history/engines/futures/markets/forts/securities/{TICKER}.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Блок `history`: BOARDID, TRADEDATE, OPEN/HIGH/LOW/CLOSE, OPENPOSITION, NUMTRADES, VOLUME, VALUE, …

## Edge cases

- Для экспирированных контрактов история обрывается на дате экспирации.
