# RBAC 2.0 for Togglr

Goal: design a robust, extensible RBAC that fits our Clean Architecture and current codebase
without forcing per-endpoint ad-hoc checks. This document defines the target architecture,
a migration plan, contracts, repository/service design, middleware/policy approach, examples,
and a rollout/testing plan.

Contents:
- Motivation and current state
- Data model and migrations
- Contracts (interfaces) and packages
- Repository layer (SQL and patterns)
- Service layer (permission checks and helpers)
- Authorization middleware/policy architecture for REST (ogen)
- Integration examples (projects, features, rules)
- Testing plan (unit + functional)
- Rollout plan and backfill


## 1. Motivation and current state

Today we have:
- A stub permissions service (internal/services/permissions/service.go) that only checks
  superuser flag and user presence.
- Endpoint-specific middleware (internal/api/rest/middlewares/permissions.go) that parses
  project ID from URL and decides based on HTTP method (GET/POST). This is too coarse and
  leads to per-endpoint special cases.
- Handlers sometimes call permissions service directly (e.g., project/feature endpoints).

We want:
- Explicit roles and permissions, scoped per project, with a global superuser bypass.
- Centralized, declarative policy mapping per operation, not per path/method.
- Extensibility: ability to add roles/permissions later without code changes.
- Clean layering: use cases depend on contracts; repositories implement contracts.


## 2. Data model and migrations

Tables (Postgres). IDs for roles/permissions are UUIDs; user_id stays int; project_id is our
existing domain.ProjectID type (UUID string in DB). Membership is (user, project) -> role.

Create migrations under migrations/ (numbers are examples, adjust ordering as needed):

-- 0001_create_rbac_tables.sql

```sql
-- built-in roles catalog (seedable)
create table if not exists roles (
    id          uuid primary key default gen_random_uuid(),
    key         text not null unique,     -- e.g. 'project_owner', 'project_member'
    name        text not null,
    description text,
    created_at  timestamptz not null default now()
);

-- permissions catalog
create table if not exists permissions (
    id   uuid primary key default gen_random_uuid(),
    key  text not null unique,            -- e.g. "project.view", "feature.toggle"
    name text not null
);

-- role -> permission mapping
create table if not exists role_permissions (
    id            uuid primary key default gen_random_uuid(),
    role_id       uuid not null references roles(id) on delete cascade,
    permission_id uuid not null references permissions(id) on delete cascade,
    constraint role_permissions_unique unique (role_id, permission_id)
);

-- project memberships with a role
create table if not exists memberships (
    id         uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    user_id    integer not null references users(id) on delete cascade,
    role_id    uuid not null references roles(id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint membership_unique unique (project_id, user_id)
);

-- simple audit for membership changes
create table if not exists membership_audit (
    id             bigserial primary key,
    membership_id  uuid,
    actor_user_id  integer,
    action         text not null,    -- 'create', 'update', 'delete'
    old_value      jsonb,
    new_value      jsonb,
    created_at     timestamptz not null default now()
);
```

Seed defaults (add a 0002 seed migration or use app-time seeding):

```sql
-- roles
insert into roles (key, name, description) values
    ('project_owner',   'Project Owner',   'Full control of project'),
    ('project_manager', 'Project Manager', 'Manage features and rules'),
    ('project_member',  'Project Member',  'Toggle features'),
    ('project_viewer',  'Project Viewer',  'Read-only')
on conflict (key) do nothing;

-- permissions
insert into permissions (key, name) values
    ('project.view',       'View project'),
    ('project.manage',     'Manage project'),
    ('feature.view',       'View features'),
    ('feature.toggle',     'Toggle features'),
    ('feature.manage',     'Manage features'),
    ('rule.manage',        'Manage rules'),
    ('audit.view',         'View audit'),
    ('membership.manage',  'Manage memberships')
on conflict (key) do nothing;

-- grant permissions to default roles
with r as (select id, key from roles), p as (select id, key from permissions)
insert into role_permissions (role_id, permission_id)
select r.id, p.id
from r
join p on (
    -- Owners: everything
    (r.key = 'project_owner') or
    -- Managers: most feature/rule manage + view project
    (r.key = 'project_manager' and p.key in (
        'project.view','feature.view','feature.toggle','feature.manage','rule.manage','audit.view'
    )) or
    -- Members: feature.view + feature.toggle
    (r.key = 'project_member' and p.key in ('feature.view','feature.toggle','project.view')) or
    -- Viewers: project.view + feature.view
    (r.key = 'project_viewer' and p.key in ('project.view','feature.view'))
)
on conflict do nothing;
```

