-- 1. Add environment_id column (non-null, references environments)
alter table public.feature_algorithms
    add column environment_id bigint not null default 0
        references environments on delete cascade;

comment on column public.feature_algorithms.environment_id is
    'Environment where this algorithm configuration applies.';

-- 2. Drop old primary key (feature_id, algorithm_id)
alter table public.feature_algorithms
    drop constraint if exists feature_algorithms_pkey;

-- 3. Recreate primary key with environment_id
alter table public.feature_algorithms
    add constraint feature_algorithms_pkey
        primary key (feature_id, algorithm_id, environment_id);

-- 4. Add nullable reference to flag_variants (one algorithm can target a specific variant)
alter table public.feature_algorithms
    add column flag_variant_id uuid
        references public.flag_variants
            on delete cascade;

comment on column public.feature_algorithms.flag_variant_id is
    'Optional reference to a specific flag variant. NULL means the algorithm applies to the entire feature.';

-- 5. Ensure only one algorithm per variant within environment
alter table public.feature_algorithms
    add constraint feature_algorithms_unique_variant_per_env
        unique (feature_id, flag_variant_id, environment_id);
