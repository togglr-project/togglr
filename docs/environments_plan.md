# Environments migration plan (prod/dev/stage)

This document summarizes the current state of the codebase after introducing environments and proposes a step‑by‑step plan to fully migrate the code to the new database schema with environments. The goal is to ensure consistency, reduce duplicated joins by using DB views (v_features_full and v_projects_full), and make environment a first‑class parameter across API, use cases, and repositories.


## 1) Summary of DB changes (027_environments.up.sql)

Key points implemented by migration 027:
- New table environments (per project): keys dev, stage, prod, unique per project; each environment has its own api_key.
- features: enabled and default_variant columns moved to feature_params per environment.
- feature_params(feature_id, environment_id, enabled, default_value, timestamps); PK (feature_id, environment_id).
- rules, flag_variants, feature_schedules now contain environment_id and appropriate FKs and unique constraints are scoped by (feature_id, environment_id, ...).
- audit_log has environment_id.
- pending_changes has environment_id.
- Views introduced:
  - v_features_full(feature, environment_id, environment_key, enabled, default_value, etc.).
  - v_projects_full(project fields + environment_key, api_key).
- Triggers create default environments per project (prod, stage, dev) and default feature_params per feature.

Implications:
- Any query that needs environment-specific attributes must either:
  - Use v_features_full/v_projects_full filtered by environment_key (preferred), or
  - Join base tables with environment_id explicitly (only when truly needed).


## 2) Current state and findings (code audit)

2.1 API layer (internal/api/backend)
- feature_timeline_get.go and feature_test_timeline.go:
  - Both use hardcoded environment "prod" and contain TODOs: "Get environment_key from request parameters".
  - OpenAPI-generated params (internal/generated/server/oas_parameters_gen.go:GetFeatureTimelineParams) do NOT include environment_key.
  - Result: Even if handlers were updated, the parameter is absent from the contract; must start from OpenAPI.

- Other endpoints: Many were already updated to accept environment_key (based on webui code and other handlers), but these two are missing.

2.2 OpenAPI specification (specs/server.yml)
- Generated code confirms that /api/v1/features/{feature_id}/timeline and /api/v1/features/{feature_id}/timeline/test do not have environment_key (query or path). That contradicts the new model.
- Action: Add environment_key to both endpoints; regenerate server with ogen.

2.3 Use cases (internal/usecases)
- Features service:
  - GetExtendedByID(id, environmentKey) correctly loads env-aware feature + lists variants/rules/schedules with env ID.
  - List and ListByProjectID return features for a specific environment (repo filters by environment), OK.
  - ListExtendedByProjectID and ListExtendedByProjectIDFiltered still fetch variants/rules/schedules WITHOUT environment context (ListByFeatureID instead of env-aware methods). This mixes data across environments and is inconsistent with the new model. Needs to switch to env-aware versions (ListByFeatureIDWithEnvID / ...WithEnvID) after resolving envID from environmentKey.
  - Delete/Toggle/UpdateWithChildren correctly resolve env and pass envID to repo and guard logic.

2.4 Repositories
- features/repository.go:
  - GetByID and GetByKey read from v_features_full (good), but these return rows for a single feature+env row. For GetByID (without env), the comment indicates it’s used for audit where env-specific fields aren’t needed. It currently uses v_features_full WHERE id=$1 LIMIT 1 – this returns an arbitrary env row. For environment‑agnostic contexts (pure feature info), reading from base features (without env joins) would be less ambiguous; or use DISTINCT ON id. See Improvements below.
  - GetByIDWithEnvironment / List / ListByProjectID / ListByProjectIDFiltered build explicit joins to feature_params with resolved environment_id. This is valid, but we could simplify/readability by using v_features_full with environment_key filter when appropriate.
- rules/repository.go and feature_schedules repositories: use environment_id already in CRUD and list with env ID (good).
- projects/repository.go:
  - Still selects from projects directly (no v_projects_full usage). For endpoints where environment api_key is required, using v_projects_full filtered by environment_key would simplify.

2.5 Generated code and WebUI
- internal/generated/server shows absence of environment in timeline params.
- WebUI generated client contains many endpoints with environment_key present already, but timeline endpoints seem not aligned.


## 3) What to fix first (priorities)

P1 — Contract correctness
- Add environment_key parameter to:
  - GET /api/v1/features/{feature_id}/timeline
  - POST /api/v1/features/{feature_id}/timeline/test
