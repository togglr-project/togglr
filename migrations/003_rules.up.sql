create table rules (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    condition jsonb not null,      -- e.g. {"attribute":"country","op":"=","value":"RU"}
    flag_variant_id uuid not null references flag_variants(id) on delete cascade,
    priority int not null default 0,
    created_at timestamptz not null default now()
);

create index if not exists idx_rules_feature_id on rules(feature_id);
create index if not exists idx_rules_flag_variant_id on rules(flag_variant_id);
-- create index if not exists idx_rules_condition_gin on rules using gin (condition);
