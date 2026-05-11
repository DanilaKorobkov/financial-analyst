---
endpoint: /api/fm/v2/insider_transactions
---

# `insider_transactions`

**Назначение:** лента сделок инсайдеров **по всему рынку** (полный формат, ~24 поля). По одной бумаге в коротком формате — раздел `insiderTransactions` в [`stocks/{exchange}:{code}`](stock.md).

**URL:** `GET /api/fm/v2/insider_transactions`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)).

## Пример запроса
```
GET /api/fm/v2/insider_transactions?limit=5&updated_in_days=7
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` / `currency` | str | Тикер, биржа, валюта сделки |
| `insider` | str | Имя/название инсайдера |
| `insider_role` | str | Роль (`SUBSIDIARY`, `OFFICER`, `DIRECTOR`, …) |
| `insider_title` | str | Должность (если применимо) |
| `transaction_date` | date | Дата сделки |
| `filling_date` | date | Дата подачи раскрытия |
| `transaction_type` | str | `P` — purchase, `S` — sale |
| `amount` | float | Изменение позиции в штуках (подписанное: `+` покупка, `−` продажа) |
| `shares_before` / `shares_after` | float | Кол-во акций до/после сделки (штуки) |
| `own_before` / `own_after` | float | Доля владения до/после (%) |
| `own_group` / `own_type` | str | Категория владения (`MINOR` и т.п.) |
| `price` | float | Цена сделки за акцию |
| `value` | float | Стоимость сделки в `currency` |
| `market_trade` | bool | True — биржевая сделка, False — внебиржа/REPO |
| `approximate` | bool | True — данные приблизительные (например, REPO) |
| `reason` | str | Текстовая причина (`РЕПО ценных бумаг. Первая часть.` и т.п.) |
| `link` | url | Ссылка на e-disclosure |
| `changed_at` | datetime | Время обновления записи |

## Пример ответа
```json
[
  {
    "code": "SOFL",
    "exchange": "MOEX",
    "currency": "RUB",
    "insider": "ООО «Инвестпроекты»",
    "insider_role": "SUBSIDIARY",
    "insider_title": "",
    "transaction_date": "2026-04-29",
    "filling_date": "2026-04-30",
    "transaction_type": "P",
    "amount": 273320.0,
    "shares_before": 8712334.0,
    "shares_after": 8985654.0,
    "own_before": 2.178,
    "own_after": 2.8625,
    "own_group": "MINOR",
    "price": 67.3,
    "value": 18394400.0,
    "market_trade": false,
    "approximate": true,
    "reason": "РЕПО ценных бумаг. Вторая часть.",
    "link": "https://www.e-disclosure.ru/portal/event.aspx?EventId=YM8SlnGA1EGmY7L5f7YSmA-B-B",
    "changed_at": "2026-04-30T18:26:29"
  }
]
```

## Edge cases
- Пары REPO часто дают эффект "купил / продал то же количество" в ленте — фильтруй по `reason`/`market_trade`, если нужны только реальные покупки/продажи.
- `amount` подписанный: `+` для покупки, `−` для продажи; `value` всегда положительный.
- `transaction_type=P` иногда приходит даже для возвратной части REPO — смотри `reason`.
