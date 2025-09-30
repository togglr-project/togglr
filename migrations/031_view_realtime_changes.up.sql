DROP VIEW IF EXISTS v_realtime_events CASCADE;
CREATE OR REPLACE VIEW v_realtime_events AS
(
    SELECT
        'pending'::text AS source,
        pc.id::text     AS event_id,
        pc.project_id,
        pc.environment_id,
        e.key           AS environment_key,
        pce.entity,
        pce.entity_id,
        pc.status       AS action,
        pc.created_at
    FROM pending_changes pc
    JOIN environments e ON e.id = pc.environment_id
    JOIN pending_change_entities pce ON pce.pending_change_id = pc.id
    WHERE pc.created_at > now() - interval '1 hour'
)
UNION ALL
(
    SELECT
        'audit'::text AS source,
        grouped.event_id,
        grouped.project_id,
        grouped.environment_id,
        e.key         AS environment_key,
        'feature'     AS entity,
        grouped.feature_id AS entity_id,
        grouped.action,
        grouped.created_at
    FROM (
             SELECT
                 request_id::text AS event_id,
                 project_id,
                 environment_id,
                 feature_id,
                 max(action)     AS action,
                 max(created_at) AS created_at
             FROM audit_log
             WHERE created_at > now() - interval '1 hour'
             GROUP BY project_id, environment_id, request_id, feature_id
         ) grouped
         JOIN environments e ON e.id = grouped.environment_id
);
