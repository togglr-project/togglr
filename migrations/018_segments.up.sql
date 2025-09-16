create table segments
(
    id          uuid primary key default gen_random_uuid(),
    project_id  uuid not null references projects on delete cascade,
    name        text not null,
    description text,
    conditions  jsonb not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    unique (project_id, name)
);

alter table rules
    add column segment_id uuid references segments on delete set null,
    add column is_customized boolean not null default false;
