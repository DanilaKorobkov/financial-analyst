---
endpoint: /iss/history/engines/currency/markets/selt/securities/{TICKER}.json
block: history
paginated: true
---

# `get_currency_history`

**Назначение:** дневная история торгов валютной пары на рынке SELT.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/currency/markets/selt/securities/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/history/engines/currency/markets/selt/securities/{TICKER}.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Обязательно | Описание             |
| ----------------- | ---- | ----------- | -------------------- |
| `{TICKER}` (path) | str  | да          | SECID валютной пары  |
| `from`            | date | нет         | Начало диапазона     |
| `till`            | date | нет         | Конец диапазона      |
| `<block>.columns` | csv  | нет         | Подмножество колонок |

## Поля JSON-ответа

| Поле                  | Тип   | Смысл                          |
| --------------------- | ----- | ------------------------------ |
| `BOARDID`             | str   | Режим (CETS, AUCB, CNGD, LICU) |
| `TRADEDATE`           | date  | Дата                           |
| `SECID/SHORTNAME`     | str   | Тикер/имя                      |
| `OPEN/LOW/HIGH/CLOSE` | float | OHLC                           |
| `NUMTRADES`           | int   | Сделок                         |
| `WAPRICE`             | float | Взвешенная цена                |

## Edge cases

- На одну дату приходится несколько строк (по числу режимов).
- Для `USD000UTSTOM` история с 2003 года.
