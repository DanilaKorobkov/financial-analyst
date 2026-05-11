---
endpoint: /iss/history/engines/stock/markets/bonds/boards/{board}/yields/{TICKER}.json
block: history_yields
paginated: true
---

# `get_bond_yields`

**Назначение:** история доходностей одной облигации в режиме (с пагинацией).

**Reference ISS:** <https://iss.moex.com/iss/history/engines/stock/markets/bonds/boards/{BOARD}/yields/{SECID}>

## URL

```http
GET https://iss.moex.com/iss/history/engines/stock/markets/bonds/boards/{board}/yields/{TICKER}.json
```

**Форма ответа:** блок `history_yields` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Обязательно | Дефолт | Описание        |
| ----------------- | ---- | ----------- | ------ | --------------- |
| `{TICKER}` (path) | str  | да          | —      | SECID облигации |
| `{board}` (path)  | str  | нет         | `TQOB` | Режим           |
| `from`            | date | нет         | —      | Начало периода  |
| `till`            | date | нет         | —      | Конец периода   |

## Поля JSON-ответа (блок `history_yields`)

| Поле                  | Тип   | Смысл                                        |
| --------------------- | ----- | -------------------------------------------- |
| `TRADEDATE`           | date  | Дата                                         |
| `SECID/BOARDID`       | str   | Тикер/режим                                  |
| `YIELDDATE`           | date  | Дата погашения / оферты                      |
| `YIELDDATETYPE`       | str   | `MATDATE` / `OFFERDATE`                      |
| `PRICE/WAPRICE`       | float | Цена / WAP                                   |
| `ACCINT`              | float | Накопленный купонный доход                   |
| `EFFECTIVEYIELD`      | float | Эффективная доходность к погашению/оферте, % |
| `DURATION`            | int   | Дюрация в днях                               |
| `ZSPREADBP/GSPREADBP` | int   | Z-спред / G-спред в б.п.                     |
| `YIELDTOOFFER`        | float | Доходность к оферте                          |
| `YIELDLASTCOUPON`     | float | Доходность к последнему купону               |

## Edge cases

- ОФЗ обычно торгуются в режиме TQOB.
