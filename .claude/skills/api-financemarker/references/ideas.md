---
endpoint: /api/fm/v2/ideas
---

# `ideas`

**Назначение:** лента инвест-идей аналитиков **по всему рынку**. По одной бумаге — раздел `ideas` в [`stocks/{exchange}:{code}`](stock.md). По одному `id` — [`ideas/{id}`](idea.md).

**URL:** `GET /api/fm/v2/ideas`

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `limit` / `offset` / `sort_by` / `sort_order` / `updated_in_days` | — | optional | Стандартный набор пагинации. Полезные значения `sort_by`: `date_in`, `profit_potential`, `changed_at`. |

## Пример запроса
```
GET /api/fm/v2/ideas?api_token=$FINANCE_MARKER_TOKEN&limit=5&sort_by=profit_potential&sort_order=desc
GET /api/fm/v2/ideas?api_token=$FINANCE_MARKER_TOKEN&updated_in_days=7&limit=20
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `id` | int | Идентификатор идеи в FM. Нужен для [`ideas/{id}`](idea.md). |
| `code` | str | Тикер бумаги, по которой выставлена идея. |
| `exchange` | str | Биржа (`MOEX`). |
| `community` | str | Имя автора-аналитика / брокера (`Финам`, `Велес Капитал`, `АКБФ Инвестиции`). |
| `community_id` | int | Числовой ID автора (`40581` = Финам). Используется как ключ связи с [`experts`](experts.md). |
| `idea` | str | Краткий заголовок идеи (одна строка). |
| `date_in` | date (`YYYY-MM-DD`) | Дата публикации идеи (вход). |
| `date_out` | date (`YYYY-MM-DD`) | Дата, к которой автор ожидает достижения цели (плановый выход). |
| `duration_in_month` | int | Срок идеи в полных месяцах (`date_out − date_in`). |
| `price_in` | float | Цена входа, заявленная автором (в валюте `code`). |
| `price_out` | float | Целевая цена (`target price`). |
| `price_day` | float | Текущая цена бумаги на момент снимка ленты. |
| `profit_potential` | float | Потенциал к target от `price_in` в %: `(price_out − price_in) / price_in × 100`. |
| `profit_actual` | float | Фактическая доходность от `price_in` к `price_day` (%). Может быть отрицательной. |
| `stop_loss` | float / null | Стоп-лосс, если задан автором. |
| `system_status` | str | Статус идеи: `ACTIVE` (открыта), `CLOSED` (закрыта), и т.п. |
| `close_date` | date / null | Дата фактического закрытия идеи или последнего снимка цены. |
| `close_price` | float / null | Цена закрытия / последнего снимка. |
| `close_comment` | str | Комментарий автора при закрытии. |
| `close_link` | url | Ссылка на пост о закрытии. |
| `update_date` | date / null | Дата последнего апдейта (изменение target/stop). |
| `update_price` | float / null | Цена в момент апдейта. |
| `changed_at` | datetime (MSK) | Время последнего изменения записи. |

## Пример ответа (реальный, 2026-05-11, `limit=1`)
```json
[
  {
    "changed_at": "2026-05-11T01:40:02",
    "close_comment": "",
    "close_date": "2026-05-08",
    "close_link": "",
    "close_price": 505.5,
    "code": "SIBN",
    "community": "Финам",
    "community_id": 40581,
    "date_in": "2026-05-07",
    "date_out": "2026-06-07",
    "duration_in_month": 1,
    "exchange": "MOEX",
    "id": 6246,
    "idea": "Общая слабость рынка",
    "price_day": 505.5,
    "price_in": 505.0,
    "price_out": 464.15,
    "profit_actual": 0.0,
    "profit_potential": 8.0,
    "stop_loss": 521.2,
    "system_status": "ACTIVE",
    "update_date": null,
    "update_price": null
  }
]
```

## Edge cases
- `close_date`/`close_price` могут быть заполнены даже для `system_status=ACTIVE` — это снимок «на дату последнего обновления цены», а не закрытие идеи.
- Для **детального текста** идеи (`description` с тезисами, ссылка на источник, апдейты) нужен отдельный вызов [`ideas/{id}`](idea.md).
- `profit_potential` для SHORT-идей может быть положительным даже при `price_out < price_in` — FM считает потенциал по абсолютной разнице (логика автора, не математика). Проверяй знак `price_out − price_in` отдельно.
- Серверный фильтр по `community_id` отсутствует — для выборки идей одного аналитика фильтруй на клиенте.
