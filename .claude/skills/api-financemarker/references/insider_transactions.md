---
endpoint: /api/fm/v2/insider_transactions
---

# `insider_transactions`

**Назначение:** лента сделок инсайдеров **по всему рынку** (полный формат, ~24 поля). По одной бумаге в коротком формате (6 полей) — раздел `insiderTransactions` в [`stocks/{exchange}:{code}`](stock.md).

**URL:** `GET /api/fm/v2/insider_transactions`

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `limit` / `offset` / `sort_by` / `sort_order` / `updated_in_days` | — | optional | Стандартный набор пагинации. Полезные `sort_by`: `transaction_date`, `filling_date`, `value`. |

## Пример запроса
```
GET /api/fm/v2/insider_transactions?api_token=$FINANCE_MARKER_TOKEN&limit=5&updated_in_days=7
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` | str | Тикер бумаги, по которой сделка. |
| `exchange` | str | Биржа (`MOEX`). |
| `currency` | str (ISO 4217) | Валюта сделки (`RUB`). |
| `insider` | str | Имя/название инсайдера (физлицо, аффилированное юрлицо, эмитент при выкупе). `"Информация не раскрывается"` — данные скрыты по требованию ЦБ. |
| `insider_role` | str | Роль инсайдера: `OFFICER` — топ-менеджер, `DIRECTOR` — член СД, `SUBSIDIARY` — дочернее общество, `OWNER` — крупный акционер, `ISSUER` — сам эмитент (buyback), и т.п. |
| `insider_title` | str | Должность (если применимо: «Генеральный директор», «Председатель Правления»). Может быть пустой. |
| `transaction_date` | date (`YYYY-MM-DD`) | **Дата фактической сделки**. |
| `filling_date` | date (`YYYY-MM-DD`) | Дата подачи раскрытия в e-disclosure (может быть позже `transaction_date`). |
| `transaction_type` | str (1 char) | `P` — Purchase (покупка / получение), `S` — Sale (продажа / отчуждение). |
| `amount` | float | Изменение позиции **в штуках**, подписанное: `+` для покупки, `−` для продажи. |
| `shares_before` | float | Количество акций до сделки (штуки). |
| `shares_after` | float | Количество акций после сделки. |
| `own_before` | float | Доля владения до сделки (% от общего числа акций). |
| `own_after` | float | Доля после сделки (%). |
| `own_group` | str | Категория владения: `MAJOR` — мажоритарий, `MINOR` — миноритарий, и т.п. |
| `own_type` | str / null | Подтип владения (`COMMON`, `PREFERRED`, …). Может быть `null`. |
| `price` | float | Цена сделки за одну акцию в `currency`. Для безденежных передач может быть 0. |
| `value` | float | Общая стоимость сделки = `\|amount\| × price` в `currency`. Всегда положительное число. |
| `market_trade` | bool | `true` — биржевая сделка, `false` — внебиржевая (договор/наследование/REPO). |
| `approximate` | bool | `true` — данные приблизительные (FM реконструирует цифры, например для REPO). |
| `reason` | str | Текстовая причина сделки: `РЕПО ценных бумаг. Первая часть.`, `Решение Арбитражного суда…`, и т.п. |
| `link` | url | Ссылка на e-disclosure-страницу события. |
| `changed_at` | datetime (MSK) | Время обновления записи. |

## Пример ответа (реальный, 2026-05-11, `limit=1`)
```json
[
  {
    "amount": 342466.0,
    "approximate": true,
    "changed_at": "2026-05-08T09:46:52",
    "code": "GAZP",
    "currency": "RUB",
    "exchange": "MOEX",
    "filling_date": "2026-05-08",
    "insider": "Газпром",
    "insider_role": "OWNER",
    "insider_title": "",
    "link": "https://www.e-disclosure.ru/portal/event.aspx?EventId=UV402XqkQUuoa4ik4EsETA-B-B",
    "market_trade": false,
    "own_after": 0.00145,
    "own_before": 0.0,
    "own_group": "MINOR",
    "own_type": null,
    "price": 118.97,
    "reason": "Решение Арбитражного суда города Санкт-Петербурга и Ленинградской области от 13.03.2026 по делу № А56-112286/2025",
    "shares_after": 342466.0,
    "shares_before": 0.0,
    "transaction_date": "2026-05-04",
    "transaction_type": "P",
    "value": 40743200.0
  }
]
```

## Edge cases
- Пары REPO часто дают эффект «купил / продал то же количество» в ленте — фильтруй по `reason` (`РЕПО`) и `market_trade=false`, если нужны только реальные покупки/продажи.
- `amount` подписанный (`+`/`−`), `value` всегда положительный.
- `transaction_type=P` иногда приходит даже для возвратной части REPO — смотри `reason`.
- `insider="Информация не раскрывается"` — это валидное значение (не пропуск): эмитент скрыл персональные данные по правилу ЦБ. Поля `insider_role`/`shares_*`/`own_*` при этом обычно тоже обнулены.
- Для одной бумаги в коротком формате (без `price`/`value`/`shares_*`/`own_*`) используй `stocks/{exchange}:{code}?include=insiderTransactions` — там 6 полей вместо 24.
