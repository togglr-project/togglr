## 📦 Основные сущности

### 1. **features**

Главная таблица — фичефлаги.

```sql
create table features (
    id uuid primary key default gen_random_uuid(),
    key text not null unique,         -- machine name, e.g. "new_ui"
    name text not null,               -- human readable name
    description text,                 -- optional description
    kind text not null,               -- "simple" | "multivariant"
    default_variant text not null,    -- "on"/"off" for simple, or variant name
    enabled boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    project_id uuid not null references projects(id) on delete cascade
);
```

---

### 2. **variants**

Для multivariant флагов (A/B/n-тесты).

```sql
create table variants (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    name text not null,            -- e.g. "A", "B"
    rollout_percent int not null,  -- % of traffic (0..100)
    constraint variants_unique unique(feature_id, name)
);
```

👉 Для boolean-флагов в `variants` можно ничего не хранить.

---

### 3. **rules**

Условия таргетинга (по пользователю, региону, плану и т.д.).
Сохраняем в виде **JSON Logic** или простого JSON-DSL.

```sql
create table rules (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    condition jsonb not null,      -- e.g. {"attribute":"country","op":"=","value":"RU"}
    flag_variant_id uuid not null references flag_variants(id) on delete cascade,
    priority int not null default 0,
    created_at timestamptz not null default now(),
    
    constraint rules_flag_variants_unique unique(feature_id, flag_variant_id, condition)
);
```

👉 Механика:

* SDK берёт флаг
* проверяет правила (по приоритету)
* если условие подходит → возвращает variant
* иначе — default\_variant

---

### 4. **audit\_log**

История изменений (кто, что и когда поменял).

```sql
create table audit_log (
    id bigserial primary key,
    feature_id uuid not null references features(id) on delete cascade,
    actor text not null,         -- user/system
    action text not null,        -- "create", "update", "delete"
    old_value jsonb,
    new_value jsonb,
    created_at timestamptz not null default now()
);
```

---

## 📊 Пример данных

### features

| id | key      | name        | type         | default\_variant | enabled |
| -- | -------- | ----------- | ------------ | ---------------- | ------- |
| 1  | new\_ui  | New UI flag | boolean      | off              | true    |
| 2  | checkout | Checkout AB | multivariant | A                | true    |

### variants

| id | feature\_id | name | rollout\_percent |
| -- | -------- | ---- | ---------------- |
| 1  | 2        | A    | 50               |
| 2  | 2        | B    | 50               |

### rules

| id | feature\_id | condition                                     | variant | rollout\_percent |
| -- | -------- | --------------------------------------------- | ------- | ---------------- |
| 1  | 1        | {"attribute":"country","op":"=","value":"RU"} | on      | 100              |
| 2  | 2        | {"attribute":"plan","op":"=","value":"pro"}   | B       | 100              |

---

Таким образом у нас получается гибкая модель:

* **boolean flags** → простые on/off
* **multivariant flags** → A/B/n тесты с процентами
* **rules** → таргетинг и прогрессивные раскатки
* **audit\_log** → история изменений

