CREATE TABLE environments (
    id BIGSERIAL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key VARCHAR(20) NOT NULL, -- dev, stage, prod
    name VARCHAR(50) NOT NULL,
    api_key UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (project_id, key),
    UNIQUE (api_key)
);

ALTER TABLE projects DROP COLUMN api_key; -- in environments for a project now.

CREATE TABLE feature_params (
    feature_id UUID NOT NULL REFERENCES features(id) ON DELETE CASCADE,
    environment_id BIGINT NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT false,
    default_value VARCHAR(128) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY (feature_id, environment_id)
);

CREATE OR REPLACE FUNCTION update_feature_params_modified()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_feature_params_updated
    BEFORE UPDATE ON feature_params
    FOR EACH ROW
EXECUTE FUNCTION update_feature_params_modified();

ALTER TABLE features DROP COLUMN enabled;
ALTER TABLE features DROP COLUMN default_variant;

ALTER TABLE rules ADD COLUMN environment_id BIGINT NOT NULL default 0;

ALTER TABLE rules
    ADD CONSTRAINT rules_environment_id_fkey
        FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;

ALTER TABLE rules
    DROP CONSTRAINT IF EXISTS rules_unique_priority;

ALTER TABLE rules
    ADD CONSTRAINT rules_unique_priority_per_env
        UNIQUE (feature_id, environment_id, priority);

ALTER TABLE flag_variants ADD COLUMN environment_id BIGINT NOT NULL DEFAULT 0;

ALTER TABLE flag_variants
    ADD CONSTRAINT flag_variants_environment_id_fkey
        FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;

ALTER TABLE flag_variants
    DROP CONSTRAINT flag_variants_unique;

ALTER TABLE flag_variants
    ADD CONSTRAINT flag_variants_unique
        UNIQUE (feature_id, environment_id, name);

ALTER TABLE feature_schedules ADD COLUMN environment_id BIGINT NOT NULL DEFAULT 0;

ALTER TABLE feature_schedules
    ADD CONSTRAINT feature_schedules_environment_id_fkey
        FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;

DROP INDEX IF EXISTS idx_feature_schedules_feature_id;
CREATE INDEX idx_feature_schedules_feature_id ON feature_schedules(feature_id, environment_id);

ALTER TABLE feature_schedules
    DROP CONSTRAINT feature_schedules_no_overlap_guard;

ALTER TABLE feature_schedules
    ADD CONSTRAINT feature_schedules_no_overlap_guard
        EXCLUDE USING gist (
        feature_id WITH =,
        environment_id WITH =,
        tstzrange(starts_at, ends_at, '[]'::text) WITH &&
    ) WHERE ((cron_expr IS NULL));

-- ALTER TABLE segments ADD COLUMN environment_id BIGINT;
-- ALTER TABLE segments ALTER COLUMN environment_id SET NOT NULL;
-- ALTER TABLE segments
--     ADD CONSTRAINT segments_environment_id_fkey
--         FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;
-- ALTER TABLE segments
--     DROP CONSTRAINT segments_project_id_name_key;
-- ALTER TABLE segments
--     ADD CONSTRAINT segments_project_env_name_key
--         UNIQUE (project_id, environment_id, name);

ALTER TABLE audit_log ADD COLUMN environment_id BIGINT;
ALTER TABLE audit_log ALTER COLUMN environment_id SET NOT NULL;

ALTER TABLE audit_log
    ADD CONSTRAINT audit_log_environment_id_fkey
        FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;

ALTER TABLE pending_changes ADD COLUMN environment_id BIGINT NOT NULL DEFAULT 0;

ALTER TABLE pending_changes
    ADD CONSTRAINT pending_changes_environment_id_fkey
        FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE;

CREATE OR REPLACE VIEW v_features_full AS
SELECT f.id,
       f.project_id,
       fp.environment_id,
       e.key as environment_key,
       f.key,
       f.kind,
       f.rollout_key,
       fp.enabled,
       fp.default_value,
       f.name,
       f.description,
       f.created_at,
       f.updated_at
FROM features f
         JOIN feature_params fp ON fp.feature_id = f.id
         JOIN environments e ON e.id = fp.environment_id;

CREATE OR REPLACE VIEW v_projects_full AS
SELECT p.id,
       p.name,
       e.key as environment_key,
       e.api_key,
       p.created_at
FROM projects p
         JOIN environments e ON e.project_id = p.id;

CREATE OR REPLACE FUNCTION create_default_environments()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO environments (project_id, key, name, api_key)
    VALUES
        (NEW.id, 'prod',  'Production', gen_random_uuid()),
        (NEW.id, 'stage', 'Staging',    gen_random_uuid()),
        (NEW.id, 'dev',   'Development',gen_random_uuid());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_create_envs_on_project_insert ON projects;

CREATE TRIGGER trg_create_envs_on_project_insert
    AFTER INSERT ON projects
    FOR EACH ROW
EXECUTE FUNCTION create_default_environments();

CREATE OR REPLACE FUNCTION create_default_feature_params()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO feature_params (feature_id, environment_id, enabled, default_value)
    SELECT
        NEW.id,
        e.id,
        false,
        ''
    FROM environments e
    WHERE e.project_id = NEW.project_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_create_params_on_feature_insert ON features;

CREATE TRIGGER trg_create_params_on_feature_insert
    AFTER INSERT ON features
    FOR EACH ROW
EXECUTE FUNCTION create_default_feature_params();

---

CREATE OR REPLACE FUNCTION create_feature_params_for_env()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO feature_params (feature_id, environment_id, enabled, default_value)
    SELECT f.id, NEW.id, false, ''
    FROM features f
    WHERE f.project_id = NEW.project_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_create_feature_params_on_env_insert ON environments;

CREATE TRIGGER trg_create_feature_params_on_env_insert
    AFTER INSERT ON environments
    FOR EACH ROW
EXECUTE FUNCTION create_feature_params_for_env();
