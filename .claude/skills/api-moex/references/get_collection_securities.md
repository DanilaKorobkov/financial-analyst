---
endpoint: /iss/securitygroups/{securitygroup}/collections/{collection}/securities.json
block: securities
paginated: true
---

# `get_collection_securities`

**Назначение:** все бумаги коллекции (с пагинацией).

**Reference ISS:** <https://iss.moex.com/iss/securitygroups/{securitygroup}/collections/{collection}/securities>

## URL

```http
GET https://iss.moex.com/iss/securitygroups/{securitygroup}/collections/{collection}/securities.json
```

**Форма ответа:** блок `securities` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание                           |
| ----------------- | --- | ----------- | ---------------------------------- |
| `SECURITYGROUP`   | str | да          | См. `get_securitygroups`           |
| `COLLECTION`      | str | да          | См. `get_collections`              |
| `<block>.columns` | csv | нет         | Подмножество колонок (пусто = все) |

## Поля JSON-ответа

Список бумаг (SECID, ISIN, REGNUMBER, NAME, BOARDID, ...) — поля зависят от группы.

## Edge cases

- Коллекция `stock_shares_one` / `_two` / `_three` — бумаги из соответствующего уровня листинга.
- `stock_index_all` — все индексы MOEX как «бумаги» индексного рынка.
