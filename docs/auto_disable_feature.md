# Концепт: auto-disable **только для помеченных** фич

**Правило:** автоотключение срабатывает только если у фичи есть тег `auto-disable` (slug = `'auto-disable'`).
Если тега нет — репорты ошибок сохраняются, но auto-disable не выполняется.

## 🎯 Цель

Добавить поддержку приёма ошибок от SDK по эндпоинту
`POST /sdk/v1/features/{feature_key}/report-error`
и сохранения их в TimescaleDB (`monitoring.error_reports`), а также реализовать бизнес-логику auto-disable в зависимости от настроек проекта.

---

## 📑 Архитектурные составляющие

### 1. Таблица хранения

Используем `monitoring.error_reports` (см. миграцию, уже создана: migrations/034_error_reports.up.sql).
Назначение: хранение всех репортов об ошибках с привязкой к проекту, окружению и фиче.

Ключевые поля:

* `event_id` (UUID, бизнес-ключ)
* `project_id`
* `environment_id`
* `feature_id`
* `error_type`, `error_message`, `context`
* `created_at`

⚡ Особенности: hypertable + retention 30 дней (сами данные чистятся Timescale).

---

### 2. Репозиторий

Создать **репозиторий `ErrorReportRepository`**, методы:

```go
type ErrorReportRepository interface {
    Insert(ctx context.Context, report *domain.ErrorReport) error
    CountRecent(ctx context.Context, featureID domain.FeatureID, envID domain.EnironmentID, window time.Duration) (int, error)
    GetHealth(ctx context.Context, featureID domain.FeatureID, envID domain.EnironmentID) (*domain.FeatureHealth, error)
}
```

* `Insert` — вставка нового отчёта.
* `CountRecent` — посчитать количество ошибок за заданный интервал (например, 5 минут).
* `GetHealth` — агрегированное состояние фичи (ошибки, disabled, статус).

Реализация: `error_report_pg.go` (Postgres/Timescale).

---

### 3. Domain-модели

```go
type ErrorReport struct {
    EventID       string
    ProjectID     domain.ProjectID
    FeatureID     domain.FeatureID
    EnvironmentID domain.EnironmentID
    ErrorType     string
    ErrorMessage  string
    Context       map[string]any
    CreatedAt     time.Time
}

type FeatureHealth struct {
    FeatureID     domain.FeatureID
    EnvironmentID domain.EnironmentID
    Enabled       bool
    Status        string    // healthy | degraded | disabled
    ErrorRate     float64   // ошибки в %
    LastErrorAt   time.Time
}
```

---

### 4. UseCase

Создать **`ErrorReportUseCase`** (или `FeatureHealthUseCase`).

Методы:

* `ReportError(ctx, featureKey, envKey, report)`

    1. Найти проект/фичу/окружение по `feature_key` + `api_key` (идентификация SDK-запроса).
    2. Сохранить в `error_reports` через репозиторий.
    3. Проверить **правила auto-disable**:

        * Получить настройку проекта `auto_disable_requires_approval`.
        * Посчитать количество ошибок за последний интервал (например, 1 мин или N запросов).
        * Если порог превышен:

            * Если **автоотключение без approval**: сразу обновить `feature_params.enabled = false`.
            * Если **требуется approval**: создать запись в `pending_changes`.
    4. Вернуть `FeatureHealth` (200) или `202` (если change pending).

* `GetFeatureHealth(ctx, featureKey, envKey)`

    1. Посчитать статистику ошибок из `error_reports` (за окно, например, 1 мин).
    2. Проверить состояние фичи (`enabled/disabled`).
    3. Вернуть агрегат `FeatureHealth`.

---

### 5. Проверки и защита

* Валидация входных данных:

    * `error_type` не пустой, ограничение длины.
    * `context` должен быть JSON.
* Ограничение частоты (rate-limit per SDK client → опционально).
* Транзакционность:

    * Для `Insert + auto-disable` использовать **одну транзакцию** (чтобы не получилось: ошибка сохранилась, а автоотключение не применилось).

