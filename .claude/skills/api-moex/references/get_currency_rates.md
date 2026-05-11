---
endpoint: /iss/statistics/engines/currency/markets/selt/rates.json
---

# `get_currency_rates`

**Назначение:** курсы валют MOEX/ЦБ — текущие.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/currency/markets/selt/rates>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/currency/markets/selt/rates.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Объединение блоков `cbrf` (курс ЦБ + дельта) и `wap_rates` (взвешенный курс торгов).

| Поле                     | Тип   | Смысл                        |
| ------------------------ | ----- | ---------------------------- |
| `CBRF_USD_LAST`          | float | Курс ЦБ USD/RUB              |
| `CBRF_EUR_LAST`          | float | Курс ЦБ EUR/RUB              |
| `CBRF_*_LASTCHANGEPRCNT` | float | Изменение к предыдущему      |
| `USDTOM_UTS_CLOSEPRICE`  | float | Закрытие USDRUB_TOM с торгов |

В блоке `wap_rates` — поля по каждой валюте: `secid, shortname, price, lasttoprevprice, nominal, decimals`.

## Edge cases

- Полезно для пересчёта валютной выручки экспортёров (LKOH, GMKN, NLMK, PHOR).