- Regenerate server with make generate-backend.

P1 — API handlers correctness
- Use the provided environment_key in the two handlers instead of hardcoded "prod". Convert environment_key to envID via environments repo or delegate to usecase helpers.

P2 — Use case consistency for listing extended feature sets
- In Features service:
  - ListExtendedByProjectID and ListExtendedByProjectIDFiltered should resolve envID by (projectID, environmentKey) and then use env-aware list methods for variants, rules, and schedules (ListByFeatureIDWithEnvID equivalents) to avoid mixing environments.

P2 — Repository simplification and view usage
- Where environment_key is available, prefer v_features_full and v_projects_full to avoid manual joins, unless write/update logic or specialized filters require base tables.
- Specific suggestions:
  - features.Repository:
    - For GetByID used only for audit/log, consider switching to base features table or SELECT DISTINCT ON (f.id) from v_features_full to avoid arbitrary env row.
    - For GetByIDWithEnvironment/List/ListByProjectID, consider rewriting using v_features_full filtered by environment_key (e.key) or environment_id when it simplifies SQL (optional optimization).
  - projects.Repository:
    - Introduce methods that read via v_projects_full when environment api_key or env-specific info is needed.

P3 — Tests and coverage
- Update integration tests for new query param in timeline endpoints (cases under tests/cases/...).
- Add unit tests for Features service list extended methods to verify env scoping.


## 4) API/OpenAPI changes (API‑first)

4.1 Update specs/server.yml
- GET /api/v1/features/{feature_id}/timeline
  - Add query parameter environment_key (string, enum dev|stage|prod or free string with validation; consistent with other endpoints).
- POST /api/v1/features/{feature_id}/timeline/test
  - Add query parameter environment_key.
- Ensure descriptions and examples reflect that timelines are environment-specific.

