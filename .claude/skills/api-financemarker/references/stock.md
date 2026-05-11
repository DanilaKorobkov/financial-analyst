---
endpoint: /api/fm/v2/stocks/{exchange}:{code}
---

# `stocks/{exchange}:{code}`

**Назначение:** **главный эндпоинт по эмитенту**. Возвращает агрегированный объект `Stock` с разделами: базовая карточка, сводные метрики, мультипликаторы, отчётность, дивиденды, инвест-идеи, инсайдеры, операционные показатели, акционеры, выпуски акций. Какие разделы вернуть — задаётся query-параметром `include`.

**URL:** `GET /api/fm/v2/stocks/{exchange}:{code}`

## Path-параметры
| Имя | Тип | Описание |
|---|---|---|
| `exchange` | str | Код биржи (`MOEX`) |
| `code` | str | Тикер (`SBER`, `YDEX`, `LKOH`, `OZON`) |

В URL пара передаётся через двоеточие: `/api/fm/v2/stocks/MOEX:SBER`.

## Query-параметры
| Имя | Тип | Дефолт | Описание |
|---|---|---|---|
| `include` | str (CSV) | сервер | Список разделов через запятую (см. таблицу ниже). Пример: `include=info,summary,ratios`. |

## Разделы

| Раздел | Тип | Что внутри |
|---|---|---|
| `info` | object | `StockInfo`: имя, тикер, биржа, сектор, валюта, сайт, описание |
| `summary` | object | `StockSummary`: сводные метрики (`eps`, `peg`, `growth_*`, `dividend_*`, `idea_consensus`, `insider_consensus`, `peter_lynch_target`, `graham_target` и т.п.) |
| `ratios` | array | Мультипликаторы по периодам (`pe`, `pbv`, `evebitda`, `roe`, `roic`, `debt_equity`, `gross_margin`, `net_margin`, `operation_margin` и т.п. с `year`/`period`) |
| `reports` | array | Финансовая отчётность (`revenue`, `earnings`, `ebitda`, `cfo`, `cfi`, `cff`, `total_assets`, `equity`, `total_debt`, `link` на исходный отчёт и т.п.) |
| `dividends` | array | История и прогноз дивидендов (см. также [`dividends`](dividends.md)) |
| `ideas` | array | Инвест-идеи аналитиков по бумаге |
| `insiderTransactions` | array | Сделки инсайдеров (узкий формат: insider, title, дата, тип) |
| `operations` | array | Значения операционных метрик по `operation_metric_id` (см. [`operation_metrics`](operation_metrics.md)) |
| `owners` | array | Структура акционеров (`owner`, `own` % на дату) |
| `shares` | array | Количество выпущенных акций по периодам |

## Пример запроса
```
GET /api/fm/v2/stocks/MOEX:SBER?include=info,ratios
GET /api/fm/v2/stocks/MOEX:YDEX?include=dividends
```

## Пример ответа (`include=info`, MOEX:YDEX)
```json
{
  "info": {
    "code": "YDEX",
    "name": "Яндекс",
    "exchange": "MOEX",
    "country": "Россия",
    "currency": "RUB",
    "sector": "Информационные технологии",
    "sector_id": 45,
    "industry": "Предоставление услуг в сфере информационных технологий",
    "industry_id": 451020,
    "sub_industry": "Предоставление интернет-услуг и инфраструктуры",
    "sub_industry_id": 45102030,
    "report_frequency": "Q",
    "site": "https://ir.yandex/",
    "spb": false,
    "changed_at": "2025-07-26T12:15:25"
  }
}
```

## Пример ответа (`include=dividends`, MOEX:YDEX, фрагмент)
```json
{
  "dividends": [
    {
      "code": "YDEX",
      "exchange": "MOEX",
      "div_amount": 110.0,
      "div_curr": "RUB",
      "div_percent": 2.5674,
      "last_buy_date": "2026-04-24",
      "last_buy_price": 4284.5,
      "reestr_close_date": "2026-04-27",
      "type": "Y",
      "year": 2025,
      "link": "https://www.e-disclosure.ru/portal/event.aspx?EventId=zrpJgmOFCkacqCUClBXEag-B-B",
      "changed_at": "2026-04-29T03:10:32"
    }
  ]
}
```

## Edge cases
- Раздел запрошен, но FM его не вернул (нет данных) → ключ раздела отсутствует или пустой массив. Не путать с ошибкой.
- Если бумага не существует в FM (например, новый тикер) → 400/404. Перед вызовом проверь через [`stocks`](stocks.md).
- Раздел `operations` отдаёт **значения** метрик; справочник самих метрик (что значит `aprs`, `gmv_incl_services`) — в [`operation_metrics`](operation_metrics.md).
- Раздел `insiderTransactions` здесь возвращается в коротком виде (6 полей). Полный список полей сделок — в [`insider_transactions`](insider_transactions.md).
- `reports` — большой ответ (>80 полей). Если нужны только мультипликаторы — запрашивай `include=ratios` (1 вызов = 1 единица `day_limit`, но трафик меньше).
