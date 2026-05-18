# Tooling

Все CLI-утилиты и языковые рантаймы проекта ставятся **только через [mise](https://mise.jdx.dev/)**.

У репозитория **три независимых файла `mise.toml`** — по одному на модуль и
один на «ничью землю»:

- `server/mise.toml` — Go-toolchain (`go`, `golangci-lint`, `buf`, `mockery`,
  protoc-плагины, `govulncheck`, `go-test-coverage`, `structfielddocs`,
  `go-arch-lint`) и линтеры scope `server/` (`yamllint`, `markdownlint-cli2`,
  `prettier`, `cspell` + `@cspell/dict-ru_ru`, `lychee`).
- `web/mise.toml` — фронт-toolchain (`node`, `buf`, `npm:@bufbuild/protoc-gen-es`)
  и линтеры scope `web/`.
- `mise.toml` (корень) — линтеры «ничьей земли» репозитория (`.github/`,
  `.claude/**`, корневой `README.md`): `yamllint`, `markdownlint-cli2`,
  `prettier`, `cspell`, `lychee`, `actionlint`.

Версии — единственная точка правды (локально и в CI). Дублирование версий
утилит между `server/mise.toml`, `web/mise.toml` и корневым `mise.toml` —
сознательная цена изоляции toolchain между модулями (расхождение ловит CI на
любом расхождении версий правил).

**Запрещено:**

- `go install`, `go tool`, `brew install`, `curl | sh`, `apt install` для проектных утилит.
- Задачи `setup` / инсталляторы binaries в `Taskfile.yml`.
- Дублирование версий в CI workflow (CI ставит через `jdx/mise-action@v2` с
  `working_directory: <module>` для модульных workflow).

**Исключения** допустимы только для того, что mise поставить не может (сам mise как пример).