---

### 6. Настройки

Используем таблицу `project_settings`:

* name=`auto_disable_requires_approval` (bool)

    * `false` (default): автоотключение выполняется немедленно.
    * `true`: создаётся pending change (approve/reject).

⚡ В будущем можно добавить:

* `auto_disable_threshold` (кол-во ошибок).
* `auto_disable_window` (интервал времени).

---

## ⚙️ Поток данных (пример)

1. SDK вызывает `POST /sdk/v1/features/payment/report-error`.
2. Бэкенд извлекает `project_id`, `env_id` по API ключу.
3. Создаёт `error_reports` запись.
4. Считает ошибки за 1 мин → 25.
5. Порог = 20. Настройка проекта `auto_disable_requires_approval = false`.
6. Транзакция:

    * вставка error_report
    * update feature_params.enabled=false
7. Возврат: `200 { "status": "disabled", "errorRate": 0.7 }`.

---

## 📌 Реализация по слоям

* `internal/repository/error-reports`

    * `Insert` → `INSERT INTO monitoring.error_reports ...`
    * `CountRecent` → `SELECT count(*) FROM monitoring.error_reports WHERE feature_id=$1 AND environment_id=$2 AND created_at > now()-$3`
    * `GetHealth` → агрегаты (count, last_error_at).

* `internal/usecases/error-reports`

    * бизнес-логика.

* internal/api/sdk/feature_report_error.go, internal/api/sdk/feature_health.go

    * `ReportFeatureError` (POST)
    * `GetFeatureHealth` (GET).

---

Отлично 👍 Давай зафиксируем **алгоритм auto-disable** как бизнес-логику в ТЗ.

---

# 🔥 Алгоритм Auto-Disable Feature

### 🎯 Цель

Защитить систему от фич, которые начали массово падать на стороне SDK/клиента.
Если в течение короткого времени (окно) превышен порог ошибок → фича автоматически отключается или отправляется в pending change (по настройке проекта).

---

## 1. Входные параметры (project_settings)

Храним в `project_settings` (чтобы было настраиваемо для каждого проекта):

* `auto_disable_enabled` (bool, default = true) — включено/выключено.
* `auto_disable_requires_approval` (bool, default = false) — если true → создаётся pending change, а не прямое отключение.
* `auto_disable_error_threshold` (int, default = 20) — минимальное количество ошибок, чтобы сработало отключение.
* `auto_disable_time_window_sec` (int, default = 60) — окно времени в секундах для подсчёта ошибок.
* `auto_disable_error_rate` (float, default = 0.5) — доля ошибок среди запросов (error_rate > 50%).

    * ❗ для error_rate нам нужны ещё метрики **успешных запросов** → их можно добавлять в `feature_eval_reports` (MVP можно ограничиться только `threshold`).

---

## 2. Алгоритм проверки

При каждом `ReportFeatureError` (POST /sdk/v1/features/{feature_key}/report-error):

1. **Сохранить запись** в `monitoring.error_reports`.
2. Посчитать количество ошибок по `(feature_id, environment_id)` за `auto_disable_time_window_sec`.

   ```sql
   SELECT count(*) 
   FROM monitoring.error_reports
   WHERE feature_id = $1
     AND environment_id = $2
     AND created_at > now() - interval 'N seconds';
   ```
3. Если `count >= auto_disable_error_threshold`, то:

    * Если `auto_disable_requires_approval = false`:

        * Обновить `feature_params.enabled=false` (немедленно).
        * Записать событие в `audit_log`.
    * Если `auto_disable_requires_approval = true`:

        * Создать запись в `pending_changes` с action=disable.
        * Статус = pending → ждёт approve/reject.

---

## 3. Error Rate (опционально в будущем)

Если есть метрики **успешных запросов**:

* error_rate = errors / (errors + successes).
* Если `error_rate >= auto_disable_error_rate` (например, >50%) → срабатывает автоотключение.

👉 Пока можно реализовать **threshold-only** (по количеству ошибок), а error_rate добавить после интеграции метрик.

---

