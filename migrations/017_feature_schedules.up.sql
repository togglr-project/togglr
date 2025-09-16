create table feature_schedules
(
    id          uuid primary key default gen_random_uuid(),
    project_id  uuid not null references projects(id) on delete cascade,
    feature_id  uuid not null references features(id) on delete cascade,
    starts_at   timestamp with time zone,
    ends_at     timestamp with time zone,
    cron_expr   text, -- cron-like (optional, for example, "0 9 * * MON-FRI")
    timezone    text default 'UTC' not null,
    action      text not null check (action in ('enable', 'disable')),
    created_at  timestamp with time zone default now() not null
);
