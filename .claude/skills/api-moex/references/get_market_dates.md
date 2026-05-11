---
endpoint: /iss/history/engines/{engine}/markets/{market}/dates.json
block: dates
---

# `get_market_dates`

**Назначение:** диапазон дат истории всего рынка (`from`..`till`).

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/dates>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/dates.json
```

**Форма ответа:** блок `dates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле   | Тип  | Смысл                   |
| ------ | ---- | ----------------------- |
| `from` | date | Самая ранняя дата       |
| `till` | date | Последняя торговая дата |

## Edge cases

- Для `--market shares`: с 1997-03-24 (старт MMVB) по сегодня.
