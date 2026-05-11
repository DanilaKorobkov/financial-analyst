---
endpoint: /iss/referencedata/engines/{engine}/markets/all/securitieslisting.json
---

# `get_securities_listing`

**Назначение:** справочник торгуемости бумаг (referencedata).

**Reference ISS:** <https://iss.moex.com/iss/referencedata/engines/{engine}/markets/all/securitieslisting>

## URL

```http
GET https://iss.moex.com/iss/referencedata/engines/{engine}/markets/all/securitieslisting.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Дефолт  | Описание |
| ----------------- | --- | ------- | -------- |
| `{engine}` (path) | str | `stock` | Движок   |

## Поля JSON-ответа

Возвращает все блоки эндпоинта (с колонкой `_block`). Основной блок — `data`.

## Edge cases

- Похоже на `get_all_stock_securities`, но в формате referencedata-справочника.
