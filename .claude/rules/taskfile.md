# Taskfile

Каждая задача в `Taskfile.yml` обязана иметь поле `desc` — краткое описание того, что она делает.

```yaml
tasks:
  lint:
    desc: Run golangci-lint
    cmds:
      - golangci-lint run ./...
```

Без `desc` задачу добавлять нельзя.
