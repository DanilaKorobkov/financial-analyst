---
endpoint: /iss/statistics/engines/stock/quotedsecurities.json
block: quotedsecurities
---

# `get_quoted_securities`

**Назначение:** все котируемые бумаги с флагом `IS_QUOTED` и основным режимом `MAINBOARDID`.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/quotedsecurities>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/quotedsecurities.json
```

**Форма ответа:** блок `quotedsecurities`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле          | Тип  | Смысл                          |
| ------------- | ---- | ------------------------------ |
| `TRADEDATE`   | date | На какую дату актуально        |
| `SECID`       | str  | Тикер                          |
| `NAME`        | str  | Название                       |
| `ISIN`        | str  | ISIN                           |
| `REGNUMBER`   | str  | Регистрационный номер          |
| `MAINBOARDID` | str  | Основной режим (TQBR, RPMO, …) |
| `LISTLEVEL`   | int  | Уровень листинга (1, 2, 3)     |
| `IS_QUOTED`   | int  | 1 = бумага котируется сегодня  |

## Edge cases

- Самый быстрый способ получить «есть ли бумага в торгах сейчас» + основной режим.
