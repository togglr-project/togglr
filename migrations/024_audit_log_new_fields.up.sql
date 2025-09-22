ALTER TABLE audit_log
    ADD COLUMN entity_id uuid,
    ADD COLUMN username varchar(255);

CREATE INDEX idx_audit_log_entity_id
    ON audit_log (entity_id);

CREATE INDEX idx_audit_log_username
    ON audit_log (username);
