---
endpoint: /iss/engines/stock/markets/bonds/securities.json
---

# `get_bonds_securities`

**Назначение:** все облигации с текущими котировками (Blocks: securities + marketdata + marketdata_yields).

**Reference ISS:** <https://iss.moex.com/iss/engines/stock/markets/bonds/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/stock/markets/bonds/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Три блока:

- `securities` — статика (ISIN, FACEVALUE, MATDATE, COUPONPERIOD, …);
- `marketdata` — котировки (LAST, BID, OFFER, …);
- `marketdata_yields` — текущая доходность и дюрация (YIELDLAST, EFFECTIVEYIELD, DURATION).

## Edge cases

- Купоны/оферты/амортизации — нет в публичном ISS, бери из FM или e-disclosure.
