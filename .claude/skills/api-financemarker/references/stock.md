---
endpoint: /api/fm/v2/stocks/{exchange}:{code}
---

# `stocks/{exchange}:{code}`

**Назначение:** **главный эндпоинт по эмитенту**. Возвращает агрегированный объект `Stock` с 10 разделами: базовая карточка, сводные метрики, мультипликаторы, отчётность, дивиденды, инвест-идеи, инсайдеры, операционные показатели, акционеры, выпуски акций. Какие разделы наполнять данными — задаётся query-параметром `include`.

**URL:** `GET /api/fm/v2/stocks/{exchange}:{code}`

## Path-параметры

| Имя | Тип | Описание |
|---|---|---|
| `exchange` | str | Код биржи (`MOEX`). |
| `code` | str | Тикер бумаги (`SBER`, `YDEX`, `LKOH`, `OZON`). |

В URL пара передаётся через двоеточие: `/api/fm/v2/stocks/MOEX:SBER`.

## Query-параметры

| Имя | Тип | Обязательность | Описание |
|---|---|---|---|
| `api_token` | str | required | API-токен (см. SKILL.md → «Авторизация»). |
| `include` | str (CSV) | optional | Список разделов через запятую. **Без `include` все массивы вернутся пустыми, а `summary` — пустой объект** — всегда указывай нужные разделы явно. Пример: `include=info,summary,ratios`. Полный набор: `info,summary,ratios,reports,dividends,ideas,insiderTransactions,operations,owners,shares`. |

## Разделы верхнего уровня

| Раздел | Тип | Что внутри (кратко) |
|---|---|---|
| `info` | object | `StockInfo` (см. [`stocks`](stocks.md)) + `description`, `site`, `disc_link`. Возвращается всегда. |
| `summary` | object | Сводные метрики эмитента «одной строкой»: EPS, дивидендная статистика, темпы роста, консенсус идей/инсайдеров, target-цены моделей Грэма и Линча. |
| `ratios` | array | Временной ряд мультипликаторов по периодам (год / квартал). Один элемент = один отчётный период. |
| `reports` | array | Финансовая отчётность по периодам (один элемент — один отчёт). ~100 полей: P&L, баланс, cash flow, ссылки на исходники. |
| `dividends` | array | История + прогноз дивидендов (формат поля см. [`dividends`](dividends.md)). |
| `ideas` | array | Инвест-идеи аналитиков по этой бумаге (формат поля см. [`ideas`](ideas.md)). |
| `insiderTransactions` | array | Сделки инсайдеров — **узкий формат** (6 полей). Полный формат — в [`insider_transactions`](insider_transactions.md). |
| `operations` | array | Значения операционных метрик (`operation_metric_id` → `value`) по периодам. Расшифровка `*_id` — в [`operation_metrics`](operation_metrics.md). |
| `owners` | array | Структура акционеров по периодам. |
| `shares` | array | Количество выпущенных акций по периодам. |

## Поля `info` (object)

| Поле | Тип | Смысл |
|---|---|---|
| Все поля из [`stocks`](stocks.md) | — | См. там. |
| `description` | str | Описание эмитента в свободной форме (несколько абзацев русского текста). |
| `site` | url | Корпоративный сайт. |
| `disc_link` | url | Ссылка на страницу раскрытия (`investor relations`). |

## Поля `summary` (object, 37 полей)

