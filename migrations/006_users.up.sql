create table users
(
    id                  serial primary key,
    username            varchar(255)                                       not null unique,
    email               varchar(255)                                       not null unique,
    password_hash       varchar(255)                                       not null,
    is_superuser        boolean                  default false             not null,
    is_active           boolean                  default true              not null,
    created_at          timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at          timestamp with time zone default CURRENT_TIMESTAMP not null,
    last_login          timestamp with time zone default CURRENT_TIMESTAMP,
    is_tmp_password     boolean                  default true,
    two_fa_enabled      boolean                  default false             not null,
    two_fa_secret       text,
    two_fa_confirmed_at timestamp with time zone
);
