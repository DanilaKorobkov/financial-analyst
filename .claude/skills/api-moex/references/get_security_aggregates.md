---
endpoint: /iss/securities/{TICKER}/aggregates.json
block: aggregates
---

# `get_security_aggregates`

**Назначение:** агрегированные итоги торгов по бумаге за день в разрезе всех рынков (shares/repo/ndm/dark и т.д.).

**Reference ISS:** <https://iss.moex.com/iss/securities/{SECID}/aggregates>

## URL

```http
GET https://iss.moex.com/iss/securities/{TICKER}/aggregates.json
```

**Форма ответа:** блок `aggregates`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Обязательно | Дефолт  | Описание                              |
| ----------------- | ---- | ----------- | ------- | ------------------------------------- |
| `{TICKER}` (path) | str  | да          | —       | SECID бумаги                          |
| `date`            | date | нет         | сегодня | На какую торговую дату отдать агрегат |

## Поля JSON-ответа

| Поле           | Тип   | Смысл                                 |
| -------------- | ----- | ------------------------------------- |
| `market_name`  | str   | shares / repo / ndm / dark / mamc / … |
| `market_title` | str   | Человекочитаемое название рынка       |
| `engine`       | str   | Торговая система (stock)              |
| `tradedate`    | date  | Дата                                  |
| `secid`        | str   | Тикер                                 |
| `value`        | float | Дневной оборот в рублях               |
| `volume`       | int   | Дневной оборот в штуках               |
| `numtrades`    | int   | Число сделок                          |
| `updated_at`   | str   | Время последнего обновления           |

## Edge cases

- Полный дневной оборот = сумма `value` по всем строкам (один запрос).
- В выходные / праздники — пустой массив `data` в блоке.
- Для редомицилированных бумаг до даты редомициляции — пусто.
