---
endpoint: /iss/engines/stock/zcyc.json
---

# `get_zcyc`

**Назначение:** **кривая бескупонной доходности ОФЗ** (ZCYC) — текущая или на дату. Это **источник риск-фри ставки** для DCF.

**Reference ISS:** <https://iss.moex.com/iss/engines/stock/zcyc>

## URL

```http
GET https://iss.moex.com/iss/engines/stock/zcyc.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя    | Тип  | Описание                  |
| ------ | ---- | ------------------------- |
| `date` | date | Кривая на конкретную дату |

## Поля JSON-ответа (Blocks)

- `yearyields` — точки кривой в годах (period: 0.25, 0.5, 1, 2, 3, 5, 7, 10, 15, 20; value: %).
- `params` — коэффициенты модели Nelson-Siegel-Svensson (B0..B5, T1, G1..G9) для интерполяции произвольной дюрации.
- `securities` — точки на кривой по конкретным ОФЗ (clean price, yield, duration).

## Использование

Для горизонта 1–3 года:

```text
risk_free_rate ≈ yearyields[period == 3].value / 100
```

## Edge cases

- ZCYC — официальный источник РФ.
