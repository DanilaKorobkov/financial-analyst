# bootstrap

Ansible-каркас, который из чистой Ubuntu LTS делает узел, готовый к
установке Dokploy (отдельный шаг, не входит в этот каркас). На этом
шаге ставятся только системные вещи: apt-обновления, базовые пакеты,
timezone, автоматические security-патчи, swap, Docker, проверка sshd.

Никаких секретов внутри этого каркаса нет: секреты приложения живут
в **GitHub Actions Secrets** и доставляются в Dokploy через его API
из release workflow.

## Структура

```text
bootstrap/
├── ansible.cfg                              # inventory, roles_path, без vault
├── inventory/
│   ├── hosts.yml                            # группа prod, хост prod-vps (без IP)
│   ├── group_vars/all.yml                   # параметры geerlingguy.swap / .docker
│   └── host_vars/
│       └── prod-vps.yml.example             # шаблон: ansible_host, ansible_user, ключ
├── playbooks/
│   └── bootstrap.yml                        # common → swap → docker
├── roles/
│   └── common/                              # apt upgrade, timezone, unattended-upgrades, assert sshd
├── requirements.yml                         # geerlingguy.swap, geerlingguy.docker, community.general
├── tests/
│   └── smoke.yml                            # checklist «узел готов к Dokploy»
├── README.md
└── .gitignore
```

## Команды

Все задачи — из корня репо. Каркас Taskfile корня знает про `bootstrap/`
и подставляет `dir:` автоматически.

| Команда                       | Что делает                                                                |
| ----------------------------- | ------------------------------------------------------------------------- |
| `task bootstrap:deps`         | Ставит внешние роли (geerlingguy.\*) и коллекции в `bootstrap/.galaxy/`.  |
| `task bootstrap:lint`         | `ansible-lint` — статический анализ playbook'ов и ролей.                  |
| `task bootstrap:syntax-check` | `ansible-playbook --syntax-check` для `playbooks/bootstrap.yml`.          |
| `task bootstrap:check`        | Dry-run прогон playbook'а (без изменений). Требует заполненный host_vars. |
| `task bootstrap:run`          | Реальный прогон playbook'а на целевой VPS.                                |
| `task bootstrap:smoke`        | Прогон `tests/smoke.yml` — проверка состояния узла после `bootstrap:run`. |

`bootstrap:lint` и `bootstrap:syntax-check` подключены в общий
`task checks` и в матрицу CI (`.github/workflows/repo-checks.yml`).

## Локальное тестирование на Lima

Lima — лёгкая Linux-VM для macOS, нативная под Apple Silicon. Не путать
с embedded Lima внутри `avito devenv` — это отдельный публичный бинарник.

```sh
brew install lima
limactl start --name=fa-bootstrap --cpus=2 --memory=2 --disk=20 template://ubuntu-lts
# IP VM
limactl shell fa-bootstrap ip -4 addr show lima0 | awk '/inet /{print $2}' | cut -d/ -f1
```

Скопируй `inventory/host_vars/prod-vps.yml.example` в
`inventory/host_vars/fa-bootstrap.yml`, подставь полученный IP, имя
пользователя macOS (Lima пробрасывает его в VM), и путь к SSH-ключу
Lima — обычно `~/.lima/_config/user`.

```sh
task bootstrap:deps
task bootstrap:check
task bootstrap:run
task bootstrap:smoke
# Второй прогон — должен быть нулевой changed (idempotency check)
task bootstrap:run
limactl delete fa-bootstrap --force
```

Особенности Lima:

- SSH-ключ Lima создаётся при первом старте, путь — `~/.lima/_config/user`.
- `ansible_user` — твой macOS-логин (Lima пробрасывает его автоматически).
- `sudo` в VM работает без пароля по умолчанию — `become: true` сработает
  без `become_pass`.

**Альтернативы:** Multipass (`brew install --cask multipass`), Vagrant +
VirtualBox. Принцип тот же: поднимаешь чистую Ubuntu LTS, заполняешь
`host_vars/<имя>.yml`, гоняешь те же `task bootstrap:*` команды.

## Прогон на боевом VPS

1. Скопируй `inventory/host_vars/prod-vps.yml.example` → `prod-vps.yml`.
2. Заполни `ansible_host`, `ansible_user` (обычно `root` для свежего VPS
   от провайдера), `ansible_ssh_private_key_file`.
3. Прогон:

   ```sh
   task bootstrap:deps
   task bootstrap:check        # dry-run, убедиться что план изменений ожидаемый
   task bootstrap:run
   task bootstrap:smoke
   ```

4. Повторный `task bootstrap:run` обязан показать **0 changed** —
   idempotency invariant.

## sshd hardening: assert, а не enforce

Роль `common` **не** правит `/etc/ssh/sshd_config`. Она проверяет, что
эффективный sshd-конфиг уже соответствует требованиям:

- `PasswordAuthentication no`
- `PermitRootLogin prohibit-password` (или `no`)

У публичных провайдеров VPS Ubuntu cloud-образы выкатываются именно так.
Если playbook упал на этом assert'е — провайдер выдал нестандартный образ;
поправь sshd_config вручную (через консоль панели провайдера) и
перезапусти playbook. Это сознательное проектное решение, чтобы исключить класс
ошибок «playbook автоматически перезаписал sshd_config и я потерял
доступ».

## Что НЕ делается тут

- **Установка Dokploy** — отдельный шаг (PR 4).
- **Создание deploy-юзера, отключение root-логина** — на VPS одного
  разработчика второй пользователь не добавляет реальной безопасности
  и поднимает риск SSH lockout. Hardening сводится к assert'у выше.
- **ufw / firewall** — ufw на Docker-хосте обходится Docker'ом (DOCKER-USER
  chain). Правильная сетевая политика будет в PR Dokploy: bind Dokploy-панели
  на `127.0.0.1` и доступ через SSH-туннель, а не через ufw.
- **Ansible Vault и любые секреты приложения** — секреты живут в GitHub
  Actions Secrets, доставляются в Dokploy через его API из release workflow.
