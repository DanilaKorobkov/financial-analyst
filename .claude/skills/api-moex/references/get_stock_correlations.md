---
endpoint: /iss/statistics/engines/stock/markets/shares/correlations.json
block: coefficients
---

# `get_stock_correlations`

**Назначение:** beta и корреляции акций (single page по умолчанию, ~1000 строк).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/stock/markets/shares/correlations>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/stock/markets/shares/correlations.json
```

**Форма ответа:** блок `coefficients`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя    | Тип  | Описание                |
| ------ | ---- | ----------------------- |
| `date` | date | Срез на конкретную дату |

## Поля JSON-ответа (блок `coefficients`)

Поля: tradedate, secid, beta, correlation, … (пары bencmark/security).

## Edge cases

- Полный набор ~70 тыс. записей; CLI отдаёт только первую страницу. Для полной выгрузки — клиентский цикл по `--start` (не реализован в CLI).
