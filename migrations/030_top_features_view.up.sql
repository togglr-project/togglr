DROP VIEW IF EXISTS v_top_active_features CASCADE;
CREATE OR REPLACE VIEW v_top_active_features AS
WITH active_schedules AS (
    SELECT DISTINCT ON (fs.feature_id, fs.environment_id)
        fs.feature_id,
        fs.environment_id,
        fs.starts_at,
        fs.ends_at
    FROM feature_schedules fs
    WHERE fs.action = 'enable'
      AND (fs.starts_at IS NULL OR fs.starts_at <= now())
      AND (fs.ends_at IS NULL OR fs.ends_at > now())
    ORDER BY fs.feature_id, fs.environment_id, fs.starts_at NULLS FIRST
),
     scored AS (
         SELECT
             f.id AS feature_id,
             f.project_id,
             fp.environment_id,
             e.key AS environment_key,
             f.name,
             f.description,
             fp.enabled,
             s.starts_at,
             s.ends_at,
             f.updated_at,
             CASE WHEN EXISTS (
                 SELECT 1
                 FROM feature_tags ft
                          JOIN tags t ON ft.tag_id = t.id
                 WHERE ft.feature_id = f.id
                   AND t.slug = 'critical'
             ) THEN true ELSE false END AS is_critical,
             GREATEST(0, 50 - EXTRACT(EPOCH FROM (now() - f.updated_at))/3600) AS recency_score,
             COALESCE((
                          SELECT COUNT(*)
                          FROM audit_log a
                          WHERE a.entity = 'feature'
                            AND a.entity_id = f.id
                            AND a.created_at > now() - interval '7 days'
                      ),0) AS recent_activity_count
         FROM features f
                  JOIN feature_params fp ON fp.feature_id = f.id
                  JOIN environments e ON e.id = fp.environment_id
                  LEFT JOIN active_schedules s
                            ON s.feature_id = f.id
                                AND s.environment_id = fp.environment_id
         WHERE fp.enabled = true
     )
SELECT
    s.*,
    (CASE WHEN s.is_critical THEN 100 ELSE 0 END)
        + s.recency_score
        + (s.recent_activity_count * 10) AS rank_score,
    CASE
        WHEN s.is_critical AND s.recency_score > 0 AND s.recent_activity_count > 0 THEN 'critical + recent update + recent activity'
        WHEN s.is_critical AND s.recency_score > 0 THEN 'critical + recent update'
        WHEN s.is_critical AND s.recent_activity_count > 0 THEN 'critical + recent activity'
        WHEN s.is_critical THEN 'critical'
        WHEN s.recency_score > 0 AND s.recent_activity_count > 0 THEN 'recent update + recent activity'
        WHEN s.recency_score > 0 THEN 'recent update'
        WHEN s.recent_activity_count > 0 THEN 'recent activity'
        ELSE 'default'
        END AS rank_reason
FROM scored s
ORDER BY rank_score DESC;