Notes:
- Keep using users.is_superuser as a global bypass.
- Memberships are per project; no team-level RBAC for now.


## 3. Permission vocabulary (Go)

Define the permission key type at the domain layer to keep it framework-agnostic and shared
across all layers. Create internal/domain/permissions.go with constants for reuse:

```go
package domain

type PermKey string

const (
    PermProjectView      PermKey = "project.view"
    PermProjectManage    PermKey = "project.manage"
    PermFeatureView      PermKey = "feature.view"
    PermFeatureToggle    PermKey = "feature.toggle"
    PermFeatureManage    PermKey = "feature.manage"
    PermRuleManage       PermKey = "rule.manage"
    PermAuditView        PermKey = "audit.view"
    PermMembershipManage PermKey = "membership.manage"
)
```

Notes:
- We intentionally place type PermKey in internal/domain to avoid import cycles and to make it the
  single source of truth for permission identifiers.
- Repository and service layers should reference domain.PermKey.

Optional: add role keys if needed for seeding tools.


## 4. Contracts (interfaces)

Add repository contracts to internal/contract (new file internal/contract/rbac_repos.go):

```go
package contract

import (
    "context"

    "github.com/togglr-project/togglr/internal/domain"
)

type RolesRepository interface {
    GetByKey(ctx context.Context, key string) (id string, err error)
}

type PermissionsRepository interface {
    RoleHasPermission(ctx context.Context, roleID string, key domain.PermKey) (bool, error)
}

type MembershipsRepository interface {
    GetForUserProject(ctx context.Context, userID int, projectID domain.ProjectID) (roleID string, err error)
}
```

Extend PermissionsService interface (internal/contract/rbac.go) to support generic permission check
without breaking existing methods:

```go
// Existing methods stay
// Add a generic checker to support fine-grained permissions
HasProjectPermission(ctx context.Context, projectID domain.ProjectID, permKey string) (bool, error)
```

Rationale: handlers and middleware can ask for specific permissions (e.g., feature.toggle) while
legacy code can still call CanAccessProject/CanManageProject.


## 5. Repository layer

Create package internal/repository/rbac with three repositories implementing the contracts.
Use the executor pattern to support TxManager (pkg/db/Tx and context binding).

MembershipsRepository.GetForUserProject:

```sql
select role_id
from memberships
where project_id = $1 and user_id = $2
limit 1;
```

PermissionsRepository.RoleHasPermission (fast path):

```sql
select exists (
  select 1
  from role_permissions rp
  join permissions p on p.id = rp.permission_id
  where rp.role_id = $1 and p.key = $2
);
```

Implementation notes:
- Return ("", nil) when membership not found; treat as no access.
- Wrap DB errors with context (fmt.Errorf("role has permission: %w", err)).
- Convert DB UUIDs to strings; domain.ProjectID is already a string type.


## 6. Service layer

Replace the stubbed permissions service (internal/services/permissions/service.go) with a proper
implementation backed by the new repos. Keep the public shape backward compatible and add the
new helper.

```go
package permissions

import (
    "context"

    etx "github.com/togglr-project/togglr/internal/context"
    "github.com/togglr-project/togglr/internal/contract"
    "github.com/togglr-project/togglr/internal/domain"
)

type Service struct {
    projects contract.ProjectsRepository
    roles    contract.RolesRepository
    perms    contract.PermissionsRepository
    member   contract.MembershipsRepository
}

func New(
    projects contract.ProjectsRepository,
    roles contract.RolesRepository,
    perms contract.PermissionsRepository,
    member contract.MembershipsRepository,
) *Service {
    return &Service{projects: projects, roles: roles, perms: perms, member: member}
}

func (s *Service) isSuper(ctx context.Context) bool { return etx.IsSuper(ctx) }

func (s *Service) HasProjectPermission(
    ctx context.Context,
    projectID domain.ProjectID,
    permKey string,
) (bool, error) {
    if s.isSuper(ctx) {
        return true, nil
    }

    // Verify project exists (preserve current behavior and error mapping)
    if _, err := s.projects.GetByID(ctx, projectID); err != nil {
        return false, err
    }

    userID := etx.UserID(ctx)
    if userID == 0 {
        return false, domain.ErrUserNotFound
    }

    roleID, err := s.member.GetForUserProject(ctx, int(userID), projectID)
    if err != nil || roleID == "" {
        return false, err
    }

    return s.perms.RoleHasPermission(ctx, roleID, domain.PermKey(permKey))
}

func (s *Service) CanAccessProject(ctx context.Context, projectID domain.ProjectID) error {
    ok, err := s.HasProjectPermission(ctx, projectID, string(domain.PermProjectView))
    if err != nil { return err }
    if !ok { return domain.ErrPermissionDenied }
    return nil
}

func (s *Service) CanManageProject(ctx context.Context, projectID domain.ProjectID) error {
    ok, err := s.HasProjectPermission(ctx, projectID, string(domain.PermProjectManage))
    if err != nil { return err }
    if !ok { return domain.ErrPermissionDenied }
    return nil
}

func (s *Service) GetAccessibleProjects(
    ctx context.Context,
    projects []domain.Project,
) ([]domain.Project, error) {
    // Optional: filter with PermProjectView, superuser gets all
    if s.isSuper(ctx) { return projects, nil }
    out := make([]domain.Project, 0, len(projects))
    for _, p := range projects {
        ok, err := s.HasProjectPermission(ctx, p.ID, string(domain.PermProjectView))
        if err != nil { return nil, err }
        if ok { out = append(out, p) }
    }
    return out, nil
}
```

