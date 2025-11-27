-- Custom WASM algorithms storage (global)
CREATE TABLE public.custom_algorithms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    kind VARCHAR(50) NOT NULL CHECK (kind IN ('bandit', 'optimizer', 'contextual_bandit')),
    wasm_binary BYTEA NOT NULL,
    wasm_hash VARCHAR(64) NOT NULL,
    default_settings JSONB DEFAULT '{}'::jsonb NOT NULL,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE INDEX idx_custom_algorithms_kind ON public.custom_algorithms(kind);
CREATE INDEX idx_custom_algorithms_created_by ON public.custom_algorithms(created_by);

COMMENT ON TABLE public.custom_algorithms IS 'User-defined WASM algorithms for feature flag optimization';
COMMENT ON COLUMN public.custom_algorithms.kind IS 'Algorithm type: bandit, optimizer, or contextual_bandit';
COMMENT ON COLUMN public.custom_algorithms.wasm_binary IS 'Compiled WASM module binary';
COMMENT ON COLUMN public.custom_algorithms.wasm_hash IS 'SHA256 hash of wasm_binary for caching';

-- Custom algorithm stats for storing algorithm state
CREATE TABLE monitoring.custom_algorithm_stats (
    project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES public.features ON DELETE CASCADE,
    environment_id BIGINT NOT NULL REFERENCES public.environments ON DELETE CASCADE,
    algorithm_id UUID NOT NULL REFERENCES public.custom_algorithms ON DELETE CASCADE,
    variant_key VARCHAR(100) NOT NULL DEFAULT '',
    feature_key VARCHAR(50) NOT NULL,
    environment_key VARCHAR(20) NOT NULL,
    state JSONB DEFAULT '{}'::jsonb NOT NULL,
    evaluations BIGINT DEFAULT 0 NOT NULL,
    successes BIGINT DEFAULT 0 NOT NULL,
    failures BIGINT DEFAULT 0 NOT NULL,
    metric_sum NUMERIC(24, 6) DEFAULT 0 NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    PRIMARY KEY (feature_id, environment_id, algorithm_id, variant_key)
);

CREATE INDEX idx_custom_algorithm_stats_project ON monitoring.custom_algorithm_stats(project_id);
CREATE INDEX idx_custom_algorithm_stats_feature ON monitoring.custom_algorithm_stats(feature_key, environment_key);
CREATE INDEX idx_custom_algorithm_stats_algorithm ON monitoring.custom_algorithm_stats(algorithm_id);

COMMENT ON TABLE monitoring.custom_algorithm_stats IS 'State and statistics for custom WASM algorithms';
COMMENT ON COLUMN monitoring.custom_algorithm_stats.state IS 'Arbitrary JSON state maintained by the WASM algorithm';
COMMENT ON COLUMN monitoring.custom_algorithm_stats.variant_key IS 'Empty string for optimizer algorithms, variant key for bandits';

-- Add custom_algorithm_id to feature_algorithms table
ALTER TABLE public.feature_algorithms
    ADD COLUMN custom_algorithm_id UUID REFERENCES public.custom_algorithms(id) ON DELETE SET NULL;

CREATE INDEX idx_feature_algorithms_custom ON public.feature_algorithms(custom_algorithm_id)
    WHERE custom_algorithm_id IS NOT NULL;

COMMENT ON COLUMN public.feature_algorithms.custom_algorithm_id IS 'Reference to custom WASM algorithm (NULL if using built-in algorithm)';

-- Add constraint: either algorithm_slug or custom_algorithm_id must be set
ALTER TABLE public.feature_algorithms
    ADD CONSTRAINT chk_algorithm_type 
    CHECK (
        (algorithm_slug IS NOT NULL AND custom_algorithm_id IS NULL) OR
        (algorithm_slug IS NULL AND custom_algorithm_id IS NOT NULL)
    );

