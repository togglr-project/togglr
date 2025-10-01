CREATE OR REPLACE VIEW v_feature_health AS
SELECT
    f.project_id,
    f.id                AS feature_id,
    f.key               AS feature_key,
    e.id                AS environment_id,
    e.key               AS environment_key,
    fp.enabled,
    -- количество ошибок за последние 10 минут
    COUNT(er.event_id) FILTER (
        WHERE er.created_at > now() - interval '10 minutes'
        )::int AS error_count_last_10m,
    -- последнее время ошибки
    MAX(er.created_at) FILTER (
        WHERE er.created_at > now() - interval '10 minutes'
        ) AS last_error_at,
    -- условный error_rate: ошибки / 600 секунд
    (COUNT(er.event_id) FILTER (
        WHERE er.created_at > now() - interval '10 minutes'
        )::float / NULLIF(600, 0)) AS error_rate,
    -- настройка: автоотключение требует ли approval
    COALESCE(
            (ps.value->>0)::boolean,
            false
    ) AS requires_approval
FROM features f
         JOIN feature_params fp ON fp.feature_id = f.id
         JOIN environments e ON e.id = fp.environment_id
         LEFT JOIN monitoring.error_reports er
                   ON er.feature_id = f.id AND er.environment_id = e.id
         LEFT JOIN project_settings ps
                   ON ps.project_id = f.project_id AND ps.name = 'auto_disable_requires_approval'
GROUP BY f.project_id, f.id, f.key, e.id, e.key, fp.enabled, ps.value;
