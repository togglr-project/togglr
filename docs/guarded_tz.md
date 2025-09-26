Делаем полноценное техническое задание по **Guarded workflow** с учётом всех деталей: `pending_changes.change` как массив `entities`, гарантии безопасности для 1-пользователя (пароль / TOTP), запрет параллельных pending на одну сущность, транзакционность `audit_log`, поведение существующих ручек (`PUT/DELETE /api/v1/features/{id}` и т.д.), таблицы для настроек/аппруверов и API для создания/аппрува/отклонения.

Ниже — детальная спецификация, SQL DDL, последовательности (flows), примеры ответов API, правила безопасности и предложения по реализации.

---

# Техническое задание: Guard Engine / pending\_changes workflow

## Ключевые требования (кратко)

1. Все изменения, затрагивающие сущности, помеченные тегом `guarded`, идут через `pending_changes`.
2. `pending_changes.change` = JSON с полем `entities[]` (каждый — { entity, entity\_id, action, changes{old,new} }).
3. Для проектов с >1 участника изменение требует APPROVAL другого пользователя. Для проектов с 1 участником изменение также проходит через `pending_changes`, но допускается **авто-аппрув** при проверке пароля или TOTP (MFA).
4. Нельзя иметь одновременно >1 `pending`-записей, которые затрагивают одну и ту же сущность (entity+entity\_id).
5. Применение `pending`-изменения (approve) выполняется в одной транзакции: изменения в целевых таблицах + запись(и) в `audit_log` + перевод статуса pending→approved.
6. API existing PUT/DELETE должны возвращать понятный ответ, позволяющий фронту показать модал «требуется подтверждение» и / или запросить пароль / TOTP.
7. Проект должен иметь конфигурацию, кто может аппрувить (по умолчанию — project superusers). Хранится в `project_settings` / `project_approvers`.

---

## 1. Структура БД — DDL

### Таблица pending\_changes + сущности

```sql
CREATE TABLE pending_changes (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    requested_by varchar(255) NOT NULL,   -- username или user id
    request_user_id integer,              -- optional FK to users.id if available
    change jsonb NOT NULL,                -- see format below: { entities: [...] , meta: {...} }
    status varchar(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, cancelled
    created_at timestamptz DEFAULT now() NOT NULL,
    approved_by varchar(255),
    approved_user_id integer,
    approved_at timestamptz,
    rejected_by varchar(255),
    rejected_at timestamptz
);

-- Разделяем сущности вовнутрь для гарантии уникальности
CREATE TABLE pending_change_entities (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    pending_change_id uuid NOT NULL REFERENCES pending_changes(id) ON DELETE CASCADE,
    entity varchar(50) NOT NULL,    -- e.g., 'feature', 'rule', 'feature_schedule'
    entity_id uuid NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL
);
-- уникальность: не должно быть двух pending на одну сущность
CREATE UNIQUE INDEX ux_pending_entity_unique ON pending_change_entities (entity, entity_id)
  WHERE (TRUE); -- we will enforce "only one pending state" via check function/trigger below
```

> Примечание: мы отдельно заносим каждую сущность в `pending_change_entities`, чтобы иметь возможность обеспечить уникальность и быстрые запросы конфликтов. Статус «pending» хранится в parent (`pending_changes.status`); при изменении статуса обновляем/консистентно учитываем это.

### Таблица project\_approvers / project\_settings

```sql
CREATE TABLE project_approvers (
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role varchar(50) DEFAULT 'approver',
    PRIMARY KEY (project_id, user_id)
);

-- simple key-value settings for project (can store approver policy etc.)
CREATE TABLE project_settings (
    id serial PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    value jsonb NOT NULL,
    UNIQUE (project_id, name)
);
```

Пример `project_settings`:

* `guard.require_second_approver: true/false`
* `guard.default_approvers: [user_id,...]`

### Audit log (используем существующую таблицу audit\_log)

`audit_log` у тебя уже есть — используем её и гарантируем, что при аппруве туда пишем записи **в той же транзакции**, где изменили сущности.

