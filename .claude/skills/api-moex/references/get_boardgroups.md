---
endpoint: /iss/engines/{engine}/markets/{market}/boardgroups.json
block: boardgroups
---

# `get_boardgroups`

**Назначение:** группы режимов торгов рынка.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/boardgroups>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boardgroups.json
```

**Форма ответа:** блок `boardgroups`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле             | Тип | Смысл                                                       |
| ---------------- | --- | ----------------------------------------------------------- |
| `id`             | int | ID группы (используй как `boardgroup` в зависимых командах) |
| `name`           | str | Код группы                                                  |
| `title`          | str | Название («Т+: Акции и ДР», «Аукцион», …)                   |
| `is_default`     | int | 1 = группа по умолчанию для рынка                           |
| `board_group_id` | int | Альтернативный ID                                           |

## Edge cases

- ID `57` — основная группа `Т+: Акции и ДР - безадрес.` для рынка shares.