Notes:
- We keep existing API intact and add HasProjectPermission.
- We reuse the current project existence check to preserve error mapping.


## 7. Authorization architecture for REST (ogen)

Problem: endpoint-specific middleware leads to duplication and hidden policy.

Solution: introduce a centralized policy map keyed by OperationID (from specs/server.yml).
- Define a map in internal/api/rest/auth_policies.go:
  - key: operation ID string (as generated by ogen)
  - value: struct { required PermKey; projectIDExtractor func(ctx, req, params) (ProjectID, bool) }
- Implement a wrapper around the generated Handler that, before calling the real method,
  looks up the policy and executes the check using permissions.Service.HasProjectPermission.

Example skeleton:

```go
// internal/api/rest/auth_policies.go
package rest

import (
    "context"

    "github.com/togglr-project/togglr/internal/domain"
)

type policy struct {
    perm domain.PermKey
    extract func(ctx context.Context, params any, req any) (domain.ProjectID, bool)
}

var policies = map[string]policy{ // operationId -> policy
    "ListProjectFeatures": { // example op id
        perm: domain.PermFeatureView,
        extract: extractProjectIDFromParams,
    },
    "CreateProjectFeature": {
        perm: domain.PermFeatureManage,
        extract: extractProjectIDFromParams,
    },
}
```

Wrapper integration: ogen generates a Handler interface in internal/generated/server/oas_server_gen.go.
We can create a struct that implements this interface by delegating to our RestAPI, but with a
pre-dispatch hook that checks policy based on the operation ID. This centralizes authorization.

For operations where no policy is defined, we pass through (e.g., auth/login).

Project ID extractors: keep small helpers for common patterns (path param project_id, feature ->
project lookup via repository, etc.). If extraction requires DB (e.g., feature by id), extractors
can call a read-only use case or repository to determine the project context.

This allows us to avoid editing each handler for permission checks and removes the method-based
logic from middlewares/permissions.go. That middleware can be removed or reduced to attach
project id to context when easily available.


## 8. Integration examples

1) Project update info (PUT /api/v1/projects/{project_id})
- OperationId: UpdateProject (example; check specs/server.yml)
- Policy: perm = project.manage; extractor reads path param project_id
- Behavior: return 403 via ogen error model when check fails

2) List features by project (GET /api/v1/projects/{project_id}/features)
- Policy: perm = feature.view

3) Create feature under project
- Policy: perm = feature.manage

4) Toggle flag variant
- If toggle is a write op: perm = feature.toggle

Handler example using helper (if you choose inline checks instead of wrapper):

```go
// inside RestAPI
func (r *RestAPI) requirePerm(
    ctx context.Context,
    projectID domain.ProjectID,
    perm pkeys.PermKey,
) error {
    ok, err := r.permissionsService.HasProjectPermission(ctx, projectID, string(perm))
    if err != nil { return err }
    if !ok { return domain.ErrPermissionDenied }
    return nil
}
```

Then call r.requirePerm in handlers instead of mixing method/URL logic in middleware.


## 9. Testing

Unit tests:
- Generate mocks for RolesRepository, PermissionsRepository, MembershipsRepository.
- Test Service.HasProjectPermission for cases:
  - superuser bypass
  - no user in context -> ErrUserNotFound
  - no membership -> false, no error
  - membership present but no perm -> false
  - membership present with perm -> true

