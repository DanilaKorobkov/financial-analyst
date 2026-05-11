---
endpoint: /iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/dates.json
block: dates
---

# `get_board_history_dates`

**Назначение:** диапазон дат истории бумаги в одном режиме.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{secid}/dates>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/securities/{TICKER}/dates.json
```

**Форма ответа:** блок `dates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_market_history_dates`.

## Edge cases

- Полезно при кросс-проверке: TQBR vs SMAL — у бумаги может быть разная глубина истории.
