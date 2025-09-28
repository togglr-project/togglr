## 🟢 Данные для Dashboard

1. **`v_project_health`**
   – Общая сводка по проекту:

   * total\_features, enabled/disabled
   * auto-disable managed
   * guarded count
   * pending count
   * health\_status (green / yellow / red по порогам)

2. **`v_project_category_health`**
   – Сводка по категориям внутри проекта:

   * категории (UI/UX, backend, infra, experiment, critical, guarded…)
   * total\_features / enabled / disabled
   * pending\_guarded\_features
   * health\_status по каждой категории

3. **`v_project_recent_activity`**
   – Лента событий (батчи по `request_id`):

   * кто изменил
   * что изменил (changes JSON)
   * статус (applied / pending)
   * время

4. **`v_project_top_risky_features`**
   – Список фич с критическими тегами (`critical`, `guarded`, `auto-disable`).

   * enabled?
   * pending?
   * какие risky\_tags

5. **`v_project_pending_summary`**
   – Сводка по pending-запросам:

   * total\_pending
   * pending\_feature\_changes
   * pending\_guarded\_changes
   * oldest\_request\_at

---

## 🖼️ ASCII-схема Dashboard

```
+=====================================================================+
|                          DASHBOARD                                  |
+=====================================================================+

[ Project Selector (any\all) ] [Environment Selector, prod by default]

PROJECT OVERVIEW
+----------------------+----------------------+------------------------+
| Project              | Features             | Health Status          |
|----------------------+----------------------+------------------------|
| Shop App             | 12 total / 10 on / 2 off | 🟢 GREEN            |
| Billing Service      | 8 total / 6 on / 2 off   | 🟡 YELLOW           |
| Infra Monitor        | 15 total / 15 on / 0 off | 🔴 RED              |
+----------------------+----------------------+------------------------+

CATEGORY HEALTH
+-------------+-------+--------+---------+----------+----------+
| Category    | Total | Enabled| Disabled| Pending  | Health   |
|-------------+-------+--------+---------+----------+----------|
| UI/UX       | 5     | 4      | 1       | 0        | 🟢 GREEN |
| Backend     | 4     | 4      | 0       | 0        | 🟢 GREEN |
| Critical    | 1     | 1      | 0       | 0        | 🟢 GREEN |
| Guarded     | 2     | 1      | 1       | 1        | 🟡 YELLOW|
| Uncategorized | 0   | 0      | 0       | 0        | 🟢 GREEN |
+-------------+-------+--------+---------+----------+----------+

FEATURE ACTIVITY (top 10)
+------------------------------------------------------+
|  ├ Upcoming (Next State from schedules)              |
|  |   ⏭ new_checkout → ENABLE at 14:30 (5m left)      |
|  |   ⏭ search_bar → DISABLE at 22:00                 |
|  ├ Recent (from audit_log)                           |
|      ⏮ promo_banner → DISABLED by user1 (12:20)      |
|      ⏮ login_flow → RULE UPDATED by user2 (11:40)    |
+------------------------------------------------------+

TOP RISKY FEATURES (top 10)
+----------+-----------------+----------+-------------+-----------+
| Project  | Feature         | Enabled  | Pending     | Risk Tags |
|----------+-----------------+----------+-------------+-----------|
| Shop App | PaymentGateway  | true     | false       | critical  |
| Shop App | FeatureX        | false    | true        | guarded   |
| Infra    | Autoscaler      | true     | false       | auto-disable|
+----------+-----------------+----------+-------------+-----------+

RECENT ACTIVITY (last 5)
+----------+-------------+-----------+-----------+---------------------+
| Project  | Actor       | Status    | When      | Changes             |
|----------+-------------+-----------+-----------+---------------------|
| Shop App | user1       | applied   | 2025-09-24| [{feature:update}]  |
| Shop App | admin       | pending   | 2025-09-23| [{feature:create}]  |
| Billing  | user2       | applied   | 2025-09-22| [{rule:update}]     |
+----------+-------------+-----------+-----------+---------------------+

PENDING SUMMARY
+----------+------------------+--------------------+---------------------+
| Project  | Total Pending    | Guarded Pending    | Oldest Request Age  |
|----------+------------------+--------------------+---------------------|
| Shop App | 2                | 1                  | 2 days              |
| Billing  | 1                | 0                  | 1 hour              |
+----------+------------------+--------------------+---------------------+
```

---

## 📌 Итого

На Dashboard у тебя:

