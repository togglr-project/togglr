# 📑 ТЗ: WebSocket Broadcaster для Pending Changes и Audit Log

## 🎯 Цель

Обеспечить доставку пользователям фронтенда актуальных событий о фичах и связанных сущностях в реальном времени. Устранить неконсистентность UI, когда изменения применены в бэкенде, но не отражаются на фронтенде.

---

## 📡 Архитектура

### 1. Источник событий

* В базе данных есть объединяющая view `v_realtime_events`, которая агрегирует события из таблиц:

    * `pending_changes` (+ `pending_change_entities`);
    * `audit_log`.
* Каждая запись содержит:

    * `source` (`pending` или `audit`),
    * `event_id` (строка),
    * `project_id`,
    * `environment_id` и `environment_key`,
    * `entity` и `entity_id`,
    * `action` (`created`, `updated`, `deleted`, `pending`, и т. д.),
    * `created_at`.

### 2. Фоновый воркер

* Запускается в backend.
* Раз в N секунд выполняет запрос:

  ```sql
  SELECT * FROM v_realtime_events WHERE created_at > $last_seen ORDER BY created_at ASC;
  ```
* Обновляет `last_seen` после каждой выборки.
* Каждое новое событие трансформируется в JSON и передаётся в broadcaster.

### 3. Broadcaster

* Компонент в памяти backend.
* Хранит мапу активных WebSocket-соединений:

  ```
  (project_id, environment_id) → []connections
  ```
* При получении события воркер вызывает `Broadcast(project_id, env, event)`, и все подписанные клиенты получают сообщение.

### 4. WebSocket-соединение

* Endpoint: `/api/ws?project_id=...&env=...`
* Авторизация: `Authorization: Bearer <token>`
* Клиент получает только те события, которые относятся к указанному `project_id` и `environment`.

---

## 📡 Формат сообщения

Сообщение от сервера к клиенту:

```json
{
  "source": "audit",
  "type": "feature_updated",
  "timestamp": "2025-09-26T12:00:00Z",
  "project_id": "abc123",
  "environment": "prod",
  "entity": "feature",
  "entity_id": "f1",
  "action": "updated"
}
```

---

## 🖥 Frontend

1. Подключается к `/api/ws?project_id=abc123&env=prod`.
2. Получает сообщения и обновляет локальное состояние (store).
3. Обработка типовых событий:

    * `feature_updated` → обновить конкретную фичу в store;
    * `feature_deleted` → удалить из store;
    * `pending_change_created` → показать в UI «ожидает применения»;

---

## ✅ Критерии готовности MVP

* Backend поднимает `/api/ws`, клиенты могут подключиться.
* Воркер каждые N секунд проверяет `v_realtime_events` и передаёт новые записи.
* При изменении фичи (pending или audit) фронтенд получает событие и UI обновляется без перезагрузки.

---

