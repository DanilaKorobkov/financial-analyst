---
endpoint: /iss/securitygroups/{securitygroup}/collections.json
block: collections
---

# `get_collections`

**Назначение:** коллекции внутри группы инструментов (например, для `stock_index`: `stock_index_all`, `stock_index_shares`, `stock_index_shares_sectoral`, `stock_index_total_return`).

**Reference ISS:** <https://iss.moex.com/iss/securitygroups/{securitygroup}/collections>

## URL

```http
GET https://iss.moex.com/iss/securitygroups/{securitygroup}/collections.json
```

**Форма ответа:** блок `collections`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя             | Тип | Обязательно | Описание                              |
| --------------- | --- | ----------- | ------------------------------------- |
| `SECURITYGROUP` | str | да          | Имя группы (см. `get_securitygroups`) |

## Поля JSON-ответа

| Поле    | Тип | Смысл                               |
| ------- | --- | ----------------------------------- |
| `id`    | int | ID коллекции                        |
| `name`  | str | Код для `get_collection_securities` |
| `title` | str | Название                            |

## Edge cases

- Несуществующая группа → пустой массив `data` в блоке.
