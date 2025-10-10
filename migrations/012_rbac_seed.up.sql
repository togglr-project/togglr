-- Seed default roles
insert into roles (key, name, description) values
    ('project_owner',   'Project Owner',   'Full control of project'),
    ('project_manager', 'Project Manager', 'Manage features and rules'),
    ('project_member',  'Project Member',  'Toggle features'),
    ('project_viewer',  'Project Viewer',  'Read-only')
on conflict (key) do nothing;

-- Seed permissions
insert into permissions (key, name) values
    ('project.view',       'View project'),
    ('project.manage',     'Manage project'),
    ('feature.view',       'View features'),
    ('feature.toggle',     'Toggle features'),
    ('feature.manage',     'Manage features'),
    ('rule.manage',        'Manage rules'),
    ('audit.view',         'View audit'),
    ('membership.manage',  'Manage memberships')
on conflict (key) do nothing;

-- Grant permissions to roles
with r as (select id, key from roles), p as (select id, key from permissions)
insert into role_permissions (role_id, permission_id)
select r.id, p.id
from r
join p on (
    (r.key = 'project_owner') or
    (r.key = 'project_manager' and p.key in ('project.view','project.manage','feature.view','feature.toggle','feature.manage','rule.manage','audit.view')) or
    (r.key = 'project_member' and p.key in ('feature.view','feature.toggle','project.view')) or
    (r.key = 'project_viewer' and p.key in ('project.view','feature.view'))
)
on conflict do nothing;

-- Add global permission key for creating projects
insert into permissions (key, name)
values ('project.create', 'Create projects')
on conflict (key) do nothing;
