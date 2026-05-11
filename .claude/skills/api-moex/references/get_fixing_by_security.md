---
endpoint: /iss/statistics/engines/currency/markets/fixing/{TICKER}.json
block: history
paginated: true
---

# `get_fixing_by_security`

**Назначение:** история фиксинга конкретной валютной пары.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/currency/markets/fixing/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/currency/markets/fixing/{TICKER}.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Обязательно | Описание       |
| ----------------- | ---- | ----------- | -------------- |
| `{TICKER}` (path) | str  | да          | SECID фиксинга |
| `from`            | date | нет         | Начало периода |
| `till`            | date | нет         | Конец периода  |

## Поля JSON-ответа

См. `get_fixing` (поля tradedate, secid, rate).

## Edge cases

- В отличие от `get_fixing`, тут пагинация работает корректно.
