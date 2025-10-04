# Togglr

Togglr — это система управления feature flags и экспериментами для разработчиков. Она позволяет включать и выключать функциональность без релиза кода, проводить A/B-тестирование, управлять раскаткой по сегментам пользователей и отслеживать стабильность фич.

## Возможности

* Управление feature flags по проектам и окружениям (prod, stage, dev).
* Поддержка вариантов (variants) и правил таргетинга (rules).
* Guarded features (pending changes, approval workflow).
* Категории и теги для организации фич.
* Планировщик (schedules) для автоматического включения и выключения.
* Аудит лог изменений.
* SLA и health-мониторинг фич.
* Auto-disable при ошибках исполнения (через error reports).
* RBAC с ролями и разрешениями.
* REST API и SDK.

## Архитектура

* **Backend** — Go (PostgreSQL/TimescaleDB, NATS, REST API, WebSocket broadcaster).
* **Frontend** — React + TypeScript.
* **SDK**:

    * Go
    * Ruby
    * PHP
    * Python
    * TypeScript (Node.js и браузер)
    * Elixir

## Настройка dev-окружения

### Требования

- Docker
- Docker Compose

### Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone <repository-url>
   cd togglr
   ```

2. Настройте файлы окружения:
   ```bash
   make setup
   ```

3. Запустите dev-окружение:
   ```bash
   make dev-up
   ```

4. Откройте приложение:
   - Frontend: https://localhost
   - API: https://localhost/api/v1/
   - SDK: https://localhost/sdk/v1/

### Настройки

- **Домен**: По умолчанию приложение настроено на `localhost` в файлах `dev/config.env` и `dev/platform.env`
- **SSL сертификаты**: Требуются самоподписные сертификаты в директории `dev/nginx/ssl/`. Включены предварительно сгенерированные сертификаты, но они могут быть просрочены
- **Суперпользователь**: При первом запуске создается суперпользователь с:
  - Email: `ADMIN_EMAIL` из `dev/config.env` (по умолчанию: `admin@togglr.dev`)
  - Пароль: `ADMIN_TMP_PASSWORD` из `dev/config.env` (по умолчанию: `password543210`)
  - Эти данные можно изменить после первого входа

### Команды для разработки

- `make dev-up` - Запустить все сервисы
- `make dev-down` - Остановить все сервисы
- `make dev-clean` - Остановить сервисы и очистить volumes/images
- `make build` - Собрать приложение (требует Go 1.25+)
- `make test` - Запустить тесты (требует Go 1.25+)

## Использование

Сервер предоставляет API на `/api/v1/*`.

Метрики Prometheus доступны по `/metrics`.

WebSocket подключение для событий доступно на `/api/ws`.

SDK-интерфейс доступен по `/sdk/v1/*`.
