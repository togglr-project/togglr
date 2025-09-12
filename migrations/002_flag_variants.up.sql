create table flag_variants (
    id uuid primary key default gen_random_uuid(),
    flag_id uuid not null references flags(id) on delete cascade,
    name text not null,            -- e.g. "A", "B"
    rollout_percent int not null,  -- % of traffic (0..100)
    constraint flag_variants_unique unique(flag_id, name)
);
