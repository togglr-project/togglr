create table rules (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    condition jsonb not null,      -- e.g. {"attribute":"country","op":"=","value":"RU"}
    flag_variant_id uuid not null references flag_variants(id) on delete cascade,
    priority int not null default 0,
    created_at timestamptz not null default now(),

    constraint rules_feature_flag_variants_unique unique(feature_id, flag_variant_id, condition)
);
