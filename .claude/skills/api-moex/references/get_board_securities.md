---
endpoint: /iss/engines/{engine}/markets/{market}/boards/{board}/securities.json
block: securities
---

# `get_board_securities`

**Назначение:** получить таблицу **всех бумаг указанного режима торгов** со справочными данными.

**Reference ISS:** <https://iss.moex.com/iss/reference/32>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boards/{board}/securities.json
```

**Форма ответа:** блок `securities`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт   | Описание             |
| ----------------- | --- | ----------- | -------- | -------------------- |
| `{board}` (path)  | str | нет         | `TQBR`   | Режим торгов         |
| `{market}` (path) | str | нет         | `shares` | Рынок                |
| `{engine}` (path) | str | нет         | `stock`  | Движок               |
| `<block>.columns` | csv | нет         | все      | Подмножество колонок |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов. Полный набор полей справочника:

| Поле                  | Тип        | Обязательно | Смысл                         |
| --------------------- | ---------- | ----------- | ----------------------------- |
| `SECID`               | str        | да          | Тикер                         |
| `BOARDID`             | str        | да          | Код режима                    |
| `SHORTNAME`           | str        | да          | Краткое название              |
| `PREVPRICE`           | float      | нет         | Цена закрытия пред. сессии    |
| `LOTSIZE`             | int        | да          | Размер лота                   |
| `FACEVALUE`           | float      | нет         | Номинал                       |
| `STATUS`              | str        | нет         | Статус                        |
| `BOARDNAME`           | str        | да          | Название режима               |
| `DECIMALS`            | int        | да          | Знаков после запятой          |
| `SECNAME`             | str        | нет         | Короткое название эмитента    |
| `REMARKS`             | str/null   | нет         | Примечания                    |
| `MARKETCODE`          | str        | нет         | Код рынка                     |
| `INSTRID`             | str        | нет         | ID инструмента                |
| `SECTORID`            | str/null   | нет         | ID сектора                    |
| `MINSTEP`             | float      | да          | Минимальный шаг цены          |
| `PREVWAPRICE`         | float/null | нет         | Средневзвешенная пред. дня    |
| `FACEUNIT`            | str        | нет         | Валюта номинала (`SUR` = ₽)   |
| `PREVDATE`            | date       | нет         | Дата пред. торговой сессии    |
| `ISSUESIZE`           | int        | нет         | Объём выпуска                 |
| `ISIN`                | str        | нет         | ISIN                          |
| `LATNAME`             | str        | нет         | Латинское название            |
| `REGNUMBER`           | str        | нет         | Гос.рег.номер                 |
| `PREVLEGALCLOSEPRICE` | float/null | нет         | Юр. цена закрытия пред. дня   |
| `CURRENCYID`          | str        | нет         | Валюта инструмента            |
| `SECTYPE`             | str        | нет         | Тип бумаги                    |
| `LISTLEVEL`           | int        | нет         | Котировальный уровень (1/2/3) |
| `SETTLEDATE`          | date       | нет         | Дата расчётов                 |

## Пример ответа

```json
[
  {
    "SECID": "ABIO",
    "SHORTNAME": "iАРТГЕН ао",
    "LOTSIZE": 10,
    "PREVPRICE": 57.64,
    "BOARDID": "TQBR"
  },
  {
    "SECID": "SBER",
    "SHORTNAME": "Сбербанк",
    "LOTSIZE": 1,
    "PREVPRICE": 319.81,
    "BOARDID": "TQBR"
  },
  {
    "SECID": "YDEX",
    "SHORTNAME": "ЯНДЕКС",
    "LOTSIZE": 1,
    "PREVPRICE": 4060,
    "BOARDID": "TQBR"
  }
]
```

## Edge cases

- Текущие сессионные данные (`LAST`, `BID`, `OFFER`, `OPEN`, `LOW`, `HIGH`, `VOLTODAY`) приходят в **отдельном блоке `marketdata`** того же эндпоинта; чтобы их получить, не сужай `iss.only` до одного блока `securities` — см. [`get_market_security`](get_market_security.md).
- Дефолтный `{board}=TQBR` = акции T+2; для облигаций — `{board}=TQOB`, `{market}=bonds`; для ETF — `{board}=TQTF`.
