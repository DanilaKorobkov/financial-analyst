# financial-analyst

Программный финансовый аналитик: оценивает компании и индустрии несколькими независимыми методами.

## Установка окружения

Весь тулинг проекта (Go, golangci-lint, task, yamllint и т.д.) ставится через [mise](https://mise.jdx.dev/) — версии пинятся в [`mise.toml`](./mise.toml).

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
