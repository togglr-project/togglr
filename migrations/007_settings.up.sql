create table settings
(
    id          serial primary key,
    name        varchar(255) not null unique,
    value       jsonb        not null,
    description text,
    created_at  timestamp with time zone default now(),
    updated_at  timestamp with time zone default now()
);
