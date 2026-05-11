---
endpoint: /iss/engines/{engine}/markets/{market}/boards.json
block: boards
---

# `get_boards`

**Назначение:** справочник режимов торгов внутри рынка.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards.json
```

**Форма ответа:** блок `boards`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Дефолт   | Описание |
| ----------------- | --- | -------- | -------- |
| `{engine}` (path) | str | `stock`  | Движок   |
| `{market}` (path) | str | `shares` | Рынок    |

## Поля JSON-ответа

| Поле             | Тип | Смысл                                    |
| ---------------- | --- | ---------------------------------------- |
| `id`             | int | ID режима                                |
| `board_group_id` | int | ID группы                                |
| `boardid`        | str | Код режима для `--board` (TQBR, EQBR, …) |
| `title`          | str | Название режима                          |
| `is_traded`      | int | 1 = активен                              |
| `has_candles`    | int | 1 = поддерживаются свечи                 |
| `is_primary`     | int | 1 = основной режим рынка                 |

## Edge cases

- Для `--market shares` основной режим всегда TQBR.
