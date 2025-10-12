-- 1. Raw feedback events (source of truth)
create table monitoring.feedback_events
(
    id           bigserial,
    feature_id   uuid not null references public.features on delete cascade,
    algorithm_id uuid references public.algorithms,
    variant_key  varchar(100) not null,
    event_type   varchar(50) not null, -- evaluation, success, failure, error, custom
    reward       double precision default 0,
    context      jsonb,
    created_at   timestamptz not null default now(),

    primary key (id, created_at)
);

comment on table monitoring.feedback_events is
    'Raw feature feedback events from SDK (evaluation, success, failure, error, etc.)';
comment on column monitoring.feedback_events.reward is
    'Optional numeric value representing feature metric or reward (used by bandit algorithms)';

create index on monitoring.feedback_events (feature_id, algorithm_id, variant_key, created_at desc);

-- Convert to hypertable
select create_hypertable('monitoring.feedback_events', 'created_at', if_not_exists => true);

-- Enable a 30-day retention policy
select add_retention_policy('monitoring.feedback_events', interval '30 days');

---------------------------------------------------

-- 2. Aggregated per-variant statistics for runtime evaluation (used by bandit algorithms)
create table monitoring.feature_algorithm_stats
(
    feature_id    uuid not null references public.features on delete cascade,
    algorithm_id  uuid not null references public.algorithms on delete cascade,
    variant_key   varchar(100) not null,

    evaluations   bigint default 0 not null,     -- number of feature evaluations
    successes     bigint default 0 not null,     -- positive outcomes
    failures      bigint default 0 not null,     -- negative outcomes
    metric_sum    double precision default 0 not null, -- accumulated metric/reward values

    updated_at    timestamptz default now() not null,
    primary key (feature_id, algorithm_id, variant_key)
);

comment on table monitoring.feature_algorithm_stats is
    'Current per-variant statistics for algorithm evaluation (kept in sync with feedback events)';
comment on column monitoring.feature_algorithm_stats.evaluations is
    'Number of times the feature was evaluated (activated)';
comment on column monitoring.feature_algorithm_stats.metric_sum is
    'Sum of metric or reward values for this variant';
