create table features (
    id uuid primary key default gen_random_uuid(),
    key varchar(50) not null unique,         -- machine name, e.g. "new_ui"
    name varchar(50) not null,               -- human readable name
    description varchar(300),                -- optional description
    kind varchar(15) not null,               -- "simple" | "multivariant"
    default_variant varchar(128) not null,   -- any value for simple, or variant name for multivariant
    enabled boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

alter table features
    add constraint features_kind_allowed
        check (kind in ('simple','multivariant')) not valid;

-- Common helper to maintain updated_at
create or replace function set_updated_at() returns trigger as $$
begin
    new.updated_at := now();
    return new;
end;
$$ language plpgsql;

create trigger trg_features_set_updated_at
    before update on features
    for each row execute function set_updated_at();
