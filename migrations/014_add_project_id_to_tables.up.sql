alter table rules add column project_id uuid not null references projects(id) on delete cascade;
alter table flag_variants add column project_id uuid not null references projects(id) on delete cascade;
alter table audit_log add column project_id uuid not null references projects(id) on delete cascade;
