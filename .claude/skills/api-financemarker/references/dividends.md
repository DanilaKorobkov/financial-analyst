---
endpoint: /api/fm/v2/dividends
---

# `dividends`

**Назначение:** календарь дивидендов **по всему рынку** (не по одной бумаге). По бумаге — раздел `dividends` в [`stocks/{exchange}:{code}`](stock.md).

**URL:** `GET /api/fm/v2/dividends`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)). Специфичный параметр:

| Имя | Тип | Дефолт | Описание |
|---|---|---|---|
| `mode` | str | — (= все) | `upcoming` — только будущие выплаты. Любые другие значения (включая `last`) → HTTP 422. |

## Пример запроса
```
GET /api/fm/v2/dividends?limit=3
GET /api/fm/v2/dividends?mode=upcoming&limit=5&sort_by=last_buy_date&sort_order=asc
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Тикер и биржа |
| `div_amount` | float | Размер дивиденда на 1 акцию |
| `div_curr` | str | Валюта дивиденда |
| `div_percent` | float | Доходность к цене на дату фиксации (%) |
| `last_buy_date` | date | Последний день покупки с правом на дивиденд (Т-2 от закрытия реестра) |
| `last_buy_price` | float / nullable | Цена закрытия на `last_buy_date` (для будущих выплат пусто) |
| `reestr_close_date` | date | Дата закрытия реестра |
| `link` | url | Ссылка на источник (e-disclosure / эмитент) |
| `type` | str | `Y` — годовые, `S1`/`S2` — полугодовые, `Q1..Q4` — квартальные |
| `year` | int | Финансовый год выплаты |
| `changed_at` | datetime | Время последнего обновления записи |

## Пример ответа
```json
[
  {
    "code": "PMSB",
    "exchange": "MOEX",
    "div_amount": 32.0,
    "div_curr": "RUB",
    "div_percent": 5.6547,
    "last_buy_date": "2026-05-04",
    "last_buy_price": null,
    "reestr_close_date": "2026-05-05"
  }
]
```

## Edge cases
- `mode=last` → HTTP 422 (валидно только `upcoming` или отсутствие). Для прошедших выплат не задавай `mode`.
- Доходность `div_percent` для **будущих** дивидендов считается по последней доступной цене, а не по цене на `last_buy_date` (которая ещё не наступила) — поэтому может расходиться с финальной цифрой постфактум.
- Для одной бумаги быстрее и дешевле дёрнуть `stocks/{exchange}:{code}?include=dividends` (одна страница вместо пагинации по всему рынку).
