---
endpoint: /iss/securities/{TICKER}/dividends.json
block: dividends
---

# `get_security_dividends`

**Назначение:** история дивидендов бумаги из публичного MOEX ISS.

**Reference ISS:** <https://iss.moex.com/iss/securities/{SECID}/dividends>

## URL

```http
GET https://iss.moex.com/iss/securities/{TICKER}/dividends.json
```

**Форма ответа:** блок `dividends`.

**Базовые query-параметры ISS** (общие для всех эндпоинтов, в reference не дублируются — см. SKILL.md): `iss.json=extended`, `iss.meta=off`, `iss.only=<block>`, `<block>.columns=<csv>`, `start=<offset>` для пагинации.

### Параметры

| Имя               | Тип | Обязательно | Описание     |
| ----------------- | --- | ----------- | ------------ |
| `{TICKER}` (path) | str | да          | SECID бумаги |

## Поля JSON-ответа

| Поле                | Тип   | Смысл                                       |
| ------------------- | ----- | ------------------------------------------- |
| `secid`             | str   | Тикер                                       |
| `isin`              | str   | ISIN бумаги                                 |
| `registryclosedate` | date  | Дата закрытия реестра (≠ ex-date; ex = T-2) |
| `value`             | float | Размер дивиденда на 1 акцию                 |
| `currencyid`        | str   | Валюта (RUB, USD)                           |

## Edge cases

- ISS отдаёт **усечённую** историю (~5–10 последних событий, иногда с 2014). Глубокая история — `api-financemarker stock --section dividends`.
- Для редомицилированных бумаг история ведётся **только по новому тикеру** (например, для YDEX нет истории YNDX).
- VTBR до сплита 2024-07-15: дивиденды по старому лоту (~0.001 RUB) — для adjusted-расчёта применяй коэффициент сплита из `splits`.
- `registryclosedate` — дата фиксации реестра, **не** ex-dividend date. Ex-date в режиме T+2 = `registryclosedate − 2 рабочих дня`.
