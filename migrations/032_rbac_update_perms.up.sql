DELETE FROM permissions WHERE key IN ('rule.manage', 'variant.manage');

INSERT INTO permissions (key, name) VALUES
('segment.manage', 'Manage segments'),
('schedule.manage', 'Manage schedules'),
('category.manage', 'Manage categories'),
('tag.manage', 'Manage project tags');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.key IN ('segment.manage', 'schedule.manage', 'category.manage', 'tag.manage')
WHERE r.key IN ('project_owner', 'project_manager')
ON CONFLICT DO NOTHING;

-- project_member can only toggle/view/manage features,
-- so segment/schedule are not added to him

CREATE OR REPLACE VIEW v_role_permissions AS
SELECT
    r.id          AS role_id,
    r.key         AS role_key,
    r.name        AS role_name,
    json_agg(
            json_build_object(
                    'key', p.key,
                    'name', p.name
            ) ORDER BY p.key
    ) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
GROUP BY r.id, r.key, r.name
ORDER BY r.key;

CREATE OR REPLACE VIEW v_user_project_permissions AS
SELECT
    m.user_id,
    m.project_id,
    r.key AS role_key,
    json_agg(
            json_build_object(
                    'key', p.key,
                    'name', p.name
            ) ORDER BY p.key
    ) AS permissions
FROM memberships m
JOIN roles r ON r.id = m.role_id
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
GROUP BY m.user_id, m.project_id, r.key
ORDER BY m.user_id, m.project_id;
