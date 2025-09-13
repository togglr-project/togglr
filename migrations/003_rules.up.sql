create table rules (
    id uuid primary key default gen_random_uuid(),
    feature_id uuid not null references features(id) on delete cascade,
    condition jsonb not null,      -- e.g. {"attribute":"country","op":"=","value":"RU"}
    variant text,                  -- which variant to assign if condition matches
    rollout_percent int,           -- optional % rollout
    priority int not null default 0,
    created_at timestamptz not null default now()
);
