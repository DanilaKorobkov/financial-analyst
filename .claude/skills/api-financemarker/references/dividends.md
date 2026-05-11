---
endpoint: /api/fm/v2/dividends
---

# `dividends`

**Назначение:** календарь дивидендов **по всему рынку** (не по одной бумаге). По бумаге — раздел `dividends` в [`stocks/{exchange}:{code}`](stock.md).

**URL:** `GET /api/fm/v2/dividends`

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `mode` | str | optional | `upcoming` — только будущие выплаты. Любые другие значения (включая `last`) → HTTP 422. Отсутствие параметра — вся история. |
| `limit` / `offset` / `sort_by` / `sort_order` / `updated_in_days` | — | optional | Стандартный набор пагинации (см. SKILL.md). |

## Пример запроса
```
GET /api/fm/v2/dividends?api_token=$FINANCE_MARKER_TOKEN&limit=3
GET /api/fm/v2/dividends?api_token=$FINANCE_MARKER_TOKEN&mode=upcoming&limit=5&sort_by=last_buy_date&sort_order=asc
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` | str | Тикер бумаги, по которой выплата (`SBER`, `SBERP` — для префов отдельная запись). |
| `exchange` | str | Биржа (`MOEX`). |
| `div_amount` | float | Размер дивиденда на 1 акцию **в `div_curr`**. |
| `div_curr` | str (ISO 4217) | Валюта выплаты (`RUB`). |
| `div_percent` | float | Дивидендная доходность к цене на дату фиксации (%). Для будущих выплат считается по последней доступной цене — может отличаться от финальной цифры постфактум. |
| `last_buy_date` | date (`YYYY-MM-DD`) | Последний день покупки бумаги с правом на дивиденд (Т−2 от закрытия реестра в режиме Т+2). |
| `last_buy_price` | float / null | Цена закрытия бумаги на `last_buy_date`. Для будущих выплат `null`. |
| `reestr_close_date` | date (`YYYY-MM-DD`) | Дата закрытия реестра акционеров (отсечка). |
| `link` | url | Ссылка на первоисточник (e-disclosure / решение СД / сайт эмитента). |
| `type` | str | Тип выплаты: `Y` — годовая, `S1` / `S2` — полугодовая, `Q1`..`Q4` — квартальная, `O` — особая/разовая. |
| `year` | int | Финансовый год, по которому платится дивиденд. |
| `changed_at` | datetime (MSK) | Время последнего обновления записи. |

> В коллекционном эндпоинте без `mode=upcoming` поля `link` / `type` / `year` / `changed_at` могут отсутствовать (см. реальный пример ниже). В разделе `dividends` карточки эмитента эти поля всегда есть.

## Пример ответа (реальный, 2026-05-11, `limit=2`)
```json
[
  {
    "code": "SBERP",
    "div_amount": 37.64,
    "div_curr": "RUB",
    "div_percent": 11.7485,
    "exchange": "MOEX",
    "last_buy_date": "2026-07-17",
    "last_buy_price": null,
    "reestr_close_date": "2026-07-20"
  },
  {
    "code": "SBER",
    "div_amount": 37.64,
    "div_curr": "RUB",
    "div_percent": 11.7482,
    "exchange": "MOEX",
    "last_buy_date": "2026-07-17",
    "last_buy_price": null,
    "reestr_close_date": "2026-07-20"
  }
]
```

## Edge cases
- `mode=last` → HTTP 422 (валидно только `upcoming` или отсутствие). Для прошедших выплат не задавай `mode`.
- Доходность `div_percent` для **будущих** дивидендов считается по последней доступной цене, а не по цене на `last_buy_date` (которая ещё не наступила) — поэтому может расходиться с финальной цифрой постфактум.
- Для одной бумаги быстрее и дешевле дёрнуть `stocks/{exchange}:{code}?include=dividends` (одна страница вместо пагинации по всему рынку).
- Привилегированные и обычные акции выплачивают одинаковые дивиденды, но идут двумя отдельными записями (`SBER` и `SBERP`).
