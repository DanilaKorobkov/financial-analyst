---
endpoint: /iss/history/engines/{engine}/markets/{market}/sessions.json
block: trading_sessions
---

# `get_market_sessions`

**Назначение:** список торговых сессий рынка.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/sessions>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/sessions.json
```

**Форма ответа:** блок `trading_sessions`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле    | Тип | Смысл                                                                |
| ------- | --- | -------------------------------------------------------------------- |
| `id`    | int | 0/1/2/3/5                                                            |
| `name`  | str | morning/main/evening/total/weekend                                   |
| `title` | str | Утренняя/Основная/Вечерняя/Итого/Дополнительная сессия выходного дня |

## Edge cases

- `id=3` (total) — агрегат за день; используй для расчётов, исключающих дробление по сессиям.