| Поле | Тип | Смысл |
|---|---|---|
| `capital` | float | Рыночная капитализация на момент снимка (в `currency` info, обычно RUB, в **миллионах**). |
| `eps` | float | Прибыль на акцию (TTM) в валюте отчётности. |
| `peg` | float | P/E ÷ ожидаемый рост прибыли. <1 → «недооценено по росту». |
| `peter_lynch_target` | float | Целевая цена по модели Питера Линча (PEG-based). |
| `graham_target` | float | Целевая цена по модели Бенджамина Грэма (формула Graham number / revised). |
| `dividend_frequency` | int | Сколько выплат в год исторически (1 — годовые, 2 — полугодовые, 4 — квартальные). |
| `dividend_strike` | int | Сколько лет подряд эмитент **не пропускал** дивиденд (дивидендный «стрик»). |
| `dividend_growth` | int | Сколько лет подряд дивиденд **рос** (без снижений). |
| `dividend_index` | float | Композитный индекс дивидендной привлекательности FM (от 0 до ~1). |
| `dividend_yield_12m` / `_3y` / `_5y` | float | Дивдоходность за последние 12 мес / средняя за 3/5 лет (%). |
| `dividend_gap_last` / `_average` | int | Размер последнего/среднего дивидендного гэпа в днях до закрытия. |
| `growth_revenue_3y` / `_5y` | float | Среднегодовой темп роста выручки (CAGR, %) за 3 и 5 лет. |
| `growth_earnings_3y` / `_5y` | float | CAGR чистой прибыли (%). |
| `growth_ebitda_3y` / `_5y` | float | CAGR EBITDA (%). Для финкомпаний (Сбер) обычно `null`. |
| `growth_assets_3y` / `_5y` | float | CAGR активов (%). |
| `growth_equity_3y` / `_5y` | float | CAGR собственного капитала (%). |
| `growth_fcf_3y` / `_5y` | float | CAGR свободного денежного потока (%). |
| `growth_netdebt_3y` / `_5y` | float | CAGR чистого долга (%). Для финкомпаний обычно `null`. |
| `growth_operation_3y` / `_5y` | float | CAGR операционной прибыли (%). |
| `idea_buy` / `idea_hold` / `idea_sell` | int | Сколько активных идей с рекомендациями BUY/HOLD/SELL по этой бумаге. |
| `idea_consensus` | str | Консенсус по идеям: `BUY`/`HOLD`/`SELL`. |
| `idea_target` | float | Усреднённая целевая цена по активным идеям. |
| `idea_potential` | float | Средний потенциал к target по активным идеям (%). |
| `insider_consensus` | str | Перекос сделок инсайдеров: `BUYS`/`SELLS`/`MIXED`. |
| `changed_at` | datetime | Время пересчёта `summary`. |

## Поля `ratios[]` (array of objects, 36 полей)

Элемент массива — мультипликаторы за один отчётный период.

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Идентификация бумаги (дублируется из `info`). |
| `year` | int | Финансовый год. |
| `month` | int | Месяц окончания периода (12 — год). |
| `period` | str | `Y` — годовой, `Q` — квартальный, `H` — полугодовой. |
| `type` | str | Стандарт отчётности: `МСФО`, `РСБУ`, …. |
| `active` | bool | `true` — последний (актуальный) период. |
| `capital` | float | Рыночная капитализация на конец периода (в млн). |
| `pe` | float | P/E (Price / Earnings TTM). |
| `pbv` | float | P/BV (Price / Book Value). |
| `ps` | float | P/S (Price / Sales). |
| `pcf` | float | P/CF (Price / Operating Cash Flow). |
| `pfcf` | float | P/FCF (Price / Free Cash Flow). |
| `pffo` | float | P/FFO — для REIT/недвижимости. |
| `evs` | float | EV/Sales. |
| `evebitda` | float | EV/EBITDA. |
| `ev_ebit` | float | EV/EBIT. |
| `debtebitda` | float | Debt / EBITDA. |
| `netdebt_ebitda` | float | NetDebt / EBITDA. |
| `debt_equity` | float | Debt / Equity (D/E). |
| `debt_ratio` | float | Total Debt / Total Assets (доля долга в активах, 0..1). |
| `current_ratio` | float | Текущая ликвидность (current assets / current liabilities). |
| `interest_coverage` | float | EBIT / Interest Expense. |
| `gross_margin` | float | Валовая маржа (%). Для финкомпаний обычно `null`. |
| `operation_margin` | float | Операционная маржа (%). |
| `ebitda_margin` | float | EBITDA-маржа (%). |
| `net_margin` | float | Чистая маржа (%). |
| `ros` | float | Return on Sales. |
| `roe` | float | Return on Equity (%). |
| `roa` | float | Return on Assets (%). |
| `roic` | float | Return on Invested Capital (%). |
| `roce` | float | Return on Capital Employed (%). |
| `dpr` | float | Dividend Payout Ratio = Dividends / Earnings. |
| `capex_revenue` | float | CAPEX / Revenue (%). |
| `net_working_capital` | float | Чистый оборотный капитал (в валюте отчётности, обычно млн). |
| `changed_at` | datetime | Время пересчёта строки. |

