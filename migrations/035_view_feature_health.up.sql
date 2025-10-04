CREATE OR REPLACE VIEW v_feature_health
            (project_id, feature_id, feature_key, environment_id, environment_key, enabled,
             error_count_last_10m, last_error_at, error_rate, requires_approval, health_status)
AS
SELECT f.project_id,
       f.id                                                                                      AS feature_id,
       f.key                                                                                     AS feature_key,
       e.id                                                                                      AS environment_id,
       e.key                                                                                     AS environment_key,
       fp.enabled,
       COUNT(er.event_id)
       FILTER (WHERE er.created_at > (now() - '00:10:00'::interval))::integer                AS error_count_last_10m,
       MAX(er.created_at)
       FILTER (WHERE er.created_at > (now() - '00:10:00'::interval))                         AS last_error_at,
       COUNT(er.event_id)
       FILTER (WHERE er.created_at > (now() - '00:10:00'::interval))::double precision
           / NULLIF(600, 0)::double precision                                                    AS error_rate,
       COALESCE((ps.value ->> 0)::boolean, false)                                                AS requires_approval,
       CASE
           WHEN COUNT(er.event_id) FILTER (WHERE er.created_at > (now() - '00:10:00'::interval)) > 100
               THEN 'failing'
           WHEN COUNT(er.event_id) FILTER (WHERE er.created_at > (now() - '00:10:00'::interval)) > 0
               THEN 'degraded'
           ELSE 'ok'
           END                                                                                       AS health_status
FROM features f
         JOIN feature_params fp ON fp.feature_id = f.id
         JOIN environments e ON e.id = fp.environment_id
         LEFT JOIN monitoring.error_reports er ON er.feature_id = f.id AND er.environment_id = e.id
         LEFT JOIN project_settings ps
                   ON ps.project_id = f.project_id AND ps.name::text = 'auto_disable_requires_approval'::text
GROUP BY f.project_id, f.id, f.key, e.id, e.key, fp.enabled, ps.value;
