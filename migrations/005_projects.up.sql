create table if not exists projects (
    id uuid primary key default gen_random_uuid(),
    name varchar(128) not null unique,
    description varchar(300),
    api_key uuid not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    archived_at timestamptz
);

alter table features add column project_id uuid not null references projects(id) on delete cascade;
create index if not exists idx_features_project_id on features(project_id);

create trigger trg_projects_set_updated_at
    before update on projects
    for each row execute function set_updated_at();
