create table ldap_sync_stats
(
    id              serial primary key,
    sync_session_id uuid                                                          not null unique,
    start_time      timestamp with time zone default now()                        not null,
    end_time        timestamp with time zone,
    duration        varchar(50),
    total_users     integer                  default 0                            not null,
    synced_users    integer                  default 0                            not null,
    errors          integer                  default 0                            not null,
    warnings        integer                  default 0                            not null,
    status          varchar(20)              default 'running'::character varying not null,
    error_message   text
);

create index idx_ldap_sync_stats_sync_session_id
    on ldap_sync_stats (sync_session_id);

create index idx_ldap_sync_stats_start_time
    on ldap_sync_stats (start_time);

create index idx_ldap_sync_stats_status
    on ldap_sync_stats (status);

create table ldap_sync_logs
(
    id                 serial primary key,
    timestamp          timestamp with time zone default now() not null,
    level              varchar(10)                            not null,
    message            text                                   not null,
    username           varchar(255),
    details            text,
    sync_session_id    uuid                                   not null,
    stack_trace        text,
    ldap_error_code    integer,
    ldap_error_message text
);

create index idx_ldap_sync_logs_timestamp
    on ldap_sync_logs (timestamp);

create index idx_ldap_sync_logs_level
    on ldap_sync_logs (level);

create index idx_ldap_sync_logs_sync_session_id
    on ldap_sync_logs (sync_session_id);

create index idx_ldap_sync_logs_username
    on ldap_sync_logs (username);
