create table license
(
    id           uuid                                   not null primary key,
    license_text text                                   not null,
    issued_at    timestamp with time zone               not null,
    expires_at   timestamp with time zone               not null,
    client_id    text                                   not null,
    type         text                                   not null,
    created_at   timestamp with time zone default now() not null
);

create table license_history
(
    id           uuid                                   not null primary key,
    license_id   uuid                                   not null references license,
    license_text text                                   not null,
    issued_at    timestamp with time zone               not null,
    expires_at   timestamp with time zone               not null,
    client_id    text                                   not null,
    type         text                                   not null,
    created_at   timestamp with time zone default now() not null
);
