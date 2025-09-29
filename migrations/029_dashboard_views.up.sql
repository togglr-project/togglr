DROP VIEW IF EXISTS v_project_pending_summary CASCADE;
CREATE OR REPLACE VIEW v_project_pending_summary AS
SELECT
    pc.project_id,
    pr.name AS project_name,
    e.id AS environment_id,
    e.key AS environment_key,
    COUNT(DISTINCT pc.id) AS total_pending,
    COUNT(DISTINCT CASE WHEN pce.entity = 'feature' THEN pc.id END) AS pending_feature_changes,
    COUNT(DISTINCT CASE WHEN pce.entity = 'feature'
        AND EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = pce.entity_id
              AND t.slug = 'guarded'
        )
                            THEN pc.id END) AS pending_guarded_changes,
    MIN(pc.created_at) AS oldest_request_at
FROM pending_changes pc
         LEFT JOIN pending_change_entities pce ON pce.pending_change_id = pc.id
         JOIN environments e ON e.project_id = pc.project_id
         JOIN projects pr ON pr.id = e.project_id
WHERE pc.status = 'pending'
GROUP BY pc.project_id, project_name, e.id, e.key;

-- ------------------------------------

DROP VIEW IF EXISTS v_project_top_risky_features CASCADE;
CREATE OR REPLACE VIEW v_project_top_risky_features AS
SELECT
    f.project_id,
    pr.name AS project_name,
    fp.environment_id,
    e.key AS environment_key,
    f.id AS feature_id,
    f.name AS feature_name,
    fp.enabled,
    string_agg(DISTINCT t.slug, ', ') AS risky_tags,
    CASE WHEN EXISTS (
        SELECT 1
        FROM pending_change_entities pce
                 JOIN pending_changes pc ON pc.id = pce.pending_change_id
        WHERE pce.entity = 'feature'
          AND pce.entity_id = f.id
          AND pc.status = 'pending'
    ) THEN true ELSE false END AS has_pending
FROM features f
         JOIN feature_params fp ON fp.feature_id = f.id
         JOIN environments e ON e.id = fp.environment_id
         JOIN feature_tags ft ON ft.feature_id = f.id
         JOIN tags t ON ft.tag_id = t.id
         JOIN projects pr ON pr.id = e.project_id
WHERE t.slug IN ('critical','guarded','auto-disable')
GROUP BY f.project_id, project_name, fp.environment_id, e.key, f.id, f.name, fp.enabled;

-- ------------------------------------

DROP VIEW IF EXISTS v_project_recent_activity CASCADE;
CREATE OR REPLACE VIEW v_project_recent_activity AS
WITH audit AS (
    SELECT
        a.project_id,
        e.id AS environment_id,
        e.key AS environment_key,
        p.name AS project_name,
        a.request_id,
        a.username AS actor,
        MIN(a.created_at) AS created_at,
        jsonb_agg(
                jsonb_build_object(
                        'entity', a.entity,
                        'entity_id', a.entity_id,
                        'action', a.action
                )
                ORDER BY a.created_at
        ) AS changes,
        'applied'::text AS status
    FROM audit_log a
             JOIN projects p ON p.id = a.project_id
             JOIN environments e ON e.project_id = p.id
    GROUP BY a.project_id, e.id, e.key, p.name, a.request_id, a.username
),
     pending AS (
         SELECT
             pc.project_id,
             e.id AS environment_id,
             e.key AS environment_key,
             p.name AS project_name,
             pc.id AS request_id,
             pc.requested_by AS actor,
             MIN(pc.created_at) AS created_at,
             jsonb_agg(
                     jsonb_build_object(
                             'entity', pce.entity,
                             'entity_id', pce.entity_id,
                             'action', 'change_request'
                     )
                     ORDER BY pce.created_at
             ) AS changes,
             pc.status
         FROM pending_changes pc
                  JOIN projects p ON p.id = pc.project_id
                  JOIN environments e ON e.project_id = p.id
                  JOIN pending_change_entities pce ON pce.pending_change_id = pc.id
         GROUP BY pc.project_id, e.id, e.key, p.name, pc.id, pc.requested_by, pc.status
     )
SELECT *
FROM (
         SELECT * FROM audit
         UNION ALL
         SELECT * FROM pending
     ) combined
ORDER BY created_at DESC;