---

## 2. Формат `pending_changes.change`

Рекомендуемый формат — единый, читаемый, фронтенд-дружественный (вариант с diff):

```json
{
  "entities": [
    {
      "entity": "feature",
      "entity_id": "11111111-1111-1111-1111-111111111111",
      "action": "update",             // insert / update / delete
      "changes": {
        "enabled": { "old": true, "new": false },
        "default_variant": { "old": "A", "new": "B" }
      }
    },
    {
      "entity": "rule",
      "entity_id": "22222222-2222-2222-2222-222222222222",
      "action": "delete",
      "changes": {}
    }
  ],
  "meta": {
    "reason": "Disable due to payment outage",
    "client": "ui",
    "origin": "project-settings"
  }
}
```

* `entities` — массив, может включать любые сущности, которые могут быть изменены через UI (feature, rule, feature\_schedule, segment\_sync, etc.).
* `changes` — содержит `old` и `new` для отображения diff в UI и для проверки при apply.
* `meta` — произвольные данные для context.

---

## 3. Правила конкурентности / запрет параллельных pending

### Поведение:

* При создании новой pending (через UI/API) нужно проверить, что **никакая другая pending-запись (status = 'pending') не затрагивает ни одну из сущностей в `entities[]`**. Если такой конфликт есть — отказать с кодом 409 Conflict и сообщением, какие сущности уже заблокированы.

### Реализация:

* В транзакции при создании `pending_changes`:

  1. Вставить запись в `pending_changes` (status pending).
  2. Для каждой entity вставить row в `pending_change_entities`.
  3. Перед вставкой в `pending_change_entities` выполнить проверку:

     ```sql
     SELECT 1 FROM pending_change_entities pce
     JOIN pending_changes pc ON pc.id = pce.pending_change_id
     WHERE pce.entity = :entity AND pce.entity_id = :entity_id AND pc.status = 'pending';
     ```

     — если найдено — ROLLBACK and return 409.
* Дополнительно: можно завести триггер/функцию, которая запретит вставку в `pending_change_entities` при существующем pending — это лучше для защиты, если кто-то пишет в БД напрямую.

---

## 4. API: поведение существующих ручек и новые эндпоинты

### Поведение существующих PUT/DELETE `/api/v1/features/{feature_id}` и др.

* При вызове операции изменения:

  1. Проверить: feature принадлежит проекту; проверить наличие `guarded`-тега у фичи (SQL `EXISTS(SELECT 1 FROM feature_tags ft JOIN tags t ... WHERE t.slug='guarded' AND ft.feature_id = :id)`).
  2. Если нет — ведём себя как сейчас (выполняем изменение, 200 OK, write audit\_log).
  3. Если есть:

     * Создаём `pending_changes` со `change` = diff (entities array с entity==feature).
     * Определяем количество активных пользователей в проекте: `select count(*) from memberships where project_id=:p and is_active = true`.
     * Если `count > 1`:

       * Возвращаем `202 Accepted` или `409 Conflict`? Рекомендация: **202 Accepted** + body с `{ status: "pending", pending_id: "...", message: "Approval required" }`.
       * Frontend показывает модал, уведомляет аппруверов, не применяет его локально.
     * Если `count = 1`:

       * Не применять сразу. UI должен prompt for password/TOTP; then call `POST /pending_changes/:id/approve` with credentials OR call combined endpoint to create+approve in one request (see below).
       * If credentials valid → approve happens synchronously and endpoint returns 200 OK with updated resource.
       * If credentials invalid → 401/403 with message.
* **Важно**: API must not apply changes without pending flow for guarded features; even for single user the change must be recorded in `pending_changes` and approved.

### Recommended status codes & bodies

* Success immediate apply (no guarded): `200 OK` with updated resource.
* For guarded and multi-user (created pending): `202 Accepted`

```json
{
  "status": "pending",
  "pending_id": "uuid",
  "message": "Change is pending approval"
}
```

