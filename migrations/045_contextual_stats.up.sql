CREATE TABLE monitoring.feature_contextual_stats (
    project_id      uuid                                   NOT NULL
        REFERENCES public.projects ON DELETE CASCADE,
    feature_id      uuid                                   NOT NULL
        REFERENCES public.features ON DELETE CASCADE,
    environment_id  bigint                                 NOT NULL
        REFERENCES public.environments ON DELETE CASCADE,
    algorithm_slug  varchar(100)                           NOT NULL
        REFERENCES public.algorithms ON DELETE CASCADE,
    variant_key     varchar(100)                           NOT NULL,
    feature_key     varchar(50)                            NOT NULL,
    environment_key varchar(20)                            NOT NULL,
    feature_dim     integer                  DEFAULT 32    NOT NULL,
    matrix_a        jsonb                    DEFAULT '[]'  NOT NULL,
    vector_b        jsonb                    DEFAULT '[]'  NOT NULL,
    pulls           bigint                   DEFAULT 0     NOT NULL,
    total_reward    numeric(24, 6)           DEFAULT 0     NOT NULL,
    successes       bigint                   DEFAULT 0     NOT NULL,
    failures        bigint                   DEFAULT 0     NOT NULL,
    updated_at      timestamp with time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, environment_id, algorithm_slug, variant_key)
);

CREATE INDEX idx_contextual_stats_project ON monitoring.feature_contextual_stats(project_id);
CREATE INDEX idx_contextual_stats_feature ON monitoring.feature_contextual_stats(feature_key, environment_key);
