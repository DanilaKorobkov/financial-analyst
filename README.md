# financial-analyst

[![checks](https://img.shields.io/github/actions/workflow/status/DanilaKorobkov/financial-analyst/checks.yml?branch=main&label=checks&cacheSeconds=60)](https://github.com/DanilaKorobkov/financial-analyst/actions/workflows/checks.yml)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)

Программный финансовый аналитик: оценивает компании и индустрии несколькими независимыми методами.

## Установка окружения

Все инструменты проекта (Go, golangci-lint, task, yamllint и т.д.) ставятся через [mise](https://mise.jdx.dev/) — версии зафиксированы в [`mise.toml`](./mise.toml).

### 1. Поставить mise

macOS:

```bash
brew install mise
```

Linux/macOS (универсальный установщик):

```bash
curl https://mise.run | sh
```

### 2. Активировать shell-хук

Без хука версии не переключаются автоматически при `cd` в проект. Выбрать по своему shell — см. [официальную инструкцию](https://mise.jdx.dev/getting-started.html#activate-mise).

zsh:

```bash
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc
```

bash:

```bash
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
```

После — перезапустить терминал.

### 3. Поставить инструменты проекта

В корне репозитория:

```bash
mise install
```

Команда поднимает всё из `mise.toml`: Go, golangci-lint, task, python+pipx, yamllint.

### 4. Проверить

```bash
task fmt checks
```

Если зелёно — окружение готово.
