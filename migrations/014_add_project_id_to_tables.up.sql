alter table rules add column project_id uuid not null references projects(id) on delete cascade;
alter table flag_variants add column project_id uuid not null references projects(id) on delete cascade;
alter table audit_log add column project_id uuid not null references projects(id) on delete cascade;

create index if not exists idx_rules_project_id on rules(project_id);
create index if not exists idx_flag_variants_project_name on flag_variants (project_id, name);
create index if not exists idx_audit_log_project_id on audit_log(project_id);
