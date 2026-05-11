---
endpoint: /iss/engines/futures/markets/forts/securities.json
---

# `get_futures_securities`

**Назначение:** активные фьючерсы FORTS с marketdata (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/engines/futures/markets/forts/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/futures/markets/forts/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

`securities` (статика контракта) + `marketdata` (котировки) с колонкой `_block`.

## Edge cases

- Контракты помечены кодом серии (например, SiM6 = USD/RUB Mar-2026).
