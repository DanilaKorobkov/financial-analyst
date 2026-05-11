---
endpoint: /iss/statistics/engines/state/markets/repo/cboper.json
---

# `get_cboper_rates`

**Назначение:** WAP ставки операций ЦБ по периодам (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/statistics/engines/state/markets/repo/cboper>

## URL

```http
GET https://iss.moex.com/iss/statistics/engines/state/markets/repo/cboper.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

Блоки: `date` (текущие ставки), `dates` (история).

| Поле             | Тип   | Смысл                          |
| ---------------- | ----- | ------------------------------ |
| `DAYS`           | str   | Тенор (1_DAY, 1_WEEK, 1_MONTH) |
| `TRADEDATE`      | str   | Дата                           |
| `WADEPSRATE`     | float | Взвешенная депозитная ставка   |
| `WAREPORATE`     | float | Взвешенная REPO ставка         |
| `WAREPORATEFIXN` | float | Фиксированная REPO             |
| `TITLE`          | str   | Название тенора                |

## Edge cases

- Данные могут быть устаревшими (часть старых тенорам последний раз обновлялась годы назад).
