CREATE TABLE monitoring.feature_optimizer_stats (
    project_id      uuid                                   NOT NULL
        REFERENCES public.projects ON DELETE CASCADE,
    feature_id      uuid                                   NOT NULL
        REFERENCES public.features ON DELETE CASCADE,
    environment_id  bigint                                 NOT NULL
        REFERENCES public.environments ON DELETE CASCADE,
    algorithm_slug  varchar(100)                           NOT NULL
        REFERENCES public.algorithms ON DELETE CASCADE,
    feature_key     varchar(50)                            NOT NULL,
    environment_key varchar(20)                            NOT NULL,
    iteration       bigint                   DEFAULT 0     NOT NULL,
    current_value   numeric(24, 6)           DEFAULT 0     NOT NULL,
    best_value      numeric(24, 6)           DEFAULT 0     NOT NULL,
    best_reward     numeric(24, 6)           DEFAULT 0     NOT NULL,
    metric_sum      numeric(24, 6)           DEFAULT 0     NOT NULL,
    last_error      numeric(24, 6)           DEFAULT 0     NOT NULL,
    integral        numeric(24, 6)           DEFAULT 0     NOT NULL,
    step_size       numeric(24, 6)           DEFAULT 0.1   NOT NULL,
    temperature     numeric(24, 6)           DEFAULT 1.0   NOT NULL,
    updated_at      timestamp with time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, environment_id, algorithm_slug)
);

CREATE INDEX idx_optimizer_stats_project ON monitoring.feature_optimizer_stats(project_id);
CREATE INDEX idx_optimizer_stats_feature ON monitoring.feature_optimizer_stats(feature_key, environment_key);