* For guarded and single-user flow where auto-approval applied (credentials provided and valid): `200 OK` with updated resource.
* For conflict when another pending exists: `409 Conflict` with details which entity is blocked.

---

## 5. Endpoints for pending changes

### Create pending (generic) — usually called by server side when operation touches guarded entity

* `POST /api/v1/pending_changes`
* Body:

```json
{
  "project_id": "uuid",
  "requested_by": "alice",
  "request_user_id": 123,
  "change": { /* format above */ }
}
```

* Response:

  * `201 Created` with `{ pending_id }` (but typically PUT/DELETE route will create pending and respond 202).
  * If conflict → `409 Conflict`.

### List / Get

* `GET /api/v1/pending_changes?project_id=...&status=pending`
* `GET /api/v1/pending_changes/{id}`

### Approve

* `POST /api/v1/pending_changes/{id}/approve`
* Body:

```json
{
  "approver_user_id": 456,
  "approver_name": "bob",
  "auth": { "method": "password" | "totp", "credential": "<plaintext-or-code>" }
}
```

* Server verifies:

  * approver has rights (see project\_approvers or default superuser logic),
  * auth credential is valid for approver (password check or TOTP via stored secret).
* If valid → server **applies changes inside a DB transaction** (see apply procedure below), writes audit\_log entries, updates `pending_changes.status = 'approved'`, `approved_by`, `approved_at`.
* Response: `200 OK` with result of applied changes (or 204 No Content).

### Reject

* `POST /api/v1/pending_changes/{id}/reject`
* Body: `{ rejected_by: user, reason: "..." }`
* Effect: set status to `rejected`, write audit\_log entry.

### Combined create+approve (for single-user shortcut)

* `POST /api/v1/pending_changes?auto_approve=true`
* Body includes `auth` credentials so server can attempt to auto-approve if project has 1 active user or project setting allows auto\_approve.
* If auto-approved → apply and return `200 OK`.
* If not → create pending and return `202 Accepted`.

---

## 6. Apply algorithm for `approve` (atomic procedure)

When a pending is approved:

1. Begin DB transaction.
2. Acquire advisory lock(s) to avoid race:

   * For each entity in pending.entities, call `pg_advisory_xact_lock(hashtext(entity || ':' || entity_id))` (or lock by feature id integer). This prevents concurrent applies/changes.
3. For each entity:

   * Validate that the current DB state matches the `old` values in `changes` (optional optimistic check). If mismatch → fail and return 409 or require manual resolution.
   * Apply update/insert/delete according to `action` and `changes.new`. Use parametrized SQL.
   * Write an entry into `audit_log` for that entity change (actor = approver, action = 'approve\:change', old\_value, new\_value, project\_id, request\_id = pending\_change\_id).
4. Update `pending_changes` row: set `status='approved'`, `approved_by`, `approved_at`.
5. Commit transaction.

If any step fails → rollback transaction, leave `pending_changes.status='pending'` and return error.

**Important**: `audit_log` writes should be part of same transaction so auditing is consistent.

---

## 7. Authorization & approvers resolution

* Determine approver set:

  1. If `project_settings` has `guard_approvers` configured → use that list.
  2. Else use `project_approvers` table if populated.
  3. Else fallback to any user with membership role `project_owner` or global `is_superuser = true`.
* Approve endpoint must verify approver has permission.
* Approve endpoint must verify credentials provided:

  * password: compare hash with `users.password_hash` (use timing-attack safe compare).
  * totp: verify using `users.two_fa_secret` (if two\_fa\_enabled).

**Security note**: Do not log plaintext credentials. Use HTTPS. Rate-limit approve attempts.

---

## 8. UI / Frontend contract

* When user clicks Save on guarded entity:

  * Frontend calls backend normal PUT/DELETE endpoint.
  * Backend returns:

    * `202 Accepted` + pending\_id → UI shows message "Request pending approval" and optionally displays diff and list of approvers; also triggers UI notification to approvers (via WebSocket/email).
    * OR `200 OK` if auto-approved (single-user with valid MFA provided earlier).
    * OR `409 Conflict` if blocked by existing pending.
