---
endpoint: /api/fm/v2/disclosure
---

# `disclosure`

**Назначение:** лента раскрытия корпоративной информации эмитентами (существенные факты, заседания СД и т.п.). Без `updated_in_days`.

**URL:** `GET /api/fm/v2/disclosure`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)), кроме `updated_in_days` — этот эндпоинт его не поддерживает.

## Пример запроса
```
GET /api/fm/v2/disclosure?limit=5&sort_by=date&sort_order=desc
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Тикер и биржа |
| `category` | str | Категория события (`EVENT`, …) |
| `title` | str | Заголовок раскрытия |
| `type` | str | Подтип (`NA` — не специфицирован) |
| `period` | str | Период отчётности (`NA` для событий без периода) |
| `year` / `month` | int | Год / месяц периода (0/0 = не задано) |
| `date` | date | Дата события |
| `link` | url | Ссылка на e-disclosure / источник |

## Пример ответа
```json
[
  {
    "code": "CNRU",
    "exchange": "MOEX",
    "category": "EVENT",
    "title": "Проведение заседания совета директоров (наблюдательного совета) и его повестка дня",
    "type": "NA",
    "period": "NA",
    "year": 0,
    "month": 0,
    "date": "2026-05-04",
    "link": "https://www.e-disclosure.ru/portal/company.aspx?id=39286"
  }
]
```

## Edge cases
- В отличие от [`calendar`](calendar.md) (что **будет**), это лента уже **раскрытых** событий.
- Для запланированных событий с известной датой проведения смотри [`calendar`](calendar.md).
