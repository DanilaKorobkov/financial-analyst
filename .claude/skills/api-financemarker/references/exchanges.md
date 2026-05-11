---
endpoint: /api/fm/v2/exchanges
---

# `exchanges`

**Назначение:** справочник поддерживаемых бирж. Без пагинации.

**URL:** `GET /api/fm/v2/exchanges`

## Query-параметры
Нет.

## Пример запроса
```
GET /api/fm/v2/exchanges
```

## Поля JSON-ответа
Массив объектов:

| Поле | Тип | Смысл |
|---|---|---|
| `exchange` | str | Код биржи (используется в `<exchange>:<code>`) |
| `name` | str | Название биржи на русском |
| `country` | str | ISO-код страны (`RU` и т.п.) |
| `currency` | str | ISO-код валюты торгов биржи |
| `mic` | str | Market Identifier Code (ISO 10383) |

## Пример ответа
На текущей подписке возвращается **только MOEX**:
```json
[
  {
    "exchange": "MOEX",
    "name": "Московская биржа",
    "country": "RU",
    "currency": "RUB",
    "mic": "MISX"
  }
]
```

## Edge cases
- Если в выдаче только `MOEX` — это ограничение тарифа FM, а не баг. Зарубежные тикеры (NASDAQ, NYSE и т.п.) этим токеном не получить.
- Перед обращением к `stocks/{exchange}:{code}` имеет смысл проверить здесь, что нужный `exchange` доступен.