## 4. Пример сценариев

### ✅ Прямое отключение

* `threshold=20`, `window=60s`, `requires_approval=false`.
* За 1 минуту прилетает 25 ошибок.
* Feature → `enabled=false` в `feature_params`.
* В `audit_log`: `action=auto-disable`.
* SDK при следующей evaluate получит `disabled`.

### 🕒 Pending Change

* `threshold=10`, `window=30s`, `requires_approval=true`.
* За 30s прилетает 12 ошибок.
* Создаётся `pending_change (action=disable)`.
* Project owner видит в UI → может `approve` или `reject`.

---

## 5. Особенности реализации

* Операции `Insert error + disable feature` должны быть в **одной транзакции** (чтобы не было рассинхрона).
* Для защиты от «штормов» (много клиентов одновременно) → можно добавить `SELECT ... FOR UPDATE` на `feature_params` при отключении.
* В логике auto-disable важно писать в `audit_log`, чтобы в UI отображалась причина отключения.

---

Отлично — спасибо за уточнение. Ниже — дополнение и уточнённый, готовый к реализации план (включая SQL-проверки, порядок действий, транзакционность, обработку гонок и примеры кода), где **auto-disable применяется только к фичам, помеченным тегом `auto-disable`** (или с другим slug, задаваемым в конфиге / project_settings).

---

Ниже — полная логика и рекомендации по реализации.

---

## 1) Где проверяем тег (SQL)

Проверка наличия тега:

```sql
-- Есть таблицы feature_tags(feature_id, tag_id) и tags(id, slug)
SELECT 1
FROM feature_tags ft
JOIN tags t ON t.id = ft.tag_id
WHERE ft.feature_id = $1
  AND t.slug = 'auto-disable'
LIMIT 1;
```
---

## 2) Поток обработки `POST /sdk/v1/features/{feature_key}/report-error`

### Основные шаги (в порядке)

1. **Аутентификация / авторизация** SDK (ApiKeyAuth): определить `project_id` и `environment_id` по ключу.
2. **Найти фичу** по `feature_key` и `project_id`. Если не найден — вернуть 404.
3. **Вставить запись в `monitoring.error_reports`** (hypertable). Это всегда делаем (если feature отсутствует — 404).
4. **Проверить, помечена ли фича тегом автодизейбла**. Если нет — закончить: вернуть 200/202 с health.
5. **Получить project_settings**: `auto_disable_enabled`, `auto_disable_requires_approval`, `auto_disable_error_threshold`, `auto_disable_time_window_sec`.
6. **Если auto_disable_enabled == false** → finish.
7. **Подсчитать ошибки** за окно (N секунд): `CountRecent(feature_id, env_id, window)`.
8. **Если count >= threshold**:

    * Запустить атомарную процедуру (транзакция):

        * Блокировка `feature_params` (SELECT ... FOR UPDATE) для данного (feature_id, env_id).
        * Перепроверка состояния `enabled` (если уже false — ничего не делать).
        * Если `requires_approval = false` → `UPDATE feature_params SET enabled = false, updated_at = now()`; вставить `audit_log` запись `action = 'auto_disable'`.
        * Если `requires_approval = true` → `CREATE pending_changes` + `pending_change_entities` (entity='feature', entity_id = feature_id), статус = `pending`.
    * После коммита — **отправить уведомление** по WS (broadcaster) для `project_id + environment_id` (event: `feature_auto_disabled` или `pending_change_created`).
9. Вернуть `200` с актуальным `FeatureHealth` или `202` если создан pending change (и требуется approve).

---

## 3) Транзакционность, блокировки и idempotency

* Для операций, изменяющих `feature_params` или создающих `pending_changes`, используем **одну транзакцию**.
* Чтобы избежать гонок (несколько SDK одновременно превысили порог), перед модификацией делаем:

  ```sql
  SELECT enabled
  FROM feature_params
  WHERE feature_id = $1 AND environment_id = $2
  FOR UPDATE;
  ```

  Это гарантирует сериализацию попыток отключить одну и ту же фичу в одном окружении.
