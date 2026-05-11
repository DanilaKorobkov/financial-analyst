---
endpoint: /iss/securities/{TICKER}/indices.json
block: indices
---

# `get_security_indices`

**Назначение:** индексы, в которые входит (или входила) бумага.

**Reference ISS:** <https://iss.moex.com/iss/securities/{SECID}/indices>

## URL

```http
GET https://iss.moex.com/iss/securities/{TICKER}/indices.json
```

**Форма ответа:** блок `indices`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание     |
| ----------------- | --- | ----------- | ------------ |
| `{TICKER}` (path) | str | да          | SECID бумаги |

## Поля JSON-ответа

| Поле            | Тип   | Смысл                                           |
| --------------- | ----- | ----------------------------------------------- |
| `SECID`         | str   | Код индекса (IMOEX, RTSI, MOEXBC, …)            |
| `SHORTNAME`     | str   | Короткое название                               |
| `FROM`          | date  | Дата включения в индекс                         |
| `TILL`          | date  | Дата исключения (или последняя дата активности) |
| `CURRENTVALUE`  | float | Текущее значение индекса                        |
| `LASTCHANGEPRC` | float | Изменение значения индекса за день, %           |
| `LASTCHANGE`    | float | Изменение значения индекса в пунктах            |

## Edge cases

- Бумага не входит ни в один индекс → пустой массив `data` в блоке.
- Полезно для health-check: бумага в IMOEX → ликвидная голубая фишка.
- Параметр `FROM > today` означает запланированное включение.
