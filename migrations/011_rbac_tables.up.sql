-- RBAC tables
create table if not exists roles (
    id          uuid primary key default gen_random_uuid(),
    key         varchar(50) not null unique,
    name        varchar(50) not null,
    description varchar(300),
    created_at  timestamptz not null default now()
);

create table if not exists permissions (
    id   uuid primary key default gen_random_uuid(),
    key  varchar(50) not null unique,
    name varchar(50) not null
);

create table if not exists role_permissions (
    id            uuid primary key default gen_random_uuid(),
    role_id       uuid not null references roles(id) on delete cascade,
    permission_id uuid not null references permissions(id) on delete cascade,
    constraint role_permissions_unique unique (role_id, permission_id)
);

create table if not exists memberships (
    id         uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    user_id    integer not null references users(id) on delete cascade,
    role_id    uuid not null references roles(id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint membership_unique unique (project_id, user_id)
);

create table if not exists membership_audit (
    id             bigserial primary key,
    membership_id  uuid,
    actor_user_id  integer,
    action         varchar(50) not null,
    old_value      jsonb,
    new_value      jsonb,
    created_at     timestamptz not null default now()
);

create trigger trg_memberships_set_updated_at
    before update on memberships
    for each row execute function set_updated_at();

create index if not exists idx_memberships_project_id on memberships(project_id);
create index if not exists idx_memberships_user_id on memberships(user_id);
create index if not exists idx_memberships_role_id on memberships(role_id);
create index if not exists idx_role_permissions_permission_id on role_permissions(permission_id);
