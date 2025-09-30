DROP VIEW IF EXISTS v_realtime_events CASCADE;
CREATE OR REPLACE VIEW v_realtime_events AS
(
    SELECT
        'pending'::text AS source,
        pc.id::text AS event_id,
        pc.project_id,
        pc.environment_id,
        e.key AS environment_key,
        pce.entity,
        pce.entity_id,
        pc.status AS action,
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
        al.id::text AS event_id,
        al.project_id,
        al.environment_id,
        e.key AS environment_key,
        al.entity,
        al.entity_id,
        al.action,
        al.created_at
    FROM audit_log al
    JOIN environments e ON e.id = al.environment_id
    WHERE al.created_at > now() - interval '1 hour'
);
