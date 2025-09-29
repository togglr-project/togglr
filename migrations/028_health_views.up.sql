DROP VIEW IF EXISTS v_project_health CASCADE;
CREATE OR REPLACE VIEW v_project_health AS
WITH feature_with_category AS (
    SELECT
        f.id AS feature_id,
        f.project_id,
        fp.environment_id,
        e.key AS environment_key,
        fp.enabled,
        EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = f.id
              AND t.slug = 'auto-disable'
        ) AS under_auto_disable,
        EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = f.id
              AND t.slug = 'guarded'
        ) AS is_guarded,
        NOT EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = f.id
        ) AS is_uncategorized
    FROM features f
             JOIN feature_params fp ON fp.feature_id = f.id
             JOIN environments e ON e.id = fp.environment_id
),
     pending AS (
         SELECT DISTINCT pce.entity_id AS feature_id
         FROM pending_change_entities pce
                  JOIN pending_changes pc ON pc.id = pce.pending_change_id
         WHERE pc.status = 'pending'
           AND pce.entity = 'feature'
     )
SELECT
    fwc.project_id,
    pr.name AS project_name,
    fwc.environment_id,
    fwc.environment_key,
    COUNT(DISTINCT fwc.feature_id) AS total_features,
    COUNT(DISTINCT CASE WHEN fwc.enabled THEN fwc.feature_id END) AS enabled_features,
    COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) AS disabled_features,
    COUNT(DISTINCT CASE WHEN fwc.under_auto_disable THEN fwc.feature_id END) AS auto_disable_managed_features,
    COUNT(DISTINCT CASE WHEN fwc.is_uncategorized THEN fwc.feature_id END) AS uncategorized_features,
    COUNT(DISTINCT CASE WHEN fwc.is_guarded THEN fwc.feature_id END) AS guarded_features,
    COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) AS pending_features,
    COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) AS pending_guarded_features,
    CASE
        WHEN (100.0 * COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 20
            OR (100.0 * COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 20
            OR COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) > 1
            THEN 'red'
        WHEN (100.0 * COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 5
            OR (100.0 * COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 5
            OR (100.0 * COUNT(DISTINCT CASE WHEN fwc.under_auto_disable THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 10
            OR COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) = 1
            THEN 'yellow'
        ELSE 'green'
        END AS health_status
FROM feature_with_category fwc
LEFT JOIN pending p ON p.feature_id = fwc.feature_id
JOIN projects pr ON fwc.project_id = pr.id
GROUP BY fwc.project_id, project_name, fwc.environment_id, fwc.environment_key;

-- ------------------------------------

DROP VIEW IF EXISTS v_project_category_health CASCADE;
CREATE OR REPLACE VIEW v_project_category_health AS
WITH feature_with_category AS (
    SELECT
        f.id AS feature_id,
        f.project_id,
        fp.environment_id,
        e.key AS environment_key,
        COALESCE(c.id, '00000000-0000-0000-0000-000000000000'::uuid) AS category_id,
        COALESCE(c.name, 'Uncategorized') AS category_name,
        COALESCE(c.slug, 'uncategorized') AS category_slug,
        fp.enabled,
        EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = f.id
              AND t.slug = 'auto-disable'
        ) AS under_auto_disable,
        EXISTS (
            SELECT 1
            FROM feature_tags ft
                     JOIN tags t ON ft.tag_id = t.id
            WHERE ft.feature_id = f.id
              AND t.slug = 'guarded'
        ) AS is_guarded
    FROM features f
             JOIN feature_params fp ON fp.feature_id = f.id
             JOIN environments e ON e.id = fp.environment_id
             LEFT JOIN feature_tags ft ON ft.feature_id = f.id
             LEFT JOIN tags t ON ft.tag_id = t.id
             LEFT JOIN categories c ON t.category_id = c.id
),
     pending AS (
         SELECT DISTINCT pce.entity_id AS feature_id
         FROM pending_change_entities pce
                  JOIN pending_changes pc ON pc.id = pce.pending_change_id
         WHERE pc.status = 'pending'
           AND pce.entity = 'feature'
     )
SELECT
    fwc.project_id,
    pr.name AS project_name,
    fwc.environment_id,
    fwc.environment_key,
    fwc.category_id,
    fwc.category_name,
    fwc.category_slug,
    COUNT(DISTINCT fwc.feature_id) AS total_features,
    COUNT(DISTINCT CASE WHEN fwc.enabled THEN fwc.feature_id END) AS enabled_features,
    COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) AS disabled_features,
    COUNT(DISTINCT CASE WHEN fwc.under_auto_disable THEN fwc.feature_id END) AS auto_disable_managed_features,
    COUNT(DISTINCT CASE WHEN fwc.is_guarded THEN fwc.feature_id END) AS guarded_features,
    COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) AS pending_features,
    COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) AS pending_guarded_features,
    CASE
        WHEN (100.0 * COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 20
            OR (100.0 * COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 20
            OR COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) > 1
            THEN 'red'
        WHEN (100.0 * COUNT(DISTINCT CASE WHEN NOT fwc.enabled THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 5
            OR (100.0 * COUNT(DISTINCT CASE WHEN p.feature_id IS NOT NULL THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 5
            OR (100.0 * COUNT(DISTINCT CASE WHEN fwc.under_auto_disable THEN fwc.feature_id END) / NULLIF(COUNT(DISTINCT fwc.feature_id),0)) > 10
            OR COUNT(DISTINCT CASE WHEN fwc.is_guarded AND p.feature_id IS NOT NULL THEN fwc.feature_id END) = 1
            THEN 'yellow'
        ELSE 'green'
        END AS health_status
FROM feature_with_category fwc
LEFT JOIN pending p ON p.feature_id = fwc.feature_id
JOIN projects pr ON pr.id = fwc.project_id
GROUP BY fwc.project_id, project_name, fwc.environment_id, fwc.environment_key, fwc.category_id, fwc.category_name, fwc.category_slug;
