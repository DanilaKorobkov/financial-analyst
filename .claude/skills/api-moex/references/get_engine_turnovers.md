---
endpoint: /iss/engines/{engine}/turnovers.json
---

# `get_engine_turnovers`

**Назначение:** обороты одной торговой системы по её рынкам.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/turnovers>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/turnovers.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

См. `get_all_turnovers`.