* После `SELECT FOR UPDATE` ещё раз проверяем `enabled` и `requires_approval` и только потом делаем UPDATE или CREATE pending_change.
* **Idempotency:** если фича уже отключена — ничего не делаем и возвращаем текущий статус (200).
* Создание `pending_changes` должно быть с `ON CONFLICT DO NOTHING` если используете уникальные ключи; либо проверяйте, что уже нет активного pending change по этой сущности (SELECT ... WHERE status IN ('pending', 'approved'?) — в зависимости от логики).

---

## 4) Примеры SQL-подходов

**Вставка error_report (пример):**

```sql
INSERT INTO monitoring.error_reports
  (event_id, project_id, feature_id, environment_id, error_type, error_message, context, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now());
```

**Подсчёт ошибок (CountRecent):**

```sql
SELECT COUNT(*) FROM monitoring.error_reports
WHERE feature_id = $1
  AND environment_id = $2
  AND created_at > now() - ($3 || ' seconds')::interval;
```

**SELECT FOR UPDATE и авто-отключение (в транзакции):**

```sql
BEGIN;

-- блокируем строку feature_params
SELECT enabled
FROM feature_params
WHERE feature_id = $1 AND environment_id = $2
FOR UPDATE;

-- перепроверка enabled
-- if enabled then
UPDATE feature_params
SET enabled = false, updated_at = now()
WHERE feature_id = $1 AND environment_id = $2;

-- вставляем в audit_log
INSERT INTO audit_log (feature_id, actor, action, old_value, new_value, created_at, entity, project_id, environment_id)
VALUES ($feature_id, 'system', 'auto_disable', ...);

COMMIT;
```

**Создание pending change (если requires_approval):**

```sql
INSERT INTO pending_changes
(project_id, requested_by, request_user_id, change, status, created_at, environment_id)
VALUES ($project_id, 'system', NULL, $json_change, 'pending', now(), $environment_id)
RETURNING id;
-- затем insert into pending_change_entities (pending_change_id, entity, entity_id)
```

`$json_change` — сериализация нового состояния: например `{ "feature_params": { "feature_id": "...", "enabled": false } }` или тому подобное по контракту.

---

## 5) Где логика — слой UseCase / Repo

**Репозитории:**

* `ErrorReportRepository`:

    * `Insert(ctx, report *ErrorReport) error`
    * `CountRecent(ctx, featureID, envID, windowSec int) (int, error)`
* `FeatureRepository`:

    * `FindByKey(projectID, featureKey) -> Feature`
    * `GetFeatureParamsForUpdate(featureID, envID) -> feature_params row` (SELECT FOR UPDATE)
    * `UpdateFeatureParamsEnabled(featureID, envID, enabled bool)`
    * `IsTagged(featureID, tagSlug) -> bool` (or `GetTags(featureID)`)
* `PendingChangesRepository`:

    * `CreatePendingChange(...)`
* `AuditRepository`:

    * `Write(...)`
* `ProjectSettingsRepository`:

    * `GetBool(projectID, 'auto_disable_requires_approval')` and thresholds

**UseCase `AutoDisableUseCase` (pseudo):**

