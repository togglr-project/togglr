DELETE FROM permissions WHERE key IN ('rule.manage', 'variant.manage');

INSERT INTO permissions (key, name) VALUES
('segment.manage', 'Manage segments'),
('schedule.manage', 'Manage schedules');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.key IN ('segment.manage', 'schedule.manage')
WHERE r.key IN ('project_owner', 'project_manager')
ON CONFLICT DO NOTHING;

-- project_member может только toggle/view/manage features,
-- поэтому segment/schedule ему НЕ добавляем
