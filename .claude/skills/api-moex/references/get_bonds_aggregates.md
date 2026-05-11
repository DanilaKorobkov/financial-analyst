---
endpoint: /iss/statistics/engines/stock/markets/bonds/aggregates.json
block: aggregates
---

# `get_bonds_aggregates`

**Назначение:** агрегаты рынка облигаций по типам (ОФЗ, корпоративные, субфедеральные).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/bonds/aggregates>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/bonds/aggregates.json
```

**Форма ответа:** блок `aggregates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Поля aggregates: `value, volume, num_trades,` в разрезе классов облигаций.
