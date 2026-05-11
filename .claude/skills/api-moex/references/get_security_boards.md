---
endpoint: /iss/securities/{TICKER}.json
block: boards
---

# `get_security_boards`

**Назначение:** все режимы торгов бумаги, **включая исторические** (с датами листинга и активности).

**Reference ISS:** блок `boards` эндпоинта <https://iss.moex.com/iss/securities/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/securities/{TICKER}.json
```

**Форма ответа:** блок `boards`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание                        |
| ----------------- | --- | ----------- | ------------------------------- |
| `{TICKER}` (path) | str | да          | SECID бумаги (например, `SBER`) |

## Поля JSON-ответа

| Поле             | Тип  | Смысл                                                 |
| ---------------- | ---- | ----------------------------------------------------- |
| `secid`          | str  | Тикер                                                 |
| `boardid`        | str  | Код режима (TQBR, EQBR, SMAL, …)                      |
| `title`          | str  | Человекочитаемое название режима                      |
| `board_group_id` | int  | ID группы режимов                                     |
| `market_id`      | int  | ID рынка                                              |
| `market`         | str  | shares / bonds / index / …                            |
| `engine_id`      | int  | ID торговой системы                                   |
| `engine`         | str  | stock / currency / futures / …                        |
| `is_traded`      | int  | 1 = режим активен сегодня                             |
| `decimals`       | int  | Знаков после запятой в цене                           |
| `history_from`   | date | Начало доступной истории                              |
| `history_till`   | date | Последняя дата истории                                |
| `listed_from`    | date | Дата начала листинга                                  |
| `listed_till`    | date | Дата делистинга (если был)                            |
| `is_primary`     | int  | 1 = основной режим (на этот режим ориентируются цены) |
| `currencyid`     | str  | Валюта (RUB, USD, …)                                  |
| `unit`           | str  | Единица (M = штука)                                   |

## Edge cases

- Несуществующий тикер → пустой массив `data` в блоке.
- Для делистнутых бумаг — `is_traded=0`, но `history_*` заполнены.
- Поле `is_primary=1` помогает выбрать «правильный» режим для расчёта доходности.
