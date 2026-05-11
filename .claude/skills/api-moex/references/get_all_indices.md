---
endpoint: /iss/statistics/engines/stock/markets/index/analytics.json
block: indices
---

# `get_all_indices`

**Назначение:** список всех индексов MOEX.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics.json
```

**Форма ответа:** блок `indices`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле        | Тип  | Смысл                     |
| ----------- | ---- | ------------------------- |
| `indexid`   | str  | Код индекса               |
| `shortname` | str  | Короткое название         |
| `from`      | date | Дата запуска              |
| `till`      | date | Последняя актуальная дата |

## Edge cases

- ~290 индексов: основные (IMOEX, MOEX10, MOEXBC), отраслевые субиндексы, индексы полной доходности.
