-- 1. Add environment_id column (non-null, references environments), enabled column, project_id column, id column.
alter table public.feature_algorithms
    add column if not exists environment_id bigint not null default 0
        references public.environments on delete cascade;

alter table public.feature_algorithms
    add column if not exists enabled bool not null default false;

alter table public.feature_algorithms
    add column if not exists project_id uuid not null default '00000000-0000-0000-0000-000000000000'::uuid
        references public.projects on delete cascade;

alter table public.feature_algorithms
    add column if not exists id uuid not null default gen_random_uuid();

-- 2. Drop old primary key (feature_id, algorithm_slug)
alter table public.feature_algorithms
    drop constraint if exists feature_algorithms_pkey;

-- 3. Recreate primary key with id
alter table public.feature_algorithms
    add constraint feature_algorithms_pkey
        primary key (id);

-- 4. Create unique index
alter table public.feature_algorithms
    add constraint feature_algorithms_uniq
        unique (feature_id, environment_id);
