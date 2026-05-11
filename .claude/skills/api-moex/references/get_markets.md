---
endpoint: /iss/engines/{engine}/markets.json
block: markets
---

# `get_markets`

**Назначение:** список рынков торговой системы.

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets.json
```

**Форма ответа:** блок `markets`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Дефолт  | Описание         |
| ----------------- | --- | ------- | ---------------- |
| `{engine}` (path) | str | `stock` | Торговая система |

## Поля JSON-ответа

| Поле                                    | Тип | Смысл                                                         |
| --------------------------------------- | --- | ------------------------------------------------------------- |
| `id`                                    | int | ID рынка                                                      |
| `trade_engine_id` / `trade_engine_name` | …   | ID/код движка                                                 |
| `name`                                  | str | Код для `--market` (shares, bonds, foreignshares, index, ...) |
| `title`                                 | str | Название рынка                                                |
| `marketplace`                           | str | Площадка                                                      |
| `is_otc`                                | int | OTC флаг                                                      |
| `has_yield`                             | int | Считается ли доходность                                       |

## Edge cases

- Для `engine=stock`: shares, bonds, foreignshares, index, ndm, repo, qnv, mamc, …
