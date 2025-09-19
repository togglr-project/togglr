# Анализ схемы БД и рекомендации (по директории migrations)

Ниже — актуальный аудит схемы после пересборки миграций и перечень улучшений: где стоит отказаться от небезопасных/неограниченных `text` полей, какие индексы добавить, какие ограничения (constraints) и что ещё можно улучшить. В конце перечислены изменения, которые добавлены отдельной аддитивной миграцией (021_*), и предложения для следующих итераций.

## Обзор ключевых таблиц (кратко)

- features: ключевая сущность фич-флагов. Поля: id (uuid), key (varchar), name/description (varchar), kind (varchar), default_variant (varchar), enabled (bool), project_id (uuid FK), rollout_key (varchar), created_at/updated_at.
- flag_variants: варианты фичи. Поля: id, feature_id (FK), name (varchar), rollout_percent (int 0..100), project_id (FK), уникальность (feature_id, name).
- rules: правила таргетинга. Поля: id, feature_id (FK), condition (jsonb), flag_variant_id (FK nullable после 015), action (enum rule_action), priority (int), created_at, project_id (FK), segment_id (FK nullable), is_customized (bool).
- audit_log: аудит изменений. Поля: id (bigserial), feature_id (FK), actor (varchar), action (varchar), old_value/new_value (jsonb), entity (varchar), project_id (FK), request_id (uuid), created_at.
- projects: id (uuid), name (varchar unique), description (varchar), api_key (uuid unique), created_at/updated_at/archived_at.
- users: id (serial), username/email (varchar unique), password_hash (varchar), флаги, two_fa_secret (text).
- settings: name (varchar unique), value (jsonb), description (varchar), created_at/updated_at.
- product_info: key (varchar unique), value (text).
- license, license_history: текст лицензии и история.
- RBAC: roles/permissions/role_permissions/memberships(+audit). Есть уникальности, FK, индексы на часто используемые колонки.
- feature_schedules: расписания включений/выключений (action varchar с CHECK, cron_expr varchar, timezone varchar), индексы по времени/ссылкам.
- segments: сегменты (unique (project_id, name), conditions jsonb).
- ldap_sync_*: журналы интеграции, индексы есть.

## Основные проблемы/риски

1) Неограниченные text/varchar поля там, где домен ограничен списком значений:
   - features.kind — фактически два значения: 'simple'|'multivariant'.
   - audit_log.action — по смыслу ограниченный набор ('create','update','delete', возможно 'enable'/'disable').
   - feature_schedules.action — уже ограничено CHECK'ом (ок).

2) Check-ограничения (часть уже есть, часть — валидация NOT VALID):
   - flag_variants.rollout_percent должен быть в диапазоне 0..100 (в наличии).
   - feature_schedules: если заданы оба времени, starts_at < ends_at (в наличии).
   - license*/issued_at/expires_at — корректность временного диапазона (в наличии).

3) Уникальность и правила с NULL:
   - rules допускают NULL в flag_variant_id для include/exclude. Нужна частичная уникальность на (feature_id, action, condition) для action in ('include','exclude').

4) Недостающие индексы для частых выборок:
   - Многие внешние ключи покрыты, но полезны доп. индексы: rules(feature_id, priority desc), audit_log(request_id), segments(project_id), license_history(license_id), feature_schedules(project_id, feature_id).
   - JSONB поля (rules.condition, segments.conditions) — по решению не добавляем GIN индексы на текущем этапе.

5) Формат и длины ключевых полей:
   - features.rollout_key — лучше обеспечить уникальность среди ненулевых значений.

6) Доп. моменты качества данных:
   - membership_audit без явных FK — возможна потеря ссылочной целостности (см. рекомендации ниже).

## Что улучшено в этой итерации (миграция 021_*)

Только безопасные/аддитивные изменения, без GIN индексов:

- Индексы для частых выборок:
  - rules(feature_id, priority desc)
  - audit_log(request_id)
  - segments(project_id)
  - license_history(license_id)
  - feature_schedules(project_id, feature_id)
- Уникальности/ограничения:
  - Частичный UNIQUE для rules: (feature_id, action, condition) где action in ('include','exclude')
  - features.rollout_key — уникальность среди ненулевых значений
  - Мягкая проверка значений audit_log.action через CHECK (NOT VALID)

Файлы миграций:
- migrations/021_indexes_and_constraints.up.sql
- migrations/021_indexes_and_constraints.down.sql

## Что можно улучшить дальше (потребует обсуждения/оценки влияния)

1) Перейти с varchar на enum там, где это уместно:
   - features.kind: создать тип feature_kind и выполнить ALTER COLUMN TYPE USING.
   - audit_log.action: создать тип audit_action (если набор значений стабилен) или оставаться на CHECK.

2) Скоуп уникальности keys:
   - Рассмотреть переход от глобальной уникальности features.key к составной UNIQUE (project_id, key).

3) Валидация формата ключей:
   - Для features.key/roles.key/permissions.key — добавить CHECK на допустимый формат, регистр, длину (или делать на уровне приложения).
   - Для email/username — рассмотреть CITEXT или уникальные индексы по lower(email)/lower(username).

4) FK для membership_audit:
   - Добавить FK: membership_audit.membership_id → memberships(id) ON DELETE SET NULL; actor_user_id → users(id) ON DELETE SET NULL.
   - Добавить индексы на эти колонки и при необходимости на created_at.

5) Валидация JSON (по желанию):
   - CHECK с функцией валидации структуры rules.condition и segments.conditions (без GIN индексов на этом этапе).

6) Масштабирование логов:
   - Партиционирование audit_log/ldap_sync_logs по дате и политика хранения.

7) Безопасность данных:
   - two_fa_secret — хранить шифрованно (pgcrypto) или на уровне приложения.
   - license_text — хранить подпись/хэш; оценить вынос больших текстов.

8) Семантика дефолтов:
   - users.last_login: рассмотреть NULL до первого входа.
   - features.rollout_key: если нужен стаб. ключ для хеширования — сделать NOT NULL DEFAULT gen_random_uuid() + UNIQUE (потребует миграции данных).

## Вывод

Внедрены безопасные улучшения: добавлены полезные индексы, частичная уникальность для правил include/exclude, уникальность rollout_key и мягкая нормализация audit_log.action. Специально не добавляли GIN индексы для JSONB условий по принятому решению. В рекомендациях — более строгие изменения, требующие обсуждения и оценки совместимости.
