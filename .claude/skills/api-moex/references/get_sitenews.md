---
endpoint: /iss/sitenews.json
block: sitenews
---

# `get_sitenews`

**Назначение:** свежие новости биржи (~50 последних, single page).

**Reference ISS:** <https://iss.moex.com/iss/sitenews>

## URL

```http
GET https://iss.moex.com/iss/sitenews.json
```

**Форма ответа:** блок `sitenews`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя    | Тип | Дефолт | Описание  |
| ------ | --- | ------ | --------- |
| `lang` | str | `ru`   | `ru`/`en` |

## Поля JSON-ответа

| Поле           | Тип | Смысл             |
| -------------- | --- | ----------------- |
| `id`           | int | ID новости        |
| `tag`          | str | Тег (site и т.п.) |
| `title`        | str | Заголовок         |
| `published_at` | str | Дата публикации   |
| `modified_at`  | str | Дата изменения    |

## Edge cases

- Полный текст новости — `get_news_item <ID>`.
- Архив > 50 не подгружаем (single page по дизайну).
