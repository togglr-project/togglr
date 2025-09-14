-- Add global permission key for creating projects
insert into permissions (key, name)
values ('project.create', 'Create projects')
on conflict (key) do nothing;
