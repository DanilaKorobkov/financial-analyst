---
endpoint: /api/fm/v2/stocks
---

# `stocks`

**Назначение:** список компаний с базовой карточкой `StockInfo` (тикер, название, биржа, страна, валюта, сектор/отрасль). Поддерживает пагинацию и фильтр по обновлению.

**URL:** `GET /api/fm/v2/stocks`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md) → «Query-параметры пагинации»). Серверной фильтрации по сектору/стране нет — фильтруй на клиенте.

## Пример запроса
```
GET /api/fm/v2/stocks?limit=3
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `code` | str | Тикер |
| `name` | str | Название эмитента |
| `exchange` | str | Код биржи |
| `country` | str | Страна |
| `currency` | str | Валюта торгов |
| `sector` / `sector_id` | str / int | Сектор и его код |
| `industry_group` / `industry_group_id` | str / int | Группа отраслей |
| `industry` / `industry_id` | str / int | Отрасль |
| `sub_industry` / `sub_industry_id` | str / int | Под-отрасль |
| `primary_report_code` / `primary_report_exchange` | str | Первичный тикер/биржа (для дублей типа SBER/SBERP) |
| `report_frequency` | str | `Q`/`Y` — частота публикации отчётов |
| `spb` | bool | Торгуется ли на СПБ-бирже |
| `changed_at` | datetime | Время последнего обновления карточки |

## Пример ответа
```json
[
  {
    "code": "SBER",
    "name": "Сбербанк",
    "exchange": "MOEX",
    "country": "Россия",
    "currency": "RUB",
    "sector": "Финансы",
    "sector_id": 40,
    "industry_group": "Банковская деятельность",
    "industry_group_id": 4010,
    "industry": "Банковская деятельность",
    "industry_id": 401010,
    "sub_industry": "Диверсифицированные банки",
    "sub_industry_id": 40101010,
    "primary_report_code": "SBER",
    "primary_report_exchange": "MOEX",
    "report_frequency": "Q",
    "spb": false,
    "changed_at": "2026-05-04T03:32:11"
  }
]
```

## Edge cases
- Привилегированные акции (SBERP) и обычные (SBER) — две отдельные записи; `primary_report_code` указывает на ту, к которой привязаны отчёты.
- Для разовой карточки одного эмитента не вызывай этот эндпоинт — иди в [`stocks/{exchange}:{code}`](stock.md) со специфичным идентификатором.
