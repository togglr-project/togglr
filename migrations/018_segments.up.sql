create table segments
(
    id          uuid primary key default gen_random_uuid(),
    project_id  uuid not null references projects on delete cascade,
    name        varchar(255) not null,
    description varchar(255),
    conditions  jsonb not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    unique (project_id, name)
);

create index if not exists idx_segments_project_id on segments(project_id);

alter table rules
    add column segment_id uuid references segments on delete set null,
    add column is_customized boolean not null default false;

create trigger trg_segments_set_updated_at
    before update on segments
    for each row execute function set_updated_at();
