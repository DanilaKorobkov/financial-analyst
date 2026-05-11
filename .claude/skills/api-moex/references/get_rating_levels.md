---
endpoint: /iss/cci/reference/rating-levels.json
block: rating_levels
paginated: true
---

# `get_rating_levels`

**Назначение:** уровни рейтингов (~800 записей с описаниями).

**Reference ISS:** <https://iss.moex.com/iss/cci/reference/rating-levels>

## URL

```http
GET https://iss.moex.com/iss/cci/reference/rating-levels.json
```

**Форма ответа:** блок `rating_levels` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                       | Тип | Смысл                                          |
| -------------------------- | --- | ---------------------------------------------- |
| `rating_level_id`          | int | ID уровня                                      |
| `name_short_ru`            | str | Краткий код (AAA(RU), BBB+(RU), …)             |
| `rating_group_code`        | str | Группа (AAA, AA, A, BBB, BB, B, CCC, CC, C, D) |
| `rating_level_description` | str | Полное описание уровня                         |
| `rating_book_id`           | int | ID шкалы (см. `get_rating_books`)              |

## Edge cases

- Пагинация автоматическая.
