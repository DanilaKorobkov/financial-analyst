---
endpoint: /api/fm/v2/ideas
---

# `ideas`

**Назначение:** лента инвест-идей аналитиков **по всему рынку**. По одной бумаге — раздел `ideas` в [`stocks/{exchange}:{code}`](stock.md). По одному `id` — [`ideas/{id}`](idea.md).

**URL:** `GET /api/fm/v2/ideas`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)). Полезные значения `sort_by`: `date_in`, `profit_potential`, `changed_at`.

## Пример запроса
```
GET /api/fm/v2/ideas?limit=5&sort_by=profit_potential&sort_order=desc
GET /api/fm/v2/ideas?updated_in_days=7&limit=20
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `id` | int | ID идеи (нужен для [`ideas/{id}`](idea.md)) |
| `code` / `exchange` | str | Тикер бумаги, по которой идея |
| `community` / `community_id` | str / int | Брокер/аналитик-автор |
| `idea` | str | Заголовок идеи (одна строка) |
| `date_in` / `date_out` | date | Дата публикации / запланированный выход |
| `duration_in_month` | int | Срок идеи в месяцах |
| `price_in` / `price_out` | float | Цена входа / целевая |
| `price_day` | float | Текущая цена бумаги на момент снимка |
| `profit_potential` | float | Потенциал к target в % |
| `profit_actual` | float | Текущая фактическая доходность от `price_in` (%) |
| `stop_loss` | float | Стоп-лосс (если задан автором) |
| `system_status` | str | `ACTIVE` / `CLOSED` / и т.п. |
| `close_date` / `close_price` / `close_comment` / `close_link` | mixed | Поля закрытия (для `CLOSED`) |
| `update_date` / `update_price` | mixed | Дата и цена последнего апдейта target-а |
| `changed_at` | datetime | Время последнего изменения записи |

## Пример ответа
```json
[
  {
    "id": 6211,
    "code": "YDEX",
    "exchange": "MOEX",
    "community": "Велес Капитал",
    "community_id": 21111,
    "idea": "Продажа сервиса объявлений Авто.ру",
    "date_in": "2026-03-18",
    "date_out": "2027-03-18",
    "duration_in_month": 12,
    "price_in": 4500.0,
    "price_out": 5727.0,
    "price_day": 4060.0,
    "profit_potential": 27.0,
    "profit_actual": -9.8,
    "system_status": "ACTIVE",
    "close_date": "2026-04-30",
    "close_price": 4060.0,
    "changed_at": "2026-05-01T01:40:03"
  }
]
```

## Edge cases
- `close_date`/`close_price` могут быть заполнены даже для `system_status=ACTIVE` — это снимок «на дату последнего обновления цены», а не закрытие идеи.
- Для детального текста идеи (`description` с тезисами) нужен отдельный вызов [`ideas/{id}`](idea.md).
