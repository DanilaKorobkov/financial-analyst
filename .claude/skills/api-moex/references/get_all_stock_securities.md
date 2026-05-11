---
endpoint: /iss/referencedata/engines/stock/markets/all/securities.json
block: securities
paginated: true
---

# `get_all_stock_securities`

**Назначение:** полный каталог фондового рынка (~5000+ записей с ISIN, REGNUMBER, INN, FACEVALUE, ISSUESIZE, LISTLEVEL, TRADESTATUS).

**Reference ISS:** <https://iss.moex.com/iss/referencedata/engines/stock/markets/all/securities>

## URL

```http
GET https://iss.moex.com/iss/referencedata/engines/stock/markets/all/securities.json
```

**Форма ответа:** блок `securities` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип  | Описание                   |
| ----------------- | ---- | -------------------------- |
| `date`            | date | Каталог на конкретную дату |
| `<block>.columns` | csv  | Подмножество колонок       |

## Поля JSON-ответа

Большой набор полей: SECID, NAME, SHORTNAME, ISIN, REGNUMBER, INN, FACEVALUE, FACEUNIT, ISSUESIZE, LISTLEVEL, TRADESTATUS, INSTRID, SECTORID, ...

## Edge cases

- Пагинация автоматическая (~55 страниц по 100 записей). Время выполнения ~5–10 секунд.
- Полезно для скрининга по уровню листинга, размеру выпуска, секторам.
