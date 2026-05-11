---
endpoint: /iss/history/engines/stock/zcyc.json
block: params
---

# `get_zcyc_history`

**Назначение:** intraday-кривые ZCYC за **текущий** торговый день (~20 тыс. записей по 10-минутному шагу).

**Reference ISS:** <https://iss.moex.com/iss/history/engines/stock/zcyc>

## URL

```http
GET https://iss.moex.com/iss/history/engines/stock/zcyc.json
```

**Форма ответа:** блок `params`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Поля блока `params`: tradedate, tradetime, b1, b2, b3, t1, g1..g9 — коэффициенты NSS на каждый момент дня.

## Edge cases

- ISS **игнорирует** `from/till` и всегда отдаёт только сегодня. Если нужна история на даты — используй `get_zcyc --date YYYY-MM-DD`.
- Single-page (пагинация на этом эндпоинте не работает).
