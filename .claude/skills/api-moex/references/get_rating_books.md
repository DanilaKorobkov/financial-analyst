---
endpoint: /iss/cci/reference/rating-books.json
block: rating_books
---

# `get_rating_books`

**Назначение:** справочник рейтинговых шкал (АКРА, Эксперт РА, НКР, НРА, ЦБ — национальные и международные).

**Reference ISS:** <https://iss.moex.com/iss/cci/reference/rating-books>

## URL

```http
GET https://iss.moex.com/iss/cci/reference/rating-books.json
```

**Форма ответа:** блок `rating_books`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле                   | Тип | Смысл                                |
| ---------------------- | --- | ------------------------------------ |
| `rating_book_id`       | int | ID шкалы                             |
| `name_short_ru`        | str | Короткое имя                         |
| `name_full_ru`         | str | Полное имя                           |
| `agency_id`            | int | ID агентства                         |
| `agency_name_short_ru` | str | Агентство (АКРА, Эксперт РА, НКР, …) |

## Edge cases

- ~50 шкал. Сами рейтинги конкретных эмитентов/выпусков — `/iss/cci/rating/...`, **платно**.
