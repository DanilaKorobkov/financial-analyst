---
endpoint: /iss/engines/{engine}/markets/{market}/trades.json
block: trades
---

# `get_market_trades`

**Назначение:** последние сделки всех бумаг рынка (~50 свежих).

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/trades>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/trades.json
```

**Форма ответа:** блок `trades`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле             | Тип   | Смысл                                  |
| ---------------- | ----- | -------------------------------------- |
| `TRADENO`        | int   | Номер сделки                           |
| `TRADETIME`      | str   | Время сделки                           |
| `BOARDID`        | str   | Режим                                  |
| `SECID`          | str   | Тикер                                  |
| `PRICE`          | float | Цена                                   |
| `QUANTITY`       | int   | Количество (лоты)                      |
| `VALUE`          | float | Сумма                                  |
| `BUYSELL`        | str   | `B` = инициатор покупка, `S` = продажа |
| `TRADINGSESSION` | int   | Сессия                                 |
| `TRADEDATE`      | date  | Дата                                   |

## Edge cases

- Стаканы (orderbook) — платные, в CLI не реализованы.