## Поля `reports[]` (array of objects, ~100 полей)

Один элемент — одна публикация отчёта (по периоду + стандарту). Ключевые поля сгруппированы.

| Группа → Поле | Тип | Смысл |
|---|---|---|
| **период**: `year` / `month` / `period` / `type` | int/int/str/str | См. `ratios[]`. `preliminary` (bool) = предварительный отчёт. |
| **валюта**: `curr` | str | Валюта отчётности эмитента (`RUB`). |
| **масштаб**: `amount` | int | Множитель: 1 = единицы, `1000000` = млн. Все значения ниже даны **в исходных единицах** отчёта, делить на `amount` для «как в отчёте». |
| **P&L**: `revenue` / `cost_of_sales` / `gross_profit` / `sel_gen_adm_expenses` / `operating_income` / `ebit` / `ebitda` / `ebitda_adjusted` / `earnings` / `earnings_wo_tax` / `earnings_comprehensive` / `earnings_stock_holders` / `earnings_continuing_operations` / `earnings_ps` | float | Стандартные строки отчёта о прибылях/убытках. `*_ps` — на акцию. Для финкомпаний многие поля `null`, есть `interest_income` / `interest_expense` / `interest_net` / `commission_income` / `commission_expense` / `commission_net`. |
| **Cash flow**: `cfo` / `cfi` / `cff` / `fcf` / `fcf_adjusted` / `net_change_in_cash` / `capex` / `ppe_purchase` / `intangibles_purchase` / `repurchase_of_stock` / `issuance_of_debt` / `payments_of_debt` / `net_issuance_of_debt` / `payments_for_dividends` / `cash_paid_for_interest` / `cash_paid_for_tax` | float | Кэш-флоу: операционный/инвест/финансовый, FCF, CAPEX и его компоненты, выплаты долга и дивидендов. |
| **баланс — активы**: `total_assets` / `current_assets` / `long_term_assets` / `cash_and_equiv` / `short_term_investments` / `cash_equiv_st_invesments` / `accounts_receivable` / `other_receivable` / `total_receivable` / `inventories` / `property_plant_equipment` / `ppe_rou` / `right_of_use_assets` / `intangible_assets` / `goodwill` / `goodwill_intangible_assets` / `intangible_and_tangible_assets` / `long_term_investments` | float | Активы по группам, плюс компоненты внеоборотных активов. |
| **баланс — пассивы**: `total_liabilities` / `current_liabilities` / `long_term_liabilities` / `current_debt` / `long_term_debt` / `total_debt` / `net_debt` / `net_debt_adjusted` / `cur_long_debt` / `cur_long_lease` / `current_lease` / `long_term_lease` / `accounts_payable` / `other_payable` / `total_payable` | float | Обязательства, разные срезы долга и лизинга. |
| **баланс — капитал**: `equity` / `equity_stock_holders` / `retained_earnings` / `share_premium` / `treasury_stock` | float | Собственный капитал и его компоненты. |
| **на акцию**: `earnings_ps` / `equity_ps` / `revenue_ps` / `ebitda_ps` / `ebitda_adjusted_ps` / `fcf_ps` / `fcf_adjusted_ps` | float | Per-share метрики. |
| **прочее**: `depr_depl_amort` / `research_and_development` / `total_expenses` / `other_operating_income` / `ffo` | float | Амортизация, R&D, FFO (для REIT). |
| **ссылки**: `link` / `link_press` / `link_update` | url | Исходный PDF/Excel отчёта, пресс-релиз, ссылка на обновление. |
| **служебное**: `changed_at` (datetime) | — | Время записи. |

