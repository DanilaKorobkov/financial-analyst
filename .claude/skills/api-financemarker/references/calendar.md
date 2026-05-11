---
endpoint: /api/fm/v2/calendar
---

# `calendar`

**Назначение:** календарь **предстоящих** корпоративных событий (заседания СД, ГОСА, отсечки и т.п.). Без `updated_in_days`.

**URL:** `GET /api/fm/v2/calendar`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)), кроме `updated_in_days` — этот эндпоинт его не поддерживает.

## Пример запроса
```
GET /api/fm/v2/calendar?limit=5&sort_by=date&sort_order=asc
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Тикер и биржа |
| `category` | str | Категория события (`BOARD_MEETING`, …) |
| `event` | str | Описание события (одна строка) |
| `type` | str | Подтип (`NA` если не уточнён) |
| `period` | str | Период (`NA` для разовых) |
| `year` / `month` | int | Год / месяц периода (0/0 если не задано) |
| `date` | date | Дата события |
| `link` | url | Ссылка на источник раскрытия |

## Пример ответа
```json
[
  {
    "code": "ZILL",
    "exchange": "MOEX",
    "category": "BOARD_MEETING",
    "event": "Об утверждении состава Комитета Совета директоров по аудиту",
    "type": "NA",
    "period": "NA",
    "year": 0,
    "month": 0,
    "date": "2026-05-04",
    "link": "https://www.e-disclosure.ru/portal/event.aspx?EventId=UhWKOg2So0KNEKzmCG4b0A-B-B"
  }
]
```

## Edge cases
- Уже произошедшие события сюда не попадают — для них [`disclosure`](disclosure.md).
- Ключ `event` ≠ `title` (в `disclosure`) — это разные поля; не объединяй вслепую.
