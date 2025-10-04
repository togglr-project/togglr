-- raw events
create table monitoring.feedback_events
(
    id           bigserial,
    feature_id   uuid not null references public.features on delete cascade,
    algorithm_id uuid references public.algorithms,
    variant_key  varchar(100) not null,
    event_type   varchar(50) not null, -- impression, conversion, error, custom
    reward       double precision default 0,
    context      jsonb,
    created_at   timestamptz not null default now(),

    primary key (id, created_at)
);

create index on monitoring.feedback_events (feature_id, algorithm_id, variant_key, created_at desc);

-- transform to hypertable
select create_hypertable('monitoring.feedback_events', 'created_at', if_not_exists => true);

-- enable retention
select add_retention_policy('monitoring.feedback_events', interval '30 days');

---------------------------------------------------

-- 2. Table for first aggregates (для Evaluate)
create table monitoring.feature_algorithm_stats
(
    feature_id    uuid not null references public.features on delete cascade,
    algorithm_id  uuid not null references public.algorithms on delete cascade,
    variant_key   varchar(100) not null,
    impressions   bigint default 0 not null,
    successes     bigint default 0 not null,
    failures      bigint default 0 not null,
    reward_sum    double precision default 0 not null,
    updated_at    timestamptz default now() not null,
    primary key (feature_id, algorithm_id, variant_key)
);

---------------------------------------------------

-- 3. Continuous aggregate (for analytics/dashboards)
create materialized view monitoring.feedback_events_agg
            with (timescaledb.continuous) as
select
    time_bucket('1 hour', created_at) as bucket,
    feature_id,
    algorithm_id,
    variant_key,
    count(*) filter (where event_type = 'impression') as impressions,
    count(*) filter (where event_type = 'conversion') as conversions,
    count(*) filter (where event_type = 'error') as errors,
    sum(reward) as reward_sum
from monitoring.feedback_events
group by bucket, feature_id, algorithm_id, variant_key
with no data;

select add_retention_policy('monitoring.feedback_events_agg', interval '180 days');

select add_continuous_aggregate_policy('monitoring.feedback_events_agg',
    start_offset => interval '2 days',
    end_offset   => interval '1 hour',
    schedule_interval => interval '5 minutes');
