---
endpoint: /iss/statistics/engines/stock/markets/index/analytics/{index}.json
block: analytics
paginated: true
---

# `get_index_analytics`

**Назначение:** **состав индекса с весами и факторами** на текущую или указанную дату.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics/{INDEXID}>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics/{index}.json
```

**Форма ответа:** блок `analytics` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя     | Тип  | Обязательно | Описание                     |
| ------- | ---- | ----------- | ---------------------------- |
| `INDEX` | str  | да          | Код индекса (IMOEX, RTSI, …) |
| `date`  | date | нет         | Состав на указанную дату     |

## Поля JSON-ответа

| Поле                 | Тип   | Смысл                  |
| -------------------- | ----- | ---------------------- |
| `indexid`            | str   | Код индекса            |
| `tradedate`          | date  | Дата                   |
| `ticker`             | str   | Тикер бумаги в составе |
| `shortnames`         | str   | Короткие имена         |
| `secids`             | str   | Алиасы тикеров         |
| `weight`             | float | **Вес в индексе, %**   |
| `tradingsession`     | int   | Сессия                 |
| `trade_session_date` | date  | Дата сессии            |

(Полная схема включает waprice, issue_size_total, cap_total, ff_factor, w_factor, issue_size_index, cap_index, value, num_trades, volatility, factora, factorb, influence, determinat — но они приходят только при детальной выгрузке.)

## Edge cases

- В отличие от `get_index_tickers` — отдаёт **фактические веса**.
- Пагинация автоматическая (шаг 20).
