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

## Использование

Сервер предоставляет API на `/api/v1/*`.

Метрики Prometheus доступны по `/metrics`.

WebSocket подключение для событий доступно на `/api/ws`.

SDK-интерфейс доступен по `/sdk/v1/*`.
