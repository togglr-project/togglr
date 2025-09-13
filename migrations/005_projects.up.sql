create table if not exists projects (
    id uuid primary key default gen_random_uuid(),
    name text not null unique,
    description text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    archived_at timestamptz
);

alter table features add column project_id uuid not null references projects(id) on delete cascade;
