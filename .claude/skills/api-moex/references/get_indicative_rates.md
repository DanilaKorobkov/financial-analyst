---
endpoint: /iss/statistics/engines/futures/markets/indicativerates/securities.json
---

# `get_indicative_rates`

**Назначение:** индикативные курсы валют срочного рынка (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/futures/markets/indicativerates/securities>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/futures/markets/indicativerates/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Блок `securities`: tradedate, tradetime, secid, rate, clearing — индикативные курсы по парам (CAD/RUB и др.).