* If server responds `202` frontend should:

  * Show modal with pending details and "View pending changes" button.
  * Optionally allow the requesting user to enter password/TOTP to auto-approve (this would call `/pending_changes/{id}/approve` and, if successful, fetch updated resource and close modal).
* Approver UI:

  * List pending changes per project (GET `/pending_changes?status=pending&project_id=...`).
  * View details: show `entities[]` diffs (old → new).
  * Approve / Reject buttons. On Approve, request TOTP/password if needed.

---

## 9. Notifications & observability

* When creating a pending for project with >1 user — send notifications to approvers:

  * push via WebSocket / server-sent events,
  * email (if configured),
  * optional Slack webhook.
* Emit events to internal event bus: `pending.created`, `pending.approved`, `pending.rejected`.
* Metrics: `pending_created_total`, `pending_approved_total`, `pending_rejected_total`, `pending_conflicts_total`.

---

## 10. Edge cases & validation

* **Partial apply failure**: if applying second entity fails after first applied — rollback entire transaction (atomic).
* **Stale-old check**: optional optimistic check: compare `old` values in change with current DB before applying; if mismatch, return error and require refresh.
* **Large changes**: one pending can include many entities — UI should warn user if many items included.
* **Retry / idempotency**: Approve endpoint should be idempotent. If called twice, second call returns 200 and indicates already approved.
* **Cleanup**: consider TTL for pending items (e.g., auto-cancel after 30 days).
* **Audit linking**: use `request_id` field in `audit_log` to link change to pending id.

---

## 11. DB-side enforcement helpers (recommended)

* Function to enforce no duplicate pending for same entity:

```sql
CREATE OR REPLACE FUNCTION ensure_no_conflicting_pending()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    cnt int;
BEGIN
    -- NEW.pending_change_id exists, check for conflicts in other pending_changes
    SELECT count(*) INTO cnt
    FROM pending_change_entities pce
    JOIN pending_changes pc ON pc.id = pce.pending_change_id
    WHERE pce.entity = NEW.entity
      AND pce.entity_id = NEW.entity_id
      AND pc.status = 'pending'
      AND pce.pending_change_id <> NEW.pending_change_id;

    IF cnt > 0 THEN
        RAISE EXCEPTION 'Entity % % is already locked by another pending change', NEW.entity, NEW.entity_id;
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_pce_no_conflict
BEFORE INSERT ON pending_change_entities
FOR EACH ROW EXECUTE FUNCTION ensure_no_conflicting_pending();
```

* Trigger to cascade-update status field to children not required, but when `pending_changes.status` changes -> no need to keep child status duplicated. For query efficiency, you can index `pending_change_entities (entity, entity_id)`.

---

## 12. Flow diagrams (step-by-step narrative)

### A. Save pressed on guarded feature (project has >1 users)

1. UI PUT /api/v1/features/{id} with payload.
2. Server detects `guarded` tag → builds `change` diff and attempts to create `pending_changes` + `pending_change_entities` in a transaction.

   * If conflict → 409 Conflict to UI.
   * If created → 202 Accepted + `{ pending_id }`.
3. Server notifies approvers (WebSocket / email).
4. Pending is displayed in approver UI.

### B. Approver approves

1. Approver clicks Approve in UI, UI calls `POST /pending_changes/{id}/approve` with approver id and TOTP.
2. Server verifies approver authority and credentials.
3. Server calls apply-procedure:

   * begin tx
   * acquire advisory locks per entity
   * optionally verify `old` values
   * apply changes (UPDATE/DELETE/INSERT)
   * insert corresponding rows in `audit_log` (actor = approver, request\_id = pending\_id)
   * update pending\_changes.status='approved'
   * commit
4. Notify requester and update UI.

### C. Single-user flow

1. UI PUT /api/v1/features/{id} → server detects `guarded`.
2. Server creates pending record (status pending).
3. UI prompts user for password/TOTP and calls `POST /pending_changes/{id}/approve` with credentials.
4. Server verifies credential, approves and applies (same as B).
5. Return `200 OK` and updated resource to UI.

