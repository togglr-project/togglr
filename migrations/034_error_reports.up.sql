CREATE TABLE monitoring.error_reports
(
    id             BIGSERIAL NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT now() NOT NULL,
    event_id       UUID DEFAULT gen_random_uuid() NOT NULL,

    project_id     UUID NOT NULL REFERENCES public.projects(id) ON DELETE CASCADE,
    feature_id     UUID NOT NULL REFERENCES public.features(id) ON DELETE CASCADE,
    environment_id BIGINT NOT NULL REFERENCES public.environments(id) ON DELETE CASCADE,

    error_type     VARCHAR(100) NOT NULL,
    error_message  TEXT,
    context        JSONB,

    PRIMARY KEY (id, created_at)
);

CREATE INDEX idx_error_reports_feature_env_time
    ON monitoring.error_reports (feature_id, environment_id, created_at DESC);

CREATE INDEX idx_error_reports_type
    ON monitoring.error_reports (error_type);

SELECT create_hypertable('monitoring.error_reports', 'created_at', if_not_exists => TRUE);
SELECT add_retention_policy('monitoring.error_reports', INTERVAL '30 days');
