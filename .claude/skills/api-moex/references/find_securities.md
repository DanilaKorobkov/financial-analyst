---
endpoint: /iss/securities.json
block: securities
---

# `find_securities`

**Назначение:** найти инструменты по подстроке кода / названия / ISIN / идентификатора эмитента / номера гос.регистрации. Типовое применение — по рег.номеру узнать предыдущие тикеры эмитента, чтобы собрать длинную историю при ребрендинге.

**Reference ISS:** <https://iss.moex.com/iss/reference/5>

## URL

```http
GET https://iss.moex.com/iss/securities.json
```

**Форма ответа:** блок `securities`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт | Описание                                           |
| ----------------- | --- | ----------- | ------ | -------------------------------------------------- |
| `QUERY`           | str | да          | —      | Подстрока: тикер / ISIN / рег.номер / имя эмитента |
| `<block>.columns` | csv | нет         | все    | Подмножество колонок через запятую                 |

## Пример запроса

```bash

```

## Поля JSON-ответа

Массив объектов. Полный набор полей:

| Поле                  | Тип      | Обязательно | Смысл                                                                                 |
| --------------------- | -------- | ----------- | ------------------------------------------------------------------------------------- |
| `secid`               | str      | да          | Тикер инструмента                                                                     |
| `shortname`           | str      | да          | Краткое название                                                                      |
| `regnumber`           | str/null | нет         | Номер гос.регистрации (null для деривативов)                                          |
| `name`                | str      | да          | Полное название инструмента                                                           |
| `isin`                | str/null | нет         | ISIN (латиницей; null для деривативов)                                                |
| `is_traded`           | int      | да          | 1 = торгуется, 0 = нет                                                                |
| `emitent_id`          | int/null | нет         | Внутренний ID эмитента в MOEX                                                         |
| `emitent_title`       | str/null | нет         | Полное наименование эмитента                                                          |
| `emitent_inn`         | str/null | нет         | ИНН эмитента                                                                          |
| `emitent_okpo`        | str/null | нет         | ОКПО эмитента                                                                         |
| `type`                | str      | да          | `common_share`, `preferred_share`, `corporate_bond`, `futures`, `option_on_shares`, … |
| `group`               | str      | да          | Группа инструмента (`stock_shares`, `stock_bonds`, `stock_index`, …)                  |
| `primary_boardid`     | str/null | нет         | Основной режим торгов (`TQBR`, `TQTF`, …)                                             |
| `marketprice_boardid` | str/null | нет         | Режим, по которому считается рыночная цена                                            |

## Пример ответа

```json
[
  {
    "secid": "YDEX",
    "shortname": "ЯНДЕКС",
    "regnumber": "1-01-16777-A",
    "name": "МКПАО ЯНДЕКС ао",
    "isin": "RU000A107T19",
    "is_traded": 1,
    "emitent_id": 15523,
    "emitent_title": "\"Международная компания публичное акционерное общество \"\"ЯНДЕКС\"\"\"",
    "emitent_inn": 3900019850,
    "emitent_okpo": "...",
    "type": "common_share",
    "group": "stock_shares",
    "primary_boardid": "TQBR",
    "marketprice_boardid": "TQBR"
  }
]
```

## Edge cases

- Поиск по подстроке — на короткий `QUERY` (`SBER`, `YD`) возвращает много шума: фьючерсы, опционы, ИОС с `regnumber=null`. Уточняй через ISIN/рег.номер.
- Поле `emitent_title` — единственный источник имени эмитента (в `find_security_description` его нет).
- Пустой массив `[]` → ничего не найдено.
