---
endpoint: /iss/sdfi/curves.json
block: curves
---

# `get_sdfi_curves`

**Назначение:** справочник своп-кривых (~70 кривых RUB/USD/EUR/CNY).

**Reference ISS:** <https://iss.moex.com/iss/sdfi/curves>

## URL

```http
GET https://iss.moex.com/iss/sdfi/curves.json
```

**Форма ответа:** блок `curves`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле          | Тип  | Смысл                          |
| ------------- | ---- | ------------------------------ |
| `curveid`     | str  | ID кривой                      |
| `date_from`   | date | Начало истории                 |
| `date_till`   | date | Последняя актуальная дата      |
| `title`       | str  | Описание                       |
| `methodology` | str  | Источник методологии           |
| `unit`        | str  | Единица (BasisPoints/Fraction) |

## Edge cases

- Детали конкретной кривой (`/iss/sdfi/curves/{id}`) — **платно**.
