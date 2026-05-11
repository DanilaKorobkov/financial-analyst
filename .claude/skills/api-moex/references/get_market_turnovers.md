---
endpoint: /iss/engines/{engine}/markets/{market}/turnovers.json
block: turnovers
---

# `get_market_turnovers`

**Назначение:** обороты одного рынка.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/turnovers>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/turnovers.json
```

**Форма ответа:** блок `turnovers`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_all_turnovers`.
