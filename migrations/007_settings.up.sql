create table settings
(
    id          serial primary key,
    name        varchar(50) not null unique,
    value       jsonb        not null,
    description varchar(300),
    created_at  timestamp with time zone default now(),
    updated_at  timestamp with time zone default now()
);

create trigger trg_settings_set_updated_at
    before update on settings
    for each row execute function set_updated_at();
