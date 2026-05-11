---
endpoint: /iss/engines/futures/markets/options/securities.json
---

# `get_options_securities`

**Назначение:** активные опционы FORTS (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/engines/futures/markets/options/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/futures/markets/options/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Аналог `get_futures_securities` — securities + marketdata.

## Edge cases

- Большой объём данных (тысячи опционных контрактов).
