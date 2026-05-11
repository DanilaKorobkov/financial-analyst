---
endpoint: /iss/engines.json
block: engines
---

# `get_engines`

**Назначение:** все торговые системы MOEX (stock, currency, futures, state, agro, otc, money, ...).

**Reference ISS:** <https://iss.moex.com/iss/engines>

## URL

```http
GET https://iss.moex.com/iss/engines.json
```

**Форма ответа:** блок `engines`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле    | Тип | Смысл                                |
| ------- | --- | ------------------------------------ |
| `id`    | int | Числовой ID системы                  |
| `name`  | str | Код для подстановки в URL (`engine`) |
| `title` | str | Человекочитаемое название            |

## Edge cases

- Используй `name` как значение `--engine` в других командах.
