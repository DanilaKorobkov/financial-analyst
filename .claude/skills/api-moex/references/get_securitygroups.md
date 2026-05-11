---
endpoint: /iss/securitygroups.json
block: securitygroups
---

# `get_securitygroups`

**Назначение:** глобальные группы инструментов (stock_shares, stock_bonds, stock_index, currency_selt, futures_forts, ...).

**Reference ISS:** <https://iss.moex.com/iss/securitygroups>

## URL

```http
GET https://iss.moex.com/iss/securitygroups.json
```

**Форма ответа:** блок `securitygroups`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле        | Тип | Смысл                                                                       |
| ----------- | --- | --------------------------------------------------------------------------- |
| `id`        | int | ID группы                                                                   |
| `name`      | str | Код для `--securitygroup` в `get_collections` / `get_collection_securities` |
| `title`     | str | Название группы                                                             |
| `is_hidden` | int | Скрыта в UI                                                                 |

## Edge cases

- Используй `name` (например `stock_shares`) для вложенных запросов коллекций.