---

## 13. Полезные SQL примеры (apply transaction pseudo)

```sql
BEGIN;

-- lock entity to avoid race
SELECT pg_advisory_xact_lock(hashtext(concat(entity, ':', entity_id)));

-- optionally check current state equals old
-- apply update/insert/delete
UPDATE features SET enabled = <new> WHERE id = <entity_id>;

-- write audit log
INSERT INTO audit_log (feature_id, actor, action, old_value, new_value, created_at, project_id, request_id)
VALUES (...);

-- update pending_changes
UPDATE pending_changes SET status='approved', approved_by='bob', approved_at=now() WHERE id = :pending_id;

COMMIT;
```

---

## 14. Таблица/функции для MFA/password verification

* Пользуйся существующей таблицей `users` (`password_hash`, `two_fa_secret`, `two_fa_enabled`).
* Реализация проверки:

  * For password: bcrypt compare with stored hash.
  * For TOTP: use standard TOTP lib to verify `credential` against `two_fa_secret`.

**Security**: validate approver identity (must be authenticated session user) and confirm `approver_user_id` matches authenticated user. Never accept arbitrary `approver` param without session check.

---

## 15. Про проектные настройки (where to configure approvers and policy)

* `project_approvers` table (explicit list of project users who can approve).
* `project_settings` holds `guard.policy`:

```json
{
  "require_second_approver": true,
  "allow_auto_approve_if_single_member": true,
  "default_approvers": [123,456]
}
```

* UI: project settings UI to configure approvers and policies.

---

## 16. Auditing and retention

* При создании pending — write a minimal audit\_log entry: action = 'pending\_created', old\_value=NULL, new\_value=change payload, request\_id = pending\_id.
* При approve/reject — write audit\_log entries for each entity change, and one summary entry for pending approval action.
* Retain pending/review history for compliance.

---

## 17. Next-step checklist (implementation plan)

1. Добавить DDL: `pending_changes`, `pending_change_entities`, `project_approvers`, `project_settings`.
2. Добавить DB trigger `ensure_no_conflicting_pending` (as above).
3. Изменить обработчики PUT/DELETE и другие mutating API:

   * detect `guarded` tag,
   * build `change` payload,
   * create `pending_changes` (and entities) or attempt auto-approve for single-user flow.
4. Реализовать endpoints:

   * `POST /pending_changes` (create)
   * `GET /pending_changes`
   * `POST /pending_changes/{id}/approve`
   * `POST /pending_changes/{id}/reject`
5. Реализовать auth checks in approve: project approver resolution, password/TOTP validation.
6. Implement apply-procedure: transactional apply + audit\_log writes + advisory locks.
7. UI integration: return codes & payloads per spec, show modals for MFA/password when needed, list pending changes for approvers.
8. Notifications: websocket & email hooks for approvers.
9. Tests:

   * unit tests for DB triggers,
   * integration tests for whole flow (create pending -> approve -> check DB changes & audit),
   * concurrency test: simultaneous creates should produce 409 conflict.
10. Monitoring & metrics: pendings count, approves per hour, conflicts, failed approvals.

---

## 18. Дополнительные замечания и рекомендации (советую, потому что ты просил честно спорить)

* **Хранение credentials**: никогда не принимай пароль/TOTP в query string or logs. Use HTTPS and secure logging.
* **Rate-limit approve attempts** to avoid brute-force password/TOTP. Log failed attempts to audit\_log (but do not store provided credentials).
* **UI/UX**: четко показывай, что действие не применено, а создана pending. Покажи diff и список approvers. Для single-user flow — prompt for password/TOTP before sending initial request OR send it in approve call — second approach is cleaner.
* **Atomicity**: обязательно wrap apply in single DB transaction and use advisory locks to avoid racing applies.
* **Backwards compatibility**: existing clients calling PUT/DELETE will get new responses; coordinate frontend update to handle 202 and show proper modals.
