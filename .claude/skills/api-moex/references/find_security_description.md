---
endpoint: /iss/securities/{TICKER}.json
block: description
---

# `find_security_description`

**Назначение:** получить **полную спецификацию инструмента** — все атрибуты эмиссии (даты, размер выпуска, листинг, флаги сессий).

**Reference ISS:** <https://iss.moex.com/iss/reference/13>

## URL

```http
GET https://iss.moex.com/iss/securities/{TICKER}.json
```

**Форма ответа:** блок `description`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Дефолт | Описание     |
| ----------------- | --- | ----------- | ------ | ------------ |
| `{TICKER}` (path) | str | да          | —      | Тикер бумаги |

## Пример запроса

```http
GET https://iss.moex.com/iss/securities/{TICKER}.json?iss.json=extended&iss.meta=off
```

## Поля JSON-ответа

ISS отдаёт блок `description` как массив строк вида `{name, title, value, type, ...}`. В таблице ниже перечислены ключи `name`, которые встречаются для эмитента; значение читается из поля `value`. Все значения приходят **строками** — приводи типы сам.

| Ключ                   | Тип              | Обязательно | Смысл                                                                 |
| ---------------------- | ---------------- | ----------- | --------------------------------------------------------------------- |
| `SECID`                | str              | да          | Тикер                                                                 |
| `ISSUENAME`            | str              | да          | Название выпуска (`Акции обыкновенные`, …)                            |
| `NAME`                 | str              | да          | Полное название бумаги                                                |
| `SHORTNAME`            | str              | да          | Краткое название                                                      |
| `ISIN`                 | str              | нет         | ISIN (всегда латиницей); может отсутствовать у некоторых инструментов |
| `REGNUMBER`            | str              | нет         | Номер гос.регистрации                                                 |
| `ISSUESIZE`            | str(int)         | да          | Объём выпуска, шт                                                     |
| `FACEVALUE`            | str(float)       | нет         | Номинал                                                               |
| `FACEUNIT`             | str              | нет         | Валюта номинала (`SUR` = ₽)                                           |
| `ISSUEDATE`            | str(date)        | да          | Дата начала торгов на MOEX                                            |
| `LATNAME`              | str              | нет         | Латинское название                                                    |
| `HASPROSPECTUS`        | str("0"/"1")     | нет         | Есть ли проспект эмиссии                                              |
| `HASDEFAULT`           | str("0"/"1")     | нет         | Был ли дефолт                                                         |
| `HASTECHNICALDEFAULT`  | str("0"/"1")     | нет         | Тех. дефолт                                                           |
| `EMITENTMISMATCHCUR`   | str("0"/"1")     | нет         | Несовпадение валют                                                    |
| `LISTLEVEL`            | str("1"/"2"/"3") | нет         | Котировальный уровень                                                 |
| `ISQUALIFIEDINVESTORS` | str("0"/"1")     | нет         | `"0"` = доступна всем, `"1"` = только квалам                          |
| `MORNINGSESSION`       | str("0"/"1")     | нет         | Утренняя сессия                                                       |
| `EVENINGSESSION`       | str("0"/"1")     | нет         | Вечерняя сессия                                                       |
| `WEEKENDSESSION`       | str("0"/"1")     | нет         | Сессия выходного дня                                                  |
| `REGISTRY_DATE`        | str(date)        | нет         | Дата регистрации выпуска                                              |
| `TYPENAME`             | str              | да          | Человекочитаемый тип (`Акция обыкновенная`)                           |
| `GROUP`                | str              | да          | Группа (`stock_shares`, …)                                            |
| `TYPE`                 | str              | да          | Тип (`common_share`, …)                                               |
| `GROUPNAME`            | str              | да          | Человекочитаемая группа (`Акции`)                                     |
| `EMITTER_ID`           | str(int)         | нет         | ID эмитента в MOEX                                                    |

## Пример ответа

```json
{
  "SECID": "YDEX",
  "ISSUENAME": "Акции обыкновенные",
  "NAME": "МКПАО ЯНДЕКС",
  "SHORTNAME": "ЯНДЕКС",
  "ISIN": "RU000A107T19",
  "REGNUMBER": "1-01-16777-A",
  "ISSUESIZE": 394967916,
  "FACEVALUE": 0.4,
  "FACEUNIT": "SUR",
  "ISSUEDATE": "2024-07-08",
  "LATNAME": "YANDEX",
  "HASPROSPECTUS": 0,
  "HASDEFAULT": 0,
  "HASTECHNICALDEFAULT": 0,
  "EMITENTMISMATCHCUR": 0,
  "LISTLEVEL": 1,
  "ISQUALIFIEDINVESTORS": 0,
  "MORNINGSESSION": 1,
  "EVENINGSESSION": 1,
  "WEEKENDSESSION": 1,
  "REGISTRY_DATE": "2023-12-25",
  "TYPENAME": "Акция обыкновенная",
  "GROUP": "stock_shares",
  "TYPE": "common_share",
  "GROUPNAME": "Акции",
  "EMITTER_ID": 15523
}
```

## Edge cases

- **`EMITENT_TITLE` / `EMITENT_INN` отсутствуют** в этом эндпоинте — для эмитента используй [`find_securities`](find_securities.md) (поля `emitent_title`, `emitent_inn`).
- Тикер не найден → пустой dict `{}`.
- Все значения — **строки** (включая числа и булевы), приводи типы сам.
