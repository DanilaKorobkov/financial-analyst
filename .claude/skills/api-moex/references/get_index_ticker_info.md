---
endpoint: /iss/statistics/engines/stock/markets/index/analytics/{index}/tickers/{TICKER}.json
block: ticker
paginated: true
---

# `get_index_ticker_info`

**Назначение:** аналитика по конкретному тикеру в индексе (история включений с весами и фактами).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics/{INDEX}/tickers/{TICKER}>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics/{index}/tickers/{TICKER}.json
```

**Форма ответа:** блок `ticker` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание        |
| ----------------- | --- | ----------- | --------------- |
| `INDEX`           | str | да          | Код индекса     |
| `{TICKER}` (path) | str | да          | Тикер в индексе |

## Поля JSON-ответа

Поля блока `ticker` (с пагинацией): дата, вес, факторы.

## Edge cases

- Если бумаги нет в индексе — пустой массив `data` в блоке.