> Полный список реальных ключей (на примере SBER 2011, МСФО, Y): `accounts_payable, accounts_receivable, amount, capex, cash_and_equiv, cash_equiv_st_invesments, cash_paid_for_interest, cash_paid_for_tax, cff, cfi, cfo, changed_at, code, commission_expense, commission_income, commission_net, cost_of_sales, cur_long_debt, cur_long_lease, curr, current_assets, current_debt, current_lease, current_liabilities, depr_depl_amort, earnings, earnings_comprehensive, earnings_comprehensive_stock_holders, earnings_continuing_operations, earnings_ps, earnings_stock_holders, earnings_wo_tax, ebit, ebitda, ebitda_adjusted, ebitda_adjusted_ps, ebitda_ps, equity, equity_ps, equity_stock_holders, exchange, fcf, fcf_adjusted, fcf_adjusted_ps, fcf_ps, ffo, goodwill, goodwill_intangible_assets, gross_profit, intangible_and_tangible_assets, intangible_assets, intangibles_purchase, interest_expense, interest_income, interest_net, inventories, issuance_of_debt, link, link_press, link_update, long_term_assets, long_term_debt, long_term_investments, long_term_lease, long_term_liabilities, month, net_change_in_cash, net_debt, net_debt_adjusted, net_issuance_of_debt, operating_income, other_operating_income, other_payable, other_receivable, payments_for_dividends, payments_of_debt, period, ppe_purchase, ppe_rou, preliminary, property_plant_equipment, repurchase_of_stock, research_and_development, retained_earnings, revenue, revenue_ps, right_of_use_assets, sel_gen_adm_expenses, share_premium, short_term_investments, total_assets, total_debt, total_expenses, total_liabilities, total_payable, total_receivable, treasury_stock, type, year`.

## Поля `dividends[]`, `ideas[]`, `insiderTransactions[]`

См. отдельные reference:

- `dividends[]` — формат как у [`dividends`](dividends.md) (для одной бумаги поле `link` и `type`/`year` могут быть пустыми).
- `ideas[]` — формат как у [`ideas`](ideas.md).
- `insiderTransactions[]` — **узкий формат**: только `code`, `exchange`, `insider`, `insider_title`, `transaction_date`, `transaction_type`. Полный набор полей (price, value, shares_before/after и т.д.) — в [`insider_transactions`](insider_transactions.md).

## Поля `operations[]`

Один элемент — значение одной операционной метрики за один период.

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Идентификация бумаги. |
| `operation_metric_id` | str | ID метрики (`car_loans`, `gmv_incl_services`, …). Расшифровка в [`operation_metrics`](operation_metrics.md). |
| `year` / `month` / `period` | int/int/str | Период (как в `ratios[]`). |
| `value` | float | Значение метрики в единицах `unit` (уже умножено на множитель `amount` из словаря). |
| `unit` | str | Единица итогового `value` (`₽`, `шт`, `чел`, `t`, `m2`, `%`). |
| `amount` | int | Множитель из словаря (`1` или `1000000`). |
| `original_value` / `original_unit` / `original_amount` | float / str / int | Значение «как опубликовал эмитент» с его единицей (например, отчёт в `млн ₽`). |
| `curs` | float | Курс пересчёта (если метрика в нерублёвой валюте). Для рублёвых обычно `1.0`. |
| `link` | url | Ссылка на исходный отчёт эмитента. |
| `link_update` | url | Ссылка на обновление (если был). |

## Поля `owners[]`

Один элемент — доля одного акционера на дату.

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Идентификация бумаги. |
| `year` / `month` | int | Дата среза. |
| `owner` | str | Имя/название акционера или категория (`Прочие`, `Free float`). |
| `own` | str (decimal) | Доля владения в % (строка с десятичным числом). |
| `link` | url | Ссылка на источник раскрытия структуры. |
| `changed_at` | datetime | Время обновления записи. |

## Поля `shares[]`

Один элемент — количество выпущенных акций на дату.

| Поле | Тип | Смысл |
|---|---|---|
| `code` / `exchange` | str | Идентификация бумаги (включая `SBERP` отдельно от `SBER`). |
| `year` / `month` | int | Дата среза. |
| `num` | int | Количество акций (в штуках). |