* **сводка проектов** (health + количество фич)
* **категории внутри проекта**
* **Feature Activity** (когда следующие изменения (вычисляется feature processor'ом при помощи метода NextStatte) и прошедшие изменения (из audit_log))
* **лента изменений (audit + pending)**
* **рискованные фичи**
* **сводка pending-запросов**

---

OpenAPI Спецификация (черновик, добавить авторизацию, вынести респонс в отдельную структуру):

```yml
paths:
  /api/v1/dashboard/overview:
    get:
      summary: Project Dashboard overview
      description: |
        Returns aggregated dashboard data for a project:
        - project health
        - category health
        - feature activity (upcoming & recent)
        - recent activity (batched by request_id)
        - risky features
        - pending summary
      parameters:
        - name: environment_key
          in: query
          required: true
          schema:
            type: string
          description: Environment key (prod/stage/dev)
        - name: project_id
          in: query
          required: false
          schema:
            type: string
            format: uuid
          description: Optional project ID to filter results
        - name: limit
          in: query
          required: false
          schema:
            type: integer
            default: 20
          description: Limit for recent activity entries
      responses:
        "200":
          description: Dashboard data
          content:
            application/json:
              schema:
                type: object
                properties:
                  projects:
                    type: array
                    description: Project-level health overview
                    items:
                      $ref: "#/components/schemas/ProjectHealth"
                  categories:
                    type: array
                    description: Per-category health
                    items:
                      $ref: "#/components/schemas/CategoryHealth"
                  feature_activity:
                    type: object
                    description: Feature-level upcoming and recent changes
                    properties:
                      upcoming:
                        type: array
                        items:
                          $ref: "#/components/schemas/FeatureUpcoming"
                      recent:
                        type: array
                        items:
                          $ref: "#/components/schemas/FeatureRecent"
                  recent_activity:
                    type: array
                    description: Recent batched changes
                    items:
                      $ref: "#/components/schemas/RecentActivity"
                  risky_features:
                    type: array
                    description: Features with risky tags (critical, guarded, auto-disable)
                    items:
                      $ref: "#/components/schemas/RiskyFeature"
                  pending_summary:
                    type: array
                    description: Summary of pending changes
                    items:
                      $ref: "#/components/schemas/PendingSummary"

components:
  schemas:
    ProjectHealth:
      type: object
      properties:
        project_id: { type: string, format: uuid }
        total_features: { type: integer }
        enabled_features: { type: integer }
        disabled_features: { type: integer }
        auto_disable_managed_features: { type: integer }
        guarded_features: { type: integer }
        pending_features: { type: integer }
        health_status: { type: string, enum: [green, yellow, red] }

    CategoryHealth:
      type: object
      properties:
        project_id: { type: string, format: uuid }
        category_id: { type: string, format: uuid }
        category_name: { type: string }
        category_slug: { type: string }
        total_features: { type: integer }
        enabled_features: { type: integer }
        disabled_features: { type: integer }
        pending_features: { type: integer }
        guarded_features: { type: integer }
        auto_disable_managed_features: { type: integer }
        health_status: { type: string, enum: [green, yellow, red] }

    FeatureUpcoming:
      type: object
      properties:
        feature_id: { type: string, format: uuid }
        feature_name: { type: string }
        next_state: { type: string, enum: [enabled, disabled] }
        at: { type: string, format: date-time }

    FeatureRecent:
      type: object
      properties:
        feature_id: { type: string, format: uuid }
        feature_name: { type: string }
        action: { type: string }
        actor: { type: string }
        at: { type: string, format: date-time }

    RecentActivity:
      type: object
      properties:
        project_id: { type: string, format: uuid }
        project_name: { type: string }
        request_id: { type: string, format: uuid }
        actor: { type: string }
        created_at: { type: string, format: date-time }
        status: { type: string, enum: [applied, pending, rejected] }
        changes:
          type: array
          items:
            type: object
            properties:
              entity: { type: string }
              entity_id: { type: string, format: uuid }
              action: { type: string }

    RiskyFeature:
      type: object
      properties:
        project_id: { type: string, format: uuid }
        feature_id: { type: string, format: uuid }
        feature_name: { type: string }
        enabled: { type: boolean }
        has_pending: { type: boolean }
        risky_tags: { type: string }

    PendingSummary:
      type: object
      properties:
        project_id: { type: string, format: uuid }
        total_pending: { type: integer }
        pending_feature_changes: { type: integer }
        pending_guarded_changes: { type: integer }
        oldest_request_at: { type: string, format: date-time }
```
