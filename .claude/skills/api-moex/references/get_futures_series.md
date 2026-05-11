---
endpoint: /iss/statistics/engines/futures/markets/forts/series.json
block: series
---

# `get_futures_series`

**Назначение:** все серии фьючерсов FORTS.

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/futures/markets/forts/series>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/futures/markets/forts/series.json
```

**Форма ответа:** блок `series`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле               | Тип  | Смысл                                           |
| ------------------ | ---- | ----------------------------------------------- |
| `secid`            | str  | Код серии (например, SiM7)                      |
| `name`             | str  | Имя серии (Si-6.27)                             |
| `start_date`       | date | Дата запуска                                    |
| `expiration_date`  | date | Дата экспирации (или 2100-01-01 для бессрочных) |
| `asset_code`       | str  | Код базового актива                             |
| `underlying_asset` | str  | Базовый актив (тикер бумаги/пары)               |
| `is_traded`        | int  | 1 = торгуется                                   |
