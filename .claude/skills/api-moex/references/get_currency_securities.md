---
endpoint: /iss/engines/currency/markets/selt/securities.json
---

# `get_currency_securities`

**Назначение:** все валютные пары рынка SELT с marketdata (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/engines/currency/markets/selt/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/currency/markets/selt/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Объединение блоков `securities` (статика инструмента) и `marketdata` (последние котировки) с колонкой `_block`.

## Edge cases

- Основные пары: USD000UTSTOM, EUR_RUB\_\_TOM, CNYRUB_TOM, GLDRUB_TOM (золото).
