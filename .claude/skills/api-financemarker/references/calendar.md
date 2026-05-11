---
endpoint: /api/fm/v2/calendar
---

# `calendar`

**Назначение:** календарь **предстоящих** корпоративных событий: заседания СД, ГОСА/ВОСА, отсечки, публикации отчётов и т.п.

**URL:** `GET /api/fm/v2/calendar`

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `limit` / `offset` / `sort_by` / `sort_order` | — | optional | Стандартный набор пагинации. Полезный `sort_by`: `date` (ascending — ближайшие сверху). |

> `updated_in_days` сервером **не поддерживается** для этого эндпоинта.

## Пример запроса

```http
GET /api/fm/v2/calendar?api_token=$FINANCE_MARKER_TOKEN&limit=5&sort_by=date&sort_order=asc
```

## Поля JSON-ответа

Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` | str | Тикер эмитента, у которого запланировано событие. |
| `exchange` | str | Биржа (`MOEX`). |
| `category` | str | Категория события: `BOARD_MEETING` (заседание СД), `GENERAL_MEETING` (ОСА — общее собрание акционеров), `DIVIDEND_RECORD` (отсечка), `REPORT` (плановая публикация отчёта), и т.п. |
| `event` | str | Текст события / повестка (одна строка). |
| `type` | str | Подтип события. `NA` если не уточнён. |
| `period` | str | Период отчётности для `REPORT` (`Y`, `Q`, `H1`/`H2`). `NA` для разовых событий. |
| `year` | int | Год периода (0 — не задано / разовое событие). |
| `month` | int | Месяц периода (0 — не задано). |
| `date` | date (`YYYY-MM-DD`) | **Дата запланированного события** (когда оно произойдёт). |
| `link` | url | Ссылка на e-disclosure (анонс события эмитентом). |

## Пример ответа (реальный, 2026-05-11, `limit=1`)

```json
[
  {
    "category": "BOARD_MEETING",
    "code": "DOMRF",
    "date": "2026-05-12",
    "event": "О проведении ГОСА",
    "exchange": "MOEX",
    "link": "https://www.e-disclosure.ru/portal/event.aspx?EventId=yyQRV4T7H0qzt537lFl7BA-B-B",
    "month": 0,
    "period": "NA",
    "type": "NA",
    "year": 0
  }
]
```

## Edge cases

- Уже произошедшие события сюда не попадают — для них [`disclosure`](disclosure.md).
- Ключ `event` (в `calendar`) ≠ `title` (в `disclosure`) — это разные поля; не объединяй вслепую.
- Для дивидендных отсечек, помимо `calendar`, есть специализированный [`dividends?mode=upcoming`](dividends.md) с цифрами выплаты и `last_buy_date`.
- `year`/`month=0` для разовых событий — это валидное значение, не пропуск.
