create table license
(
    id           uuid                                   not null primary key,
    license_text text                                   not null,
    issued_at    timestamp with time zone               not null,
    expires_at   timestamp with time zone               not null,
    client_id    uuid                                   not null,
    type         varchar(50)                            not null,
    created_at   timestamp with time zone default now() not null
);

create table license_history
(
    id           uuid                                   not null primary key,
    license_id   uuid                                   not null references license,
    license_text text                                   not null,
    issued_at    timestamp with time zone               not null,
    expires_at   timestamp with time zone               not null,
    client_id    uuid                                   not null,
    type         varchar(50)                            not null,
    created_at   timestamp with time zone default now() not null
);

-- license dates sanity
alter table license
    add constraint license_dates_range
        check (issued_at <= expires_at) not valid;

alter table license_history
    add constraint license_history_dates_range
        check (issued_at <= expires_at) not valid;