Functional tests (tests/cases/*):
- Add fixtures for roles/permissions/memberships (tests/fixtures/*).
- Add scenarios for:
  - unauthorized (401)
  - forbidden (403) when lacking permission
  - success (200) when permission exists
- Use dbChecks to assert role grants and membership effects.

Commands:
- go test ./internal/... for unit tests
- go test -tags=integration ./tests/... for functional tests


## 10. Rollout plan

1) Add migrations 0001 and 0002, deploy to staging.
2) Seed default roles and permissions.
3) Backfill memberships for existing projects:
   - Default: project_owner = the project creator user (if tracked) or admins from ops.
   - Otherwise, grant project_viewer to active users as a temporary measure.
4) Implement repositories and the new permissions service.
5) Switch REST authorization to the policy wrapper or helper method; remove URL/method-based
   middleware logic.
6) Update handlers to remove duplicated checks if wrapper is used.
7) Add tests and fixtures, run CI.
8) Monitor logs for permission denied spikes; adjust grants.


## 11. Notes and decisions

- Superuser bypass remains in effect via context (togglr/internal/context.IsSuper).
- RBAC scope is per project. Feature permissions are evaluated in the context of the feature's
  project. Extractors can map feature id -> project id using a read path.
- Performance: we can add small caches for role->permissions and membership lookups with short TTL
  and explicit invalidation on updates (future work).
- Admin UI: later we can build screens for roles, permissions, and memberships management.


## 12. Next steps (for implementers)

- Create contracts (interfaces) and mocks, run make mocks.
- Implement repositories under internal/repository/rbac/.
- Implement the new permissions service and wire it in internal/app.go DI container.
- Replace middlewares/permissions.go with the policy wrapper, or migrate checks to r.requirePerm.
- Write tests and fixtures.



## 13. JWT payload and Frontend adaptation

Context: current JWT claims are defined in internal/domain/jwt.go and include coarse-grained
UserPermissions with booleans and a ProjectPermissions map structured as read/write/manage flags.
To align with the new RBAC, we will carry explicit permission keys to the frontend while keeping
backward compatibility for a transition period.

Proposed approach (non-breaking, phase 1):
- Keep existing fields intact:
  - permissions.can_create_projects, permissions.can_manage_users
  - permissions.project_permissions[project_id] with CanRead/CanWrite/CanDelete/CanManage
- Add an optional PermissionsV2 field carrying explicit permission keys:
  - GlobalPerms: []domain.PermKey
  - ProjectPerms: map[domain.ProjectID][]domain.PermKey

Example (claims.permissionsV2):
```json
{
  "global_perms": ["audit.view"],
  "project_perms": {
    "418aba92-0877-42d7-ac5a-252ebee0d729": [
      "project.view",
      "feature.view",
      "feature.toggle"
    ]
  }
}
```

Frontend migration plan:
1) Phase 1 (backend-first):
   - Backend populates PermissionsV2 in access tokens based on RBAC evaluation.
   - Legacy booleans remain populated (derived from permission sets) to avoid breaking webui.
   - WebUI begins reading PermissionsV2 if present and falls back to legacy booleans otherwise.
2) Phase 2 (webui complete):
   - WebUI fully switches to PermissionsV2 to drive UI visibility and routing guards.
   - Remove reliance on CanRead/CanWrite/CanDelete/CanManage heuristics for project screens.
3) Phase 3 (cleanup):
   - Deprecate and later remove legacy boolean fields from TokenClaims when all clients upgraded.

Compatibility rules (deriving legacy booleans from keys):
- CanRead := has any of [project.view, feature.view]
- CanWrite := has any of [feature.toggle, feature.manage, rule.manage, project.manage]
- CanManage := has project.manage or membership.manage
- CanDelete := CanManage (or define explicitly if you add project.delete)
- CanCreateProjects/CanManageUsers := keep using superuser for now, or introduce global
  permission keys later (e.g., user.manage, project.create) and derive booleans accordingly.

Backend notes (documentation only, code not changed here):
- Keep type PermKey in internal/domain and include it in new JWT field types.
- Access token issuance code should gather permissions via the permissions service and format
  PermissionsV2 accordingly. This doc does not implement it; it sets the contract for future work.
- Maintain token size discipline: include only effective permissions for the current user to avoid
  bloat; consider project-scoped tokens where appropriate.

Testing checklist for webui migration:
- When PermissionsV2 is present, webui uses permission keys to control:
  - Visibility of “Create Feature”, “Edit Project”, “Manage Members”, and toggling switches
  - Access to routes like /projects/:id/settings and /audit
- When PermissionsV2 is absent, webui falls back to legacy fields without breaking.
- Mixed scenarios with superuser (isSuperuser = true) still allow everything.
