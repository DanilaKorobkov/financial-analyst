# Web · Code

Конвенции frontend-модуля (`web/`). Стек: Vite + React + TypeScript +
Tailwind v4 + shadcn/ui (пресет Nova: Radix + Lucide + Geist).

## UI-компоненты: всё из shadcn/ui

**Главное правило:** любую видимую часть интерфейса собираем из готовых
компонентов [ui.shadcn.com](https://ui.shadcn.com/docs/components). Свои
визуальные компоненты не пишем.

### Порядок действий, когда нужен новый элемент UI

1. **Найти в [ui.shadcn.com/docs/components](https://ui.shadcn.com/docs/components).**
   Каталог покрывает почти всё, что встречается в продуктовых интерфейсах:
   `button`, `input`, `card`, `table`, `badge`, `skeleton`, `alert`, `tabs`,
   `dialog`, `sheet`, `dropdown-menu`, `command`, `tooltip`, `popover`,
   `select`, `combobox`, `chart`, `sonner`, `form` и т.д.
2. **Установить через CLI** в `web/`:

   ```bash
   npx shadcn@latest add <component> [<component> ...]
   ```

   Файлы появятся в `web/src/components/ui/`. Их **не редактируем** — это
   зеркало registry, обновляется командой выше.

3. **Использовать через alias `@/`**:

   ```tsx
   import { Button } from '@/components/ui/button';
   import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
   ```

4. **Если нужного блока нет в основном каталоге** — посмотреть готовые
   [блоки](https://ui.shadcn.com/blocks) (dashboard-01..07, login, sidebar)
   и взять блок целиком: `npx shadcn@latest add <block-name>`.

### Запрещено

- Писать «свой Button / Input / Card / Modal / Table / Tooltip и т.п.» —
  параллельно тому, что уже есть в shadcn. Даже если хочется «чуть-чуть
  другой стиль» — настраиваем через `className` и CSS-переменные темы
  (`--primary`, `--muted-foreground`, …) в `src/index.css`.
- Тащить параллельные UI-библиотеки (MUI, Chakra, Mantine, Ant Design,
  HeroUI и т.п.). Дизайн-система одна — shadcn поверх Radix.
- Импортировать `radix-ui` / `@radix-ui/*` напрямую в продуктовый код.
  Эти примитивы уже обёрнуты в `src/components/ui/*` — пользуемся ими.

### Что считается «своим кодом» и допустимо

- **Композиция shadcn-примитивов** под конкретный экран — нормально:
  собрать `Card + CardHeader + Badge` в виде локальной функции внутри
  страницы. Это data-driven layout, не новый компонент дизайн-системы.
- **Утилиты, не относящиеся к UI** — формат значений (`formatValue`,
  `humanizeKey`), API-клиенты (`api.ts`), хуки доступа к данным. shadcn
  таких вещей не отдаёт по определению.

### Иконки

Все иконки — из `lucide-react` (`Loader2Icon`, `AlertCircleIcon`,
`InfoIcon`, …). Это иконпак, на котором собран Nova-пресет shadcn.
Свой SVG или альтернативный набор иконок (Heroicons, react-icons,
`@tabler/icons-react`) не подключаем.

### Вспомогательная функция `cn()`

Склейка Tailwind-классов всегда через `cn()` из `@/lib/utils` — он
правильно мержит конфликтующие классы Tailwind (`clsx` + `tailwind-merge`).
Не использовать голый шаблонный литерал для динамических классов.

```tsx
import { cn } from '@/lib/utils';

<div className={cn('px-4 py-2', isActive && 'bg-primary text-primary-foreground', className)} />
```

## Связь с сервером

Backend — Connect-RPC на `:8080` (см. `server/`). С фронта дёргаем по
протоколу Connect-JSON: `POST /<service>.<Method>` с
`Content-Type: application/json`. CORS-настройки на сервере нет — в dev
работает прокси Vite (`vite.config.ts`), весь `/company.v1.*` маршрут
перенаправляется на `:8080`.

TypeScript-клиент **генерируется** из proto-схем сервера через `buf` —
артефакты лежат в `src/gen/` и хранятся в репозитории. Конфигурация
кодгена — `web/buf.gen.yaml` (плагин `@bufbuild/protoc-gen-es`); запуск —
`task web:generate-proto`. Руками описывать типы запросов и ответов
запрещено: дублирующийся контракт мгновенно расходится с реальным.

Сам клиент собирается через `createClient(<Service>, transport)` из
`@connectrpc/connect` — см. `src/api.ts` (≈10 строк, никаких ручных
запросов через `fetch`).

Дрейф между proto-схемой и сгенерированным клиентом ловится в CI задачей
`task web:lint-generated-proto`: она запускает `buf generate` во временный
каталог и сравнивает с зафиксированным в репозитории `src/gen/`. Если
расходится — PR красный. Поэтому изменение proto и его клиента физически
попадают в один PR.
