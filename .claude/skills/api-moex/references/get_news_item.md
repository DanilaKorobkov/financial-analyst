---
endpoint: /iss/sitenews/{news_id}.json
---

# `get_news_item`

**Назначение:** полный текст новости (HTML-тело) по ID.

**Reference ISS:** <https://iss.moex.com/iss/sitenews/{ID}>

## URL

```http
GET https://iss.moex.com/iss/sitenews/{news_id}.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Объединение блоков (с колонкой `_block`). Основной блок `content`: `id, title, published_at, body` (body — HTML).

## Edge cases

- ID берётся из `get_sitenews`.
