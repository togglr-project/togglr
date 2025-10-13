-- 3. Continuous aggregate for analytics and dashboards (hourly rollup)
create materialized view monitoring.feedback_events_agg
            with (timescaledb.continuous) as
select
    time_bucket('1 hour', created_at) as bucket,
    feature_id,
    algorithm_id,
    variant_key,

    count(*) filter (where event_type = 'evaluation') as evaluations,
    count(*) filter (where event_type = 'success') as successes,
    count(*) filter (where event_type = 'failure') as failures,
    count(*) filter (where event_type = 'error') as errors,

    sum(reward) as metric_sum
from monitoring.feedback_events
group by bucket, feature_id, algorithm_id, variant_key
with no data;

select add_retention_policy('monitoring.feedback_events_agg', interval '180 days');
select add_continuous_aggregate_policy('monitoring.feedback_events_agg',
    start_offset => interval '6 hours',
    end_offset   => interval '5 minutes',
    schedule_interval => interval '1 minute');
