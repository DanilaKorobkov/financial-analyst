---
endpoint: /iss/statistics/engines/stock/currentprices.json
block: currentprices
---

# `get_currentprices`

**Назначение:** история текущих цен в разрезе сессий (CURPRICE/LASTPRICE/LEGALCLOSE).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/currentprices>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/currentprices.json
```

**Форма ответа:** блок `currentprices`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа (блок `currentprices`)

| Поле             | Тип   | Смысл                     |
| ---------------- | ----- | ------------------------- |
| `TRADEDATE`      | date  | Дата                      |
| `BOARDID`        | str   | Режим                     |
| `SECID`          | str   | Тикер                     |
| `TRADETIME`      | str   | Время                     |
| `CURPRICE`       | float | Текущая цена              |
| `LASTPRICE`      | float | Последняя сделка          |
| `LEGALCLOSE`     | float | Юридическая цена закрытия |
| `TRADINGSESSION` | int   | Сессия                    |

## Edge cases

- Single page.
