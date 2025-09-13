create table audit_log (
    id bigserial primary key,
    feature_id uuid not null references features(id) on delete cascade,
    actor text not null,         -- user/system
    action text not null,        -- "create", "update", "delete"
    old_value jsonb,
    new_value jsonb,
    created_at timestamptz not null default now()
);
