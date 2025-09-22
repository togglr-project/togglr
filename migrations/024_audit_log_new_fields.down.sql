ALTER TABLE audit_log
    DROP COLUMN IF EXISTS entity_id,
    DROP COLUMN IF EXISTS username;
