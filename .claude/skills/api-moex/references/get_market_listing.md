---
endpoint: /iss/history/engines/{engine}/markets/{market}/listing.json
block: securities
paginated: true
---

# `get_market_listing`

**Назначение:** все когда-либо торговавшиеся инструменты рынка (включая делистнутые) с диапазоном дат истории.

**Reference ISS:** <https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/listing>

## URL

```http
GET https://iss.moex.com/iss/history/engines/{engine}/markets/{market}/listing.json
```

**Форма ответа:** блок `securities` с пагинацией через `start` (+ `history.cursor` где применимо).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

## Поля JSON-ответа

| Поле           | Тип  | Смысл                        |
| -------------- | ---- | ---------------------------- |
| `SECID`        | str  | Тикер (включая исторические) |
| `SHORTNAME`    | str  | Короткое имя                 |
| `NAME`         | str  | Полное имя                   |
| `BOARDID`      | str  | Режим торгов                 |
| `decimals`     | int  | Знаков после запятой         |
| `history_from` | date | Первая дата истории          |
| `history_till` | date | Последняя дата истории       |

## Edge cases

- Самый полный источник делистнутых тикеров. Для скрининга «что когда-либо торговалось».
- Пагинация автоматическая, занимает несколько секунд.
