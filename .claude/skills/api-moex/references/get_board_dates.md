---
endpoint: /iss/history/engines/{engine}/markets/{market}/boards/{board}/dates.json
block: dates
---

# `get_board_dates`

**Назначение:** узнать диапазон дат, доступных в **истории** для рынка по конкретному режиму торгов. Используется, чтобы подобрать ближайшую торговую дату при пустом ответе других методов.

**Reference ISS:** <https://iss.moex.com/iss/reference/26>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/boards/{board}/dates.json
```

**Форма ответа:** блок `dates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт   | Описание     |
| ----------------- | --- | ----------- | -------- | ------------ |
| `{board}` (path)  | str | нет         | `TQBR`   | Режим торгов |
| `{market}` (path) | str | нет         | `shares` | Рынок        |
| `{engine}` (path) | str | нет         | `stock`  | Движок       |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив **из одного** объекта:

| Поле   | Тип  | Обязательно | Смысл                                   |
| ------ | ---- | ----------- | --------------------------------------- |
| `from` | date | да          | Самая ранняя дата истории на режиме     |
| `till` | date | да          | Самая поздняя (последняя торговая дата) |

## Пример ответа

```json
[
  {
    "from": "2013-03-25",
    "till": "2026-04-30"
  }
]
```
