create table flag_variants (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    name varchar(128) not null,            -- e.g. "A", "B"
    rollout_percent int not null,  -- % of traffic (0..100)
    constraint flag_variants_unique unique(feature_id, name)
);

alter table flag_variants
    add constraint flag_variants_rollout_percent_range
        check (rollout_percent between 0 and 100) not valid;

create index if not exists idx_flag_variants_feature_id on flag_variants(feature_id);
