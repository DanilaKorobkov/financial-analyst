---
endpoint: /iss/statistics/engines/futures/markets/options/assets.json
---

# `get_options_assets`

**Назначение:** опционные серии по всем активам (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/futures/markets/options/assets>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/futures/markets/options/assets.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Поля по активам: tradedate, asset, asset_type, asset_last_price, valtoday, voltoday, numtrades, openposition, oichange, …
