---
endpoint: /iss/statistics/engines/stock/markets/index/rusfar.json
block: analytics
---

# `get_rusfar`

**Назначение:** RUSFAR (Russian Secured Funding Average Rate) — индикатор ставок денежного рынка по срокам.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/index/rusfar>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/index/rusfar.json
```

**Форма ответа:** блок `analytics`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                            | Тип   | Смысл                                                          |
| ------------------------------- | ----- | -------------------------------------------------------------- |
| `tradedate`                     | date  | Дата                                                           |
| `secid`                         | str   | RUSFAR (overnight), RUSFAR1W, RUSFAR1M, RUSFAR3M, RUSFARCNY, … |
| `Numtrades`                     | int   | Сделок                                                         |
| `AvgTradesVol`                  | int   | Средний объём                                                  |
| `MinTradePrice`/`MaxTradePrice` | float | Min/Max ставка, %                                              |
| `Vol`                           | int   | Общий объём                                                    |
| `MinVol`                        | int   | Минимальный объём для расчёта                                  |
| `VolShare`                      | float | Доля объёма                                                    |

## Edge cases

- Базовый proxy безрисковой rouble ставки на коротких сроках.
