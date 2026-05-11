---
endpoint: /api/fm/v2/experts
---

# `experts`

**Назначение:** рейтинг аналитиков (брокеров/комьюнити), публикующих идеи в FM. Метрики: количество идей, % успешных, средняя доходность, средний срок.

**URL:** `GET /api/fm/v2/experts`

## Query-параметры
Стандартный набор пагинации (см. [SKILL.md](../SKILL.md)).

## Пример запроса
```
GET /api/fm/v2/experts?limit=10&sort_by=ranking_all&sort_order=asc
```

## Поля JSON-ответа
Массив объектов. Префиксы: `t_*` — total (за всё время), `c_*` — current (текущий год / период расчёта рейтинга).

| Поле | Тип | Смысл |
|---|---|---|
| `community` / `community_id` | str / int | Имя и ID аналитика (используется как `community_id` в [`ideas`](ideas.md)) |
| `category` | str | `PROF` (профи) и т.п. |
| `t_num` / `c_num` | int | Всего идей / в текущем периоде |
| `t_num_profit` / `c_num_profit` | int | Из них прибыльных |
| `t_perc_profit` / `c_perc_profit` | int | % прибыльных |
| `c_perc_success` | int | % достигших target |
| `t_av_profit` / `c_av_profit` | int | Средняя доходность по идеям (%) |
| `t_av_duration` / `c_av_duration` | int | Средний срок жизни идеи (мес) |
| `ranking_all` / `ranking_year` / `ranking_month` | float | Места в рейтингах |
| `changed_at` | datetime | Время обновления |

## Пример ответа
```json
[
  {
    "community": "Финам",
    "community_id": 40581,
    "category": "PROF",
    "t_num": 1350,
    "t_num_profit": 779,
    "t_perc_profit": 58,
    "t_av_profit": 2,
    "t_av_duration": 9,
    "c_num": 1312,
    "c_num_profit": 769,
    "c_perc_profit": 59,
    "c_perc_success": 42,
    "c_av_profit": 2,
    "c_av_duration": 8,
    "ranking_all": 3.0,
    "ranking_year": 2.5,
    "ranking_month": 2.0,
    "changed_at": "2026-05-01T01:40:03"
  }
]
```

## Edge cases
- `community_id` из этой выдачи можно использовать как фильтр для скрипта поверх [`ideas`](ideas.md) (FM-API не предоставляет фильтр `?community_id=` — фильтруй на стороне клиента).
- Числовые `*_profit` округлены до целых процентов; для точных цифр используй [`ideas`](ideas.md).