```go
func (uc *AutoDisableUseCase) ReportError(ctx, sdkCtx, featureKey, envKey, report) (FeatureHealth, statusCode, error) {
  // 1. Resolve projectID, environmentID from sdk api key / envKey
  feature := featureRepo.FindByKey(projectID, featureKey)
  if feature == nil { return 404 }

  // 2. Insert error report
  errRepo.Insert(report)

  // 3. Check tag
  tagSlug := "auto-disable"
  if !featureRepo.IsTagged(feature.ID, tagSlug) {
    // return health (no auto-disable)
  }

  // 4. Get settings
  if !projectSettings.GetBool(projectID, "auto_disable_enabled", true) {
    // return health
  }
  threshold := projectSettings.GetInt(projectID, "auto_disable_error_threshold", 20)
  window := projectSettings.GetInt(projectID, "auto_disable_time_window_sec", 60)
  requiresApproval := projectSettings.GetBool(projectID, "auto_disable_requires_approval", false)

  // 5. Count recent
  cnt := errRepo.CountRecent(feature.ID, envID, window)
  if cnt < threshold { return health }

  // 6. Do atomic change
  tx := db.BeginTx(...)
  defer tx.RollbackOnError()
  fp := featureRepo.GetFeatureParamsForUpdate(tx, feature.ID, envID)
  if !fp.Enabled {
     tx.Commit()
     return health // already disabled
  }
  if requiresApproval {
     pcRepo.CreatePendingChange(tx, ...)
     tx.Commit()
     // notify via broadcaster about pending change
     return health with 202
  } else {
     featureRepo.UpdateFeatureParamsEnabled(tx, feature.ID, envID, false)
     auditRepo.InsertAutoDisable(tx, feature, envID, details)
     tx.Commit()
     // notify via broadcaster about auto-disable
     return updated health with 200
  }
}
```

---

## 6) WS / Audit / Pending Changes integration

* После успешной транзакции:

    * Для immediate-disable: вставить в `audit_log(action='auto_disable')` и `broadcaster.Broadcast(projectID, envID, payload)` с payload `{ "event": "feature_auto_disabled", "feature_id": ... }`
    * Для pending-change: `broadcaster.Broadcast(..., { "event":"pending_change_created", "pending_change_id": ... })`

Frontend будет отображать это событие и обновлять UI.

---

## 7) Проверки и валидации входных данных

* `error_type` — non-empty, max length (e.g. 100).
* `error_message` — optional, max length (e.g. 2000).
* `context` — must be JSON; sanitize / limit size.
* Rate limit per sdk api key (to avoid DoS), e.g. 1000 reports/min per key.

---

## 8) Индексы и оптимизация

* `monitoring.error_reports` уже индексирован по `(feature_id, environment_id, created_at DESC)`.
* Добавить/проверить индекс на `feature_tags (feature_id, tag_id)` и `tags (slug)` (unique).
* Индекс на `project_settings (project_id, name)` есть (unique), быстрый lookup.
* При большом потоке ошибок — агрегировать CountRecent через Timescale continuous aggregate (в будущем) для speed.

---

## 9) Метрики / Telemetry

* `auto_disable_checks_total` (counter)
* `auto_disable_tripped_total` (counter)
* `auto_disable_pending_total` (counter)
* `error_reports_received_total` (counter)
* Высылать Prometheus metrics при каждом значимом шаге.

---

## 10) Логирование и мониторинг

* Логировать попытки auto-disable с контекстом (feature, env, project, count, threshold).
* Alert (optional): если auto-disable срабатывает > X раз за час — оповещение SRE/owners.

---

## 11) Тесты (unit + integration)

**Unit:**

* вставка error_report, отсутствие тега → нет disable.
* вставка error_report + тег + ниже threshold → нет disable.
* при достижении threshold, `requires_approval=false` → feature_params.enabled=false и audit записан.
* при `requires_approval=true` → создан pending_change, feature остаётся enabled.

**Integration:**

* End-to-end через тестовую DB (Timescale extension enabled) — проверка транзакционной атомарности.

---

## 12) Edge-cases и замечания

* Если SDK шлёт много мелких events (spam), уникальный индекс не поможет (мы решили не ставить uniq constraint на event_id). Rate-limiter обязателен.
* Для error_rate (доля ошибок) нужна метрика total requests. Пока threshold-by-count — прост и эффектив.
* Если проект хочет отключить auto-disable для конкретной фичи — можно добавлять противотег или project_settings override list (advanced).
* Поддерживать override имени тега в `project_settings` (ключ `auto_disable_tag`) удобно.

---

## 13) Примеры возвращаемых ответов

**200 — feature disabled immediately**

```json
{
  "feature_key": "new_checkout",
  "environment_key": "prod",
  "enabled": false,
  "auto_disabled": true,
  "error_count_last_window": 25,
  "threshold": 20
}
```

**202 — pending change created**
