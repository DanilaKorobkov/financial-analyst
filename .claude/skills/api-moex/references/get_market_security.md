---
endpoint: /iss/engines/{engine}/markets/{market}/securities/{TICKER}.json
---

# `get_market_security`

**Назначение:** все режимы торгов конкретной бумаги с marketdata одним вызовом (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{secid}>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/securities/{TICKER}.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание     |
| ----------------- | --- | ----------- | ------------ |
| `{TICKER}` (path) | str | да          | SECID бумаги |

## Поля JSON-ответа

Объединение блоков `securities` + `marketdata` + `marketdata_yields` (для облигаций) + `dataversion` через колонку `_block`.

## Edge cases

- Альтернатива `marketdata <TICKER>` — та возвращает один режим, эта все.
