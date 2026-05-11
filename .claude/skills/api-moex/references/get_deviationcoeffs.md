---
endpoint: /iss/statistics/engines/stock/deviationcoeffs.json
block: securities
---

# `get_deviationcoeffs`

**Назначение:** коэффициенты отклонения по бумагам (sigma, beta, f_plus, f_minus, spread).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/deviationcoeffs>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/deviationcoeffs.json
```

**Форма ответа:** блок `securities`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа (блок `securities`)

| Поле        | Тип   | Смысл                     |
| ----------- | ----- | ------------------------- |
| `tradedate` | date  | Дата                      |
| `secid`     | str   | Тикер                     |
| `liquidity` | str   | Класс ликвидности (L/M/I) |
| `sigma`     | float | Волатильность             |
| `beta`      | float | Коэффициент бета          |
| `f_plus`    | float | Фактор роста              |
| `f_minus`   | float | Фактор снижения           |
| `spread`    | float | Спред                     |

## Edge cases

- Single page (~1300 записей).
