create table audit_log (
    id bigserial primary key,
    flag_id uuid not null references flags(id) on delete cascade,
    actor text not null,         -- user/system
    action text not null,        -- "create", "update", "delete"
    old_value jsonb,
    new_value jsonb,
    created_at timestamptz not null default now()
);
