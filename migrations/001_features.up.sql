create table features (
    id uuid primary key default gen_random_uuid(),
    key text not null unique,         -- machine name, e.g. "new_ui"
    name text not null,               -- human readable name
    description text,                 -- optional description
    kind text not null,               -- "simple" | "multivariant"
    default_variant text not null,    -- any value for boolean, or variant name for multivariant
    enabled boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
