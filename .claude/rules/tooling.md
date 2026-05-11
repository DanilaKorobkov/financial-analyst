# Tooling

Все CLI-утилиты и языковые рантаймы проекта ставятся **только через [mise](https://mise.jdx.dev/)**.
Версии зафиксированы в `mise.toml` — единственная точка правды (локально и в CI).

**Запрещено:**

- `go install`, `go tool`, `brew install`, `curl | sh`, `apt install` для проектных утилит.
- Задачи `setup` / инсталляторы binaries в `Taskfile.yml`.
- Дублирование версий в CI workflow (CI ставит через `jdx/mise-action@v2`).

**Исключения** допустимы только для того, что mise поставить не может (сам mise как пример)
