---
endpoint: /iss/engines/{engine}.json
---

# `get_engine`

**Назначение:** описание одной торговой системы со списком её рынков и режимов (Blocks: engine, markets, boards).

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя      | Тип | Обязательно | Описание                                   |
| -------- | --- | ----------- | ------------------------------------------ |
| `ENGINE` | str | да          | `stock`, `currency`, `futures`, `state`, … |

## Поля JSON-ответа

CSV с дополнительной колонкой `_block` (engine / markets / boards).

## Edge cases

- Несуществующий движок → пустой массив `data` в блоке.
