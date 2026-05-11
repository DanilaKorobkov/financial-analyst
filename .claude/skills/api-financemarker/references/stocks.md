---
endpoint: /api/fm/v2/stocks
---

# `stocks`

**Назначение:** список компаний с базовой карточкой `StockInfo` (тикер, название, биржа, страна, валюта, сектор/отрасль). Поддерживает пагинацию и фильтр по дате обновления.

**URL:** `GET /api/fm/v2/stocks`

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `limit` / `offset` / `sort_by` / `sort_order` / `updated_in_days` | — | optional | Стандартный набор пагинации (см. SKILL.md → «Query-параметры пагинации»). |

Серверной фильтрации по сектору/стране/отрасли нет — фильтруй на клиенте по полям `sector`, `industry`, `country`.

## Пример запроса

```http
GET /api/fm/v2/stocks?api_token=$FINANCE_MARKER_TOKEN&limit=3
GET /api/fm/v2/stocks?api_token=$FINANCE_MARKER_TOKEN&updated_in_days=7&limit=100
```

## Поля JSON-ответа

Массив объектов `StockInfo`:

| Поле | Тип | Смысл |
|---|---|---|
| `code` | str | Тикер бумаги на бирже (`SBER`, `SBERP`, `YDEX`). Уникален в паре `(exchange, code)`. |
| `name` | str | Название эмитента / тип бумаги (`Сбербанк`, `Сбербанк привилегированные`). |
| `exchange` | str | Код биржи листинга (`MOEX`). |
| `country` | str | Страна регистрации эмитента (русское название, например `Россия`). |
| `currency` | str | Валюта торгов бумагой (ISO 4217, `RUB`). |
| `sector` | str | GICS-сектор эмитента на русском (`Финансы`, `Энергетика`). |
| `sector_id` | int | Числовой код GICS-сектора (`40` = Финансы). |
| `industry_group` | str | Группа отраслей GICS на русском (`Банковская деятельность`). |
| `industry_group_id` | int | Числовой код группы отраслей (`4010`). |
| `industry` | str | Отрасль GICS на русском. |
| `industry_id` | int | Числовой код отрасли (`401010`). |
| `sub_industry` | str | Под-отрасль GICS на русском (`Диверсифицированные банки`). |
| `sub_industry_id` | int | Числовой код под-отрасли (`40101010`). |
| `primary_report_code` | str | Тикер «основной» бумаги эмитента, к которой привязаны отчёты (для SBERP это `SBER`). |
| `primary_report_exchange` | str | Биржа этой основной бумаги. |
| `report_frequency` | str | Частота публикации отчётности эмитентом: `Q` — квартальная, `Y` — годовая. |
| `spb` | bool | `true` — бумага дополнительно листингована на СПБ-бирже (обычно `false` для нашей подписки MOEX-only). |
| `changed_at` | datetime (ISO 8601, MSK) | Время последнего изменения карточки на стороне FM. |

## Пример ответа (реальный, 2026-05-11, `limit=2`)

```json
[
  {
    "changed_at": "2026-05-11T03:32:06",
    "code": "SBERP",
    "country": "Россия",
    "currency": "RUB",
    "exchange": "MOEX",
    "industry": "Банковская деятельность",
    "industry_group": "Банковская деятельность",
    "industry_group_id": 4010,
    "industry_id": 401010,
    "name": "Сбербанк привилегированные",
    "primary_report_code": "SBER",
    "primary_report_exchange": "MOEX",
    "report_frequency": "Q",
    "sector": "Финансы",
    "sector_id": 40,
    "spb": false,
    "sub_industry": "Диверсифицированные банки",
    "sub_industry_id": 40101010
  },
  {
    "changed_at": "2026-05-11T03:32:06",
    "code": "SBER",
    "country": "Россия",
    "currency": "RUB",
    "exchange": "MOEX",
    "industry": "Банковская деятельность",
    "industry_group": "Банковская деятельность",
    "industry_group_id": 4010,
    "industry_id": 401010,
    "name": "Сбербанк",
    "primary_report_code": "SBER",
    "primary_report_exchange": "MOEX",
    "report_frequency": "Q",
    "sector": "Финансы",
    "sector_id": 40,
    "spb": false,
    "sub_industry": "Диверсифицированные банки",
    "sub_industry_id": 40101010
  }
]
```

## Edge cases

- Привилегированные акции (`SBERP`) и обычные (`SBER`) — две отдельные записи; `primary_report_code` указывает на ту, к которой привязаны отчёты.
- Для разовой карточки одного эмитента не вызывай этот эндпоинт — иди в [`stocks/{exchange}:{code}`](stock.md) со специфичным идентификатором.
- `industry` и `industry_group` могут совпадать по тексту (как у Сбера), но их `*_id` всегда разной размерности (`401010` vs `4010`).
