-- ========================================
-- Rollback Pending Changes (Guard Engine)
-- ========================================

-- Drop triggers first
DROP TRIGGER IF EXISTS trg_project_settings_set_updated_at ON project_settings;
DROP TRIGGER IF EXISTS trg_pce_no_conflict ON pending_change_entities;

-- Drop functions
DROP FUNCTION IF EXISTS set_project_settings_updated_at();
DROP FUNCTION IF EXISTS ensure_no_conflicting_pending();

-- Drop indexes
DROP INDEX IF EXISTS idx_project_settings_project_id;
DROP INDEX IF EXISTS idx_project_approvers_user_id;
DROP INDEX IF EXISTS idx_pending_change_entities_entity;
DROP INDEX IF EXISTS idx_pending_change_entities_pending_id;
DROP INDEX IF EXISTS idx_pending_changes_created_at;
DROP INDEX IF EXISTS idx_pending_changes_project_status;
DROP INDEX IF EXISTS ux_pending_entity_unique;

-- Drop tables (in reverse order due to foreign keys)
DROP VIEW IF EXISTS v_project_effective_settings;
DROP TABLE IF EXISTS project_settings;
DROP TABLE IF EXISTS project_approvers;
DROP TABLE IF EXISTS pending_change_entities;
DROP TABLE IF EXISTS pending_changes;
