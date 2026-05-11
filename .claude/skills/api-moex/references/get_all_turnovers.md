---
endpoint: /iss/turnovers.json
---

# `get_all_turnovers`

**Назначение:** сводные обороты всех рынков MOEX (Blocks: turnovers — сегодня; turnoversprevdate — вчера).

**Reference ISS:** <https://iss.moex.com/iss/turnovers>

## URL

```http
GET https://iss.moex.com/iss/turnovers.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

CSV с колонкой `_block`. Поля: `NAME, ID, VALTODAY, VALTODAY_USD, NUMTRADES, UPDATETIME, TITLE`.

## Edge cases

- Содержит итоговую строку `TOTALS`.
- Для денежных оценок состояния рынка / сравнения дней.
