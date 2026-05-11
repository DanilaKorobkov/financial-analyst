---
endpoint: /iss/events.json
---

# `get_events`

**Назначение:** активные мероприятия биржи.

**Reference ISS:** <https://iss.moex.com/iss/events>

## URL

```http
GET https://iss.moex.com/iss/events.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Объединение блоков. Обычно пусто (события публикуются по графику).
