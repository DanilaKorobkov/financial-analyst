---
endpoint: /iss/history/engines/{engine}/markets/{market}/securities.json
block: history
paginated: true
---

# `get_market_history_all`

**Назначение:** дневная история **всех** бумаг рынка за **один** день.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/securities>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/securities.json
```

**Форма ответа:** блок `history` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Обязательно | Описание                               |
| ----------------- | ---- | ----------- | -------------------------------------- |
| `date`            | date | да          | Торговая дата (без даты ISS отвергает) |
| `<block>.columns` | csv  | нет         | Подмножество колонок                   |

## Поля JSON-ответа

Аналог `get_market_history`: BOARDID, SECID, TRADEDATE, OPEN/HIGH/LOW/CLOSE, WAPRICE, VOLUME, VALUE, NUMTRADES, …

## Edge cases

- Пагинация автоматическая (~5000 строк × число режимов). Несколько секунд.
- Если `--date` приходится на нерабочий день — пустой массив `data` в блоке.
