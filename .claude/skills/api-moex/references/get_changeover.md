---
endpoint: /iss/history/engines/stock/markets/shares/securities/changeover.json
block: changeover
---

# `get_changeover`

**Назначение:** история смены торговых кодов (старый SECID → новый SECID). Полезно для редомицилированных и реорганизованных эмитентов.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/stock/markets/shares/securities/changeover>

## URL

```http
GET https://iss.moex.com/iss/history/engines/stock/markets/shares/securities/changeover.json
```

**Форма ответа:** блок `changeover`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле          | Тип  | Смысл                               |
| ------------- | ---- | ----------------------------------- |
| `action_date` | date | Дата смены                          |
| `old_secid`   | str  | Старый тикер (например, YNDX, TCSG) |
| `new_secid`   | str  | Новый тикер (YDEX, T)               |

## Edge cases

- Включает не только редомициляцию (YDEX, T, HEAD), но и переименования эмитентов.
