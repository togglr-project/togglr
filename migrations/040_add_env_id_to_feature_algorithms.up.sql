-- 1. Add environment_id column (non-null, references environments) enabled column
alter table public.feature_algorithms
    add column environment_id bigint not null default 0
        references environments on delete cascade;

alter table public.feature_algorithms
    add column enabled bool not null default false;

comment on column public.feature_algorithms.environment_id is
    'Environment where this algorithm configuration applies.';

-- 2. Drop old primary key (feature_id, algorithm_slug)
alter table public.feature_algorithms
    drop constraint if exists feature_algorithms_pkey;

-- 3. Recreate primary key with environment_id
alter table public.feature_algorithms
    add constraint feature_algorithms_pkey
        primary key (feature_id, environment_id);
