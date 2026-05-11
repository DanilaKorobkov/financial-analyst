---
endpoint: /iss/index.json
---

# `get_reference`

**Назначение:** получить перечень допустимых значений `PLACEHOLDER` в URL ISS — нужно когда строишь произвольный путь и не помнишь, какие `engine` / `market` / `board` существуют.

**Reference ISS:** <https://iss.moex.com/iss/reference/28>

## URL

```http
GET https://iss.moex.com/iss/index.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя           | Тип | Обязательно | Дефолт   | Описание                                                                                                                                       |
| ------------- | --- | ----------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `PLACEHOLDER` | str | нет         | `boards` | Что выгрузить. Допустимо: `engines`, `markets`, `boards`, `boardgroups`, `durations`, `securitytypes`, `securitygroups`, `securitycollections` |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов. Поля **зависят от `PLACEHOLDER`**. Ниже — для дефолтного `boards`:

| Поле             | Тип | Обязательно | Смысл                          |
| ---------------- | --- | ----------- | ------------------------------ |
| `id`             | int | да          | ID режима                      |
| `board_group_id` | int | да          | Группа режимов                 |
| `engine_id`      | int | да          | ID движка                      |
| `market_id`      | int | да          | ID рынка                       |
| `boardid`        | str | да          | Код режима (`TQBR`, `TQTF`, …) |
| `board_title`    | str | да          | Название режима                |
| `is_traded`      | int | да          | 1 = торгуется                  |
| `has_candles`    | int | да          | 1 = есть свечи                 |
| `is_primary`     | int | да          | 1 = основной режим для рынка   |

Для остальных значений `PLACEHOLDER` (`engines`, `markets`, …) набор полей другой.

## Пример ответа (`boards`)

```json
[
  {
    "id": 177,
    "board_group_id": 57,
    "engine_id": 1,
    "market_id": 1,
    "boardid": "TQIF",
    "board_title": "Т+: Паи - безадрес.",
    "is_traded": 1,
    "has_candles": 1,
    "is_primary": 1
  },
  {
    "id": 178,
    "board_group_id": 57,
    "engine_id": 1,
    "market_id": 1,
    "boardid": "TQTF",
    "board_title": "Т+: ETF - безадрес.",
    "is_traded": 1,
    "has_candles": 1,
    "is_primary": 1
  },
  {
    "id": 129,
    "board_group_id": 57,
    "engine_id": 1,
    "market_id": 1,
    "boardid": "TQBR",
    "board_title": "Т+: Акции и ДР - безадрес.",
    "is_traded": 1,
    "has_candles": 1,
    "is_primary": 1
  }
]
```
