---
endpoint: /iss/engines/{engine}/markets/{market}/boardgroups/{boardgroup}/securities.json
---

# `get_boardgroup_securities`

**Назначение:** все бумаги группы режимов с marketdata (Blocks).

**Reference ISS:** <https://iss.moex.com/iss/engines/{engine}/markets/{market}/boardgroups/{boardgroup}/securities>

## URL

```http
GET https://iss.moex.com/iss/engines/{engine}/markets/{market}/boardgroups/{boardgroup}/securities.json
```

**Форма ответа:** ISS Blocks (несколько именованных блоков в одном payload).

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя          | Тип | Обязательно | Описание                                         |
| ------------ | --- | ----------- | ------------------------------------------------ |
| `BOARDGROUP` | int | да          | ID группы (см. `get_boardgroups`); 57 = Т+ Акции |

## Поля JSON-ответа

Аналог `get_market_securities`, но в рамках одной группы режимов.

## Edge cases

- BOARDGROUP=57 для shares = TQBR + аукцион закрытия SPEQ + ...
