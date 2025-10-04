-- ========================================
-- Pending Changes (Guard Engine)
-- ========================================

-- Main table for pending changes
CREATE TABLE pending_changes (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    requested_by varchar(255) NOT NULL,   -- username or user id
    request_user_id integer,              -- optional FK to users.id if available
    change jsonb NOT NULL,                -- see format: { entities: [...] , meta: {...} }
    status varchar(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, cancelled
    created_at timestamptz DEFAULT now() NOT NULL,
    approved_by varchar(255),
    approved_user_id integer,
    approved_at timestamptz,
    rejected_by varchar(255),
    rejected_at timestamptz,
    rejection_reason text
);

-- Separate entities table for uniqueness constraints
CREATE TABLE pending_change_entities (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    pending_change_id uuid NOT NULL REFERENCES pending_changes(id) ON DELETE CASCADE,
    entity varchar(50) NOT NULL,    -- e.g., 'feature', 'rule', 'feature_schedule'
    entity_id uuid NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL
);

-- Unique constraint: no two pending changes can affect the same entity
-- CREATE UNIQUE INDEX ux_pending_entity_unique ON pending_change_entities (entity, entity_id)
--   WHERE (TRUE); -- we will enforce "only one pending state" via check function/trigger below

-- Indexes for performance
CREATE INDEX idx_pending_changes_project_status ON pending_changes (project_id, status);
CREATE INDEX idx_pending_changes_created_at ON pending_changes (created_at);
CREATE INDEX idx_pending_change_entities_pending_id ON pending_change_entities (pending_change_id);
CREATE INDEX idx_pending_change_entities_entity ON pending_change_entities (entity, entity_id);

-- ========================================
-- Project Approvers
-- ========================================

CREATE TABLE project_approvers (
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role varchar(50) DEFAULT 'approver',
    created_at timestamptz DEFAULT now() NOT NULL,
    PRIMARY KEY (project_id, user_id)
);

CREATE INDEX idx_project_approvers_user_id ON project_approvers (user_id);

-- ========================================
-- Project Settings (key-value for project configuration)
-- ========================================

CREATE TABLE project_settings (
    id serial PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    value jsonb NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    UNIQUE (project_id, name)
);

CREATE INDEX idx_project_settings_project_id ON project_settings (project_id);

create or replace view v_project_effective_settings as
select
    p.id as project_id,
    coalesce(
                    jsonb_object_agg(ps.name, ps.value) filter (where ps.name is not null),
                    '{}'::jsonb
    ) as settings
from projects p
         left join project_settings ps on ps.project_id = p.id
group by p.id;

-- ========================================
-- Functions and Triggers
-- ========================================

-- Function to enforce no duplicate pending for same entity
CREATE OR REPLACE FUNCTION ensure_no_conflicting_pending()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    cnt int;
BEGIN
    -- NEW.pending_change_id exists, check for conflicts in other pending_changes
    SELECT count(*) INTO cnt
    FROM pending_change_entities pce
    JOIN pending_changes pc ON pc.id = pce.pending_change_id
    WHERE pce.entity = NEW.entity
      AND pce.entity_id = NEW.entity_id
      AND pc.status = 'pending'
      AND pce.pending_change_id <> NEW.pending_change_id;

    IF cnt > 0 THEN
        RAISE EXCEPTION 'Entity % % is already locked by another pending change', NEW.entity, NEW.entity_id;
    END IF;

    RETURN NEW;
END;
$$;

-- Trigger to prevent conflicts
CREATE TRIGGER trg_pce_no_conflict
BEFORE INSERT ON pending_change_entities
FOR EACH ROW EXECUTE FUNCTION ensure_no_conflicting_pending();

-- Function to update updated_at for project_settings
CREATE OR REPLACE FUNCTION set_project_settings_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_project_settings_set_updated_at
    BEFORE UPDATE ON project_settings
    FOR EACH ROW EXECUTE FUNCTION set_project_settings_updated_at();

---

-- Function-trigger for adding default project_settings
CREATE OR REPLACE FUNCTION set_default_project_settings()
    RETURNS TRIGGER AS $$
BEGIN
    -- auto-disable approval
    INSERT INTO project_settings (project_id, name, value, created_at, updated_at)
    VALUES (NEW.id, 'auto_disable_requires_approval', 'false'::jsonb, now(), now())
    ON CONFLICT (project_id, name) DO NOTHING;

    -- auto_disable_enabled
    INSERT INTO project_settings (project_id, name, value, created_at, updated_at)
    VALUES (NEW.id, 'auto_disable_enabled', 'true'::jsonb, now(), now())
    ON CONFLICT (project_id, name) DO NOTHING;

    -- auto_disable_error_threshold
    INSERT INTO project_settings (project_id, name, value, created_at, updated_at)
    VALUES (NEW.id, 'auto_disable_error_threshold', '10'::jsonb, now(), now())
    ON CONFLICT (project_id, name) DO NOTHING;

    -- auto_disable_time_window_sec
    INSERT INTO project_settings (project_id, name, value, created_at, updated_at)
    VALUES (NEW.id, 'auto_disable_time_window_sec', '60'::jsonb, now(), now())
    ON CONFLICT (project_id, name) DO NOTHING;

    -- audit log retention
    INSERT INTO project_settings (project_id, name, value, created_at, updated_at)
    VALUES (NEW.id, 'audit_log_retention_days', '180'::jsonb, now(), now())
    ON CONFLICT (project_id, name) DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger on insert project
CREATE TRIGGER trg_set_default_project_settings
    AFTER INSERT ON projects
    FOR EACH ROW
EXECUTE FUNCTION set_default_project_settings();
