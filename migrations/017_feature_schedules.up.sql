create table feature_schedules
(
    id          uuid primary key default gen_random_uuid(),
    project_id  uuid not null references projects(id) on delete cascade,
    feature_id  uuid not null references features(id) on delete cascade,
    starts_at   timestamp with time zone,
    ends_at     timestamp with time zone,
    cron_expr   varchar(255), -- cron-like (optional, for example, "0 9 * * MON-FRI")
    timezone    varchar(50) default 'UTC' not null,
    action      varchar(15) not null check (action in ('enable', 'disable')),
    created_at  timestamp with time zone default now() not null,
    updated_at  timestamp with time zone default now() not null
);

create trigger trg_feature_schedules_set_updated_at
    before update on feature_schedules
    for each row execute function set_updated_at();

create index if not exists idx_feature_schedules_project_id on feature_schedules(project_id);
create index if not exists idx_feature_schedules_feature_id on feature_schedules(feature_id);
create index if not exists idx_feature_schedules_time on feature_schedules(starts_at, ends_at);

alter table feature_schedules
    add constraint feature_schedules_time_range
        check (starts_at is null or ends_at is null or starts_at < ends_at) not valid;
