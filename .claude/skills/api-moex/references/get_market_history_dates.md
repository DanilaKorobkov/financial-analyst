---
endpoint: /iss/history/engines/{engine}/markets/{market}/securities/{TICKER}/dates.json
block: dates
---

# `get_market_history_dates`

**Назначение:** диапазон дат истории конкретной бумаги по всем режимам рынка.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/securities/{secid}/dates>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/securities/{TICKER}/dates.json
```

**Форма ответа:** блок `dates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле   | Тип  | Смысл                       |
| ------ | ---- | --------------------------- |
| `from` | date | Первая торговая дата бумаги |
| `till` | date | Последняя торговая дата     |

## Edge cases

- Используй перед длинным запросом `get_market_history` — узнать реальные границы.
