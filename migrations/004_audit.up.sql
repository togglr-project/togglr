create table audit_log (
    id bigserial primary key,
    feature_id uuid not null references features(id) on delete no action,
    actor varchar(50) not null,         -- user/system
    action varchar(50) not null,        -- "create", "update", "delete"
    old_value jsonb,
    new_value jsonb,
    created_at timestamptz not null default now()
);

create index if not exists idx_audit_log_feature_id on audit_log(feature_id);
create index if not exists idx_audit_log_created_at on audit_log(created_at);
