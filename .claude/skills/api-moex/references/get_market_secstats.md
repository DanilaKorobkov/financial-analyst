---
endpoint: /iss/engines/{engine}/markets/{market}/secstats.json
block: secstats
---

# `get_market_secstats`

**Назначение:** промежуточные итоги дня по всем бумагам рынка в разрезе сессий.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/secstats>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/secstats.json
```

**Форма ответа:** блок `secstats`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                  | Тип   | Смысл                                               |
| --------------------- | ----- | --------------------------------------------------- |
| `SECID`               | str   | Тикер                                               |
| `BOARDID`             | str   | Режим                                               |
| `TRADINGSESSION`      | int   | 0 = утренняя, 1 = основная, 2 = вечерняя, 3 = итого |
| `TIME`                | str   | Время отчёта                                        |
| `OPEN/HIGH/LOW/LAST`  | float | Дневные O/H/L и последняя сделка в сессии           |
| `WAPRICE`             | float | WAP сессии                                          |
| `MARKETPRICE2`        | float | Альтернативный WAP MOEX                             |
| `LCURRENTPRICE`       | float | Текущая цена закрытия (моментальная)                |
| `CLOSINGAUCTIONPRICE` | float | Цена аукциона закрытия                              |
| `VOLTODAY/VALTODAY`   | int   | Накопленный объём за сессию                         |
| `NUMTRADES`           | int   | Количество сделок                                   |

## Edge cases

- Обновляется раз в N минут. Полезно для intraday-обзора.
