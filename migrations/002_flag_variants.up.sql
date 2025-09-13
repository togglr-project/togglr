create table flag_variants (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    name text not null,            -- e.g. "A", "B"
    rollout_percent int not null,  -- % of traffic (0..100)
    constraint flag_variants_unique unique(feature_id, name)
);