4.2 Regenerate server
- Run: make generate-backend
- Verify new parameters present in internal/generated/server/* for GetFeatureTimelineParams and TestFeatureTimelineParams.

4.3 Frontend alignment
- After regeneration, update WebUI client if it’s generated from the same spec (or ensure compatible usage where it already passes environment_key).


## 5) Handler changes (after codegen)

- internal/api/backend/feature_timeline_get.go
  - Read params.EnvironmentKey
  - Pass to featuresUseCase.GetExtendedByID(ctx, featureID, environmentKey)
- internal/api/backend/feature_test_timeline.go
  - Same as above.

Note: internal/services/features-processor/* stays untouched as requested.


## 6) Use case changes

- internal/usecases/features/features_service.go
  - ListExtendedByProjectID(projectID, environmentKey): resolve env via environmentsRep.GetByProjectIDAndKey, then:
    - Use repo.ListByProjectID(projectID, environmentKey) for features list (already env-aware)
    - For each feature:
      - variants := flagVariantsRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
      - rules := rulesRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
      - schedules := schedulesRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
  - ListExtendedByProjectIDFiltered: same change.

Benefits: guarantees returned extended lists are environment-scoped.


## 7) Repository improvements (optional but recommended)

- features.Repository.GetByID (env-agnostic read for audit):
  - Option A: SELECT FROM features WHERE id=$1 (no env ambiguity) and map only fields available in features.
  - Option B: SELECT DISTINCT ON (id) FROM v_features_full WHERE id=$1 ORDER BY id (explicit) to avoid arbitrary env row.
- projects.Repository:
  - Add helper methods using v_projects_full when endpoints need environment api_key or environment list by project.


## 8) Views usage guidelines

- Use v_features_full when you need environment-dependent attributes (enabled, default_value) and you know environment_key.
- Use base tables for writes and for env-agnostic reads (features without env fields).
- Use v_projects_full to fetch per-environment api_key and list environments by project through SELECT DISTINCT ON (project_id, environment_key) or filtering by environment_key.


## 9) Backward compatibility and migration risks

- Adding environment_key to timeline endpoints is a breaking API change. Coordinate with frontend. Provide defaulting behavior temporarily (e.g., default to "prod" if param omitted) if necessary; however, the preferred approach is to require the parameter.
- Ensure all caches/processors that prefetch features (features-processor) are updated separately (owner will handle).
- Verify permission checks: CanAccessProject should remain environment-agnostic (project-level). If any env-specific permission is introduced later, adjust accordingly.


## 10) Testing plan

- Unit tests:
  - Features service ListExtendedByProjectID and ListExtendedByProjectIDFiltered must return env-scoped children only.
- Integration (tests/cases):
  - Extend existing features timeline tests to pass environment_key and validate timelines reflect environment-specific schedules.
  - Add negative tests for unknown environment_key => 404.
  - Add tests for projects endpoints if they expose env-specific fields via v_projects_full.


## 11) Step-by-step rollout checklist

1. OpenAPI: add environment_key to two timeline endpoints; regenerate backend. 
2. Backend handlers: remove hardcoded "prod" and read param from request; pass to use cases. 
3. Use cases: update extended list methods to use env-aware child lists. 
4. Repositories: optional refactors to use views, and make GetByID env-agnostic. 
5. Tests: update and add tests. 
6. Review and clean any leftover TODOs and direct joins where views suffice.


## 12) Concrete locations to change

- specs/server.yml:
  - /api/v1/features/{feature_id}/timeline (GET) – add query param environment_key
  - /api/v1/features/{feature_id}/timeline/test (POST) – add query param environment_key

- internal/api/backend:
  - feature_timeline_get.go – read env param, remove hardcoded "prod"
  - feature_test_timeline.go – read env param, remove hardcoded "prod"

- internal/usecases/features/features_service.go:
  - ListExtendedByProjectID – fetch children using env-aware repo methods
  - ListExtendedByProjectIDFiltered – same as above

- internal/repository/features/repository.go:
  - Consider making GetByID env-agnostic (base table) or DISTINCT ON (id) in v_features_full.

- internal/repository/projects/repository.go:
  - Add view-backed read methods if/when env-specific fields are needed by API.


## 13) Non-goals in this iteration

- Do not change internal/services/features-processor/* (per request).
- No DB schema changes.


## 14) Done now in this MR

- Added this plan: docs/environments_plan.md with audit results and a concrete, minimal migration plan aligned with Clean Architecture and API-first process.

# Environments migration plan (prod/dev/stage)

This document summarizes the current state of the codebase after introducing environments and proposes a step‑by‑step plan to fully migrate the code to the new database schema with environments. The goal is to ensure consistency, reduce duplicated joins by using DB views (v_features_full and v_projects_full), and make environment a first‑class parameter across API, use cases, and repositories.


## 1) Summary of DB changes (027_environments.up.sql)

Key points implemented by migration 027:
- New table environments (per project): keys dev, stage, prod, unique per project; each environment has its own api_key.
- features: enabled and default_variant columns moved to feature_params per environment.
- feature_params(feature_id, environment_id, enabled, default_value, timestamps); PK (feature_id, environment_id).
- rules, flag_variants, feature_schedules now contain environment_id and appropriate FKs and unique constraints are scoped by (feature_id, environment_id, ...).
- audit_log has environment_id.
- pending_changes has environment_id.
- Views introduced:
  - v_features_full(feature, environment_id, environment_key, enabled, default_value, etc.).
  - v_projects_full(project fields + environment_key, api_key).
- Triggers create default environments per project (prod, stage, dev) and default feature_params per feature.

Implications:
- Any query that needs environment-specific attributes must either:
  - Use v_features_full/v_projects_full filtered by environment_key (preferred), or
  - Join base tables with environment_id explicitly (only when truly needed).


## 2) Current state and findings (code audit)

2.1 API layer (internal/api/backend)
- feature_timeline_get.go and feature_test_timeline.go:
  - Both use hardcoded environment "prod" and contain TODOs: "Get environment_key from request parameters".
  - OpenAPI-generated params (internal/generated/server/oas_parameters_gen.go:GetFeatureTimelineParams) do NOT include environment_key.
  - Result: Even if handlers were updated, the parameter is absent from the contract; must start from OpenAPI.

- Other endpoints: Many were already updated to accept environment_key (based on webui code and other handlers), but these two are missing.

2.2 OpenAPI specification (specs/server.yml)
- Generated code confirms that /api/v1/features/{feature_id}/timeline and /api/v1/features/{feature_id}/timeline/test do not have environment_key (query or path). That contradicts the new model.
- Action: Add environment_key to both endpoints; regenerate server with ogen.

2.3 Use cases (internal/usecases)
- Features service:
  - GetExtendedByID(id, environmentKey) correctly loads env-aware feature + lists variants/rules/schedules with env ID.
  - List and ListByProjectID return features for a specific environment (repo filters by environment), OK.
  - ListExtendedByProjectID and ListExtendedByProjectIDFiltered still fetch variants/rules/schedules WITHOUT environment context (ListByFeatureID instead of env-aware methods). This mixes data across environments and is inconsistent with the new model. Needs to switch to env-aware versions (ListByFeatureIDWithEnvID / ...WithEnvID) after resolving envID from environmentKey.
  - Delete/Toggle/UpdateWithChildren correctly resolve env and pass envID to repo and guard logic.

2.4 Repositories
- features/repository.go:
  - GetByID and GetByKey read from v_features_full (good), but these return rows for a single feature+env row. For GetByID (without env), the comment indicates it’s used for audit where env-specific fields aren’t needed. It currently uses v_features_full WHERE id=$1 LIMIT 1 – this returns an arbitrary env row. For environment‑agnostic contexts (pure feature info), reading from base features (without env joins) would be less ambiguous; or use DISTINCT ON id. See Improvements below.
  - GetByIDWithEnvironment / List / ListByProjectID / ListByProjectIDFiltered build explicit joins to feature_params with resolved environment_id. This is valid, but we could simplify/readability by using v_features_full with environment_key filter when appropriate.
- rules/repository.go and feature_schedules repositories: use environment_id already in CRUD and list with env ID (good).
- projects/repository.go:
  - Still selects from projects directly (no v_projects_full usage). For endpoints where environment api_key is required, using v_projects_full filtered by environment_key would simplify.

2.5 Generated code and WebUI
- internal/generated/server shows absence of environment in timeline params.
- WebUI generated client contains many endpoints with environment_key present already, but timeline endpoints seem not aligned.


## 3) What to fix first (priorities)

P1 — Contract correctness
- Add environment_key parameter to:
  - GET /api/v1/features/{feature_id}/timeline
  - POST /api/v1/features/{feature_id}/timeline/test
- Regenerate server with make generate-backend.

P1 — API handlers correctness
- Use the provided environment_key in the two handlers instead of hardcoded "prod". Convert environment_key to envID via environments repo or delegate to usecase helpers.

P2 — Use case consistency for listing extended feature sets
- In Features service:
  - ListExtendedByProjectID and ListExtendedByProjectIDFiltered should resolve envID by (projectID, environmentKey) and then use env-aware list methods for variants, rules, and schedules (ListByFeatureIDWithEnvID equivalents) to avoid mixing environments.

P2 — Repository simplification and view usage
- Where environment_key is available, prefer v_features_full and v_projects_full to avoid manual joins, unless write/update logic or specialized filters require base tables.
- Specific suggestions:
  - features.Repository:
    - For GetByID used only for audit/log, consider switching to base features table or SELECT DISTINCT ON (f.id) from v_features_full to avoid arbitrary env row.
    - For GetByIDWithEnvironment/List/ListByProjectID, consider rewriting using v_features_full filtered by environment_key (e.key) or environment_id when it simplifies SQL (optional optimization).
  - projects.Repository:
    - Introduce methods that read via v_projects_full when environment api_key or env-specific info is needed.

P3 — Tests and coverage
- Update integration tests for new query param in timeline endpoints (cases under tests/cases/...).
- Add unit tests for Features service list extended methods to verify env scoping.


## 4) API/OpenAPI changes (API‑first)

4.1 Update specs/server.yml
- GET /api/v1/features/{feature_id}/timeline
  - Add query parameter environment_key (string, enum dev|stage|prod or free string with validation; consistent with other endpoints).
- POST /api/v1/features/{feature_id}/timeline/test
  - Add query parameter environment_key.
- Ensure descriptions and examples reflect that timelines are environment-specific.

4.2 Regenerate server
- Run: make generate-backend
- Verify new parameters present in internal/generated/server/* for GetFeatureTimelineParams and TestFeatureTimelineParams.

4.3 Frontend alignment
- After regeneration, update WebUI client if it’s generated from the same spec (or ensure compatible usage where it already passes environment_key).


## 5) Handler changes (after codegen)

- internal/api/backend/feature_timeline_get.go
  - Read params.EnvironmentKey
  - Pass to featuresUseCase.GetExtendedByID(ctx, featureID, environmentKey)
- internal/api/backend/feature_test_timeline.go
  - Same as above.

Note: internal/services/features-processor/* stays untouched as requested.


## 6) Use case changes

- internal/usecases/features/features_service.go
  - ListExtendedByProjectID(projectID, environmentKey): resolve env via environmentsRep.GetByProjectIDAndKey, then:
    - Use repo.ListByProjectID(projectID, environmentKey) for features list (already env-aware)
    - For each feature:
      - variants := flagVariantsRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
      - rules := rulesRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
      - schedules := schedulesRep.ListByFeatureIDWithEnvID(feature.ID, env.ID)
  - ListExtendedByProjectIDFiltered: same change.

Benefits: guarantees returned extended lists are environment-scoped.


## 7) Repository improvements (optional but recommended)

- features.Repository.GetByID (env-agnostic read for audit):
  - Option A: SELECT FROM features WHERE id=$1 (no env ambiguity) and map only fields available in features.
  - Option B: SELECT DISTINCT ON (id) FROM v_features_full WHERE id=$1 ORDER BY id (explicit) to avoid arbitrary env row.
- projects.Repository:
  - Add helper methods using v_projects_full when endpoints need environment api_key or environment list by project.


## 8) Views usage guidelines

- Use v_features_full when you need environment-dependent attributes (enabled, default_value) and you know environment_key.
- Use base tables for writes and for env-agnostic reads (features without env fields).
- Use v_projects_full to fetch per-environment api_key and list environments by project through SELECT DISTINCT ON (project_id, environment_key) or filtering by environment_key.


## 9) Backward compatibility and migration risks

- Adding environment_key to timeline endpoints is a breaking API change. Coordinate with frontend. Provide defaulting behavior temporarily (e.g., default to "prod" if param omitted) if necessary; however, the preferred approach is to require the parameter.
- Ensure all caches/processors that prefetch features (features-processor) are updated separately (owner will handle).
- Verify permission checks: CanAccessProject should remain environment-agnostic (project-level). If any env-specific permission is introduced later, adjust accordingly.


## 10) Testing plan

- Unit tests:
  - Features service ListExtendedByProjectID and ListExtendedByProjectIDFiltered must return env-scoped children only.
- Integration (tests/cases):
  - Extend existing features timeline tests to pass environment_key and validate timelines reflect environment-specific schedules.
  - Add negative tests for unknown environment_key => 404.
  - Add tests for projects endpoints if they expose env-specific fields via v_projects_full.


## 11) Step-by-step rollout checklist

1. OpenAPI: add environment_key to two timeline endpoints; regenerate backend. 
2. Backend handlers: remove hardcoded "prod" and read param from request; pass to use cases. 
3. Use cases: update extended list methods to use env-aware repo methods. 
4. Repositories: optional refactors to use views, and make GetByID env-agnostic. 
5. Tests: update and add tests. 
6. Review and clean any leftover TODOs and direct joins where views suffice.


## 12) Concrete locations to change

- specs/server.yml:
  - /api/v1/features/{feature_id}/timeline (GET) – add query param environment_key
  - /api/v1/features/{feature_id}/timeline/test (POST) – add query param environment_key

- internal/api/backend:
  - feature_timeline_get.go – read env param, remove hardcoded "prod"
  - feature_test_timeline.go – read env param, remove hardcoded "prod"

- internal/usecases/features/features_service.go:
  - ListExtendedByProjectID – fetch children using env-aware repo methods
  - ListExtendedByProjectIDFiltered – same as above

- internal/repository/features/repository.go:
  - Consider making GetByID env-agnostic (base table) or DISTINCT ON (id) in v_features_full.

- internal/repository/projects/repository.go:
  - Add view-backed read methods if/when env-specific fields are needed by API.


## 13) Non-goals in this iteration

- Do not change internal/services/features-processor/* (per request).
- No DB schema changes.


## 14) Done now in this MR

- Added this plan: docs/environments_plan.md with audit results and a concrete, minimal migration plan aligned with Clean Architecture and API-first process.

### 15) Follow-up: Implemented fixes for env consistency (current session)

- API spec and handlers for timeline endpoints now require and pass environment_key. Server regenerated.
- Use cases:
  - ListExtendedByProjectID and ListExtendedByProjectIDFiltered now load child entities (variants, rules, schedules) strictly within the requested environment using ...WithEnvID repository methods.
  - Toggle now updates feature_params (enabled per environment) instead of features table and reloads the feature via GetByIDWithEnvironment.
  - UpdateWithChildren reconciles variants and rules only within the specified environment: lists existing by env, sets EnvironmentID on create/update, and deletes only within env scope.
- Repositories: feature_params/flagvariants/rules already support env IDs; used accordingly.
- Build: go build ./... passes.