## Пример запроса

```http
GET /api/fm/v2/stocks/MOEX:SBER?api_token=$FINANCE_MARKER_TOKEN&include=info,summary,ratios
GET /api/fm/v2/stocks/MOEX:YDEX?api_token=$FINANCE_MARKER_TOKEN&include=dividends
```

## Пример ответа (`include=summary`, MOEX:SBER, реальный 2026-05-11)

```json
{
  "summary": {
    "capital": 97627.3,
    "changed_at": "2026-05-11T03:32:06",
    "dividend_frequency": 1,
    "dividend_gap_average": 482,
    "dividend_gap_last": 280,
    "dividend_growth": 3,
    "dividend_index": 0.7,
    "dividend_strike": 4,
    "dividend_yield_12m": 10.65,
    "dividend_yield_3y": 10.55,
    "dividend_yield_5y": 7.5,
    "eps": 78.8,
    "graham_target": 160.46,
    "growth_assets_3y": 9.69,
    "growth_assets_5y": 10.89,
    "growth_earnings_3y": 5.62,
    "growth_earnings_5y": 7.37,
    "growth_equity_3y": 10.45,
    "growth_equity_5y": 9.72,
    "growth_fcf_5y": 7.57,
    "growth_revenue_3y": 10.59,
    "growth_revenue_5y": 13.42,
    "idea_buy": 9,
    "idea_consensus": "BUY",
    "idea_hold": 3,
    "idea_potential": 20.9125,
    "idea_sell": 0,
    "idea_target": 387.392,
    "insider_consensus": "BUYS",
    "peg": 0.56,
    "peter_lynch_target": 96.77
  }
}
```

## Пример ответа (`include=ratios`, MOEX:SBER, фрагмент — годовой 2011 МСФО)

```json
{
  "ratios": [
    {
      "active": false,
      "capex_revenue": 14.74,
      "capital": 54822.8,
      "code": "SBER",
      "debt_equity": 7.54,
      "debt_ratio": 0.88,
      "exchange": "MOEX",
      "month": 12,
      "net_margin": 42.58,
      "pbv": 1.4,
      "pcf": -7.23,
      "pe": 5.58,
      "period": "Y",
      "pfcf": -4.99,
      "ps": 2.38,
      "roa": 2.92,
      "roe": 24.91,
      "type": "МСФО",
      "year": 2011
    }
  ]
}
```

## Пример ответа (`include=owners`, MOEX:SBER, реальный)

```json
{
  "owners": [
    {
      "code": "SBER",
      "exchange": "MOEX",
      "month": 6,
      "year": 2022,
      "owner": "Прочие",
      "own": "50.0",
      "link": "https://www.sberbank.com/ru/investor-relations/share-profile",
      "changed_at": "2024-11-14T03:31:31"
    }
  ]
}
```

## Edge cases

- **`include` обязателен по смыслу**: без него `summary` — пустой объект, а все массивы — пустые. `info` возвращается всегда.
- Раздел запрошен, но FM его не вернул (нет данных для бумаги) → ключ присутствует, но массив пустой / `summary` пустой. Это **не ошибка** API.
- Если бумага не существует в FM (например, новый тикер) → 400/404. Перед вызовом проверь через [`stocks`](stocks.md).
- `reports` — большой ответ (для SBER ~92 элемента, JSON-ответ ~360 КБ при `include=summary,ratios,reports,dividends,ideas,insiderTransactions,operations,owners,shares`). Запрашивай только нужное.
- `operations` тоже объёмный (для SBER — 267 элементов). Если нужен один KPI, фильтруй на клиенте по `operation_metric_id`.
- Для банков (`SBER`, `VTBR`) поля `ebitda`/`gross_margin`/`current_ratio` в `ratios` будут `null` — это специфика отчётности кредитных организаций, не баг.
- Раздел `insiderTransactions` здесь возвращается в **коротком виде** (6 полей). Полный формат — в [`insider_transactions`](insider_transactions.md).
- 1 вызов = 1 единица `day_limit` независимо от размера `include`. Объединяй разделы в один запрос, если они нужны вместе.
