---
endpoint: /iss/events/{event_id}.json
---

# `get_event`

**Назначение:** детали мероприятия по ID.

**Reference ISS:** <https://iss.moex.com/iss/events/{ID}>

## URL

```http
GET https://iss.moex.com/iss/events/{event_id}.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Объединение блоков события (с колонкой `_block`).
