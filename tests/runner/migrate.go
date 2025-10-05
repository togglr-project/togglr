package runner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func upMigrations(ctx context.Context, connStr, migrationsDir string) error {
	slog.Info("up migrations...")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open postgres connection: %w", err)
	}

	defer func() { _ = db.Close() }()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	pgMigrate, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrations: %w", err)
	}
	if err := pgMigrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("up migrations: no changes")

			return nil
		}

		return fmt.Errorf("up: %w", err)
	}

	slog.Info("up migrations: done")

	// Disable triggers for test environment
	if err := disableTriggersForTests(ctx, db); err != nil {
		return fmt.Errorf("disable triggers: %w", err)
	}

	return nil
}

func disableTriggersForTests(ctx context.Context, db *sql.DB) error {
	slog.Info("disabling triggers for test environment...")

	// First, let's see what triggers exist
	rows, err := db.QueryContext(ctx, `
SELECT trigger_name, event_object_table 
FROM information_schema.triggers 
WHERE trigger_schema = 'public' 
AND trigger_name LIKE '%project%' OR trigger_name LIKE '%tag%' OR trigger_name LIKE '%setting%'
ORDER BY trigger_name`)
	if err != nil {
		slog.Warn("failed to query triggers", "error", err)
	} else {
		defer rows.Close()
		slog.Info("existing triggers:")
		for rows.Next() {
			var triggerName, tableName string
			if err := rows.Scan(&triggerName, &tableName); err == nil {
				slog.Info("trigger found", "name", triggerName, "table", tableName)
			}
		}

		if err := rows.Err(); err != nil {
			return err
		}
	}

	triggers := []string{
		"trg_create_envs_on_project_insert",
		"trg_create_params_on_feature_insert",
		"trg_create_feature_params_on_env_insert",
		"trg_set_default_project_settings",
		"trg_init_project_safety_tags",
	}

	for _, trigger := range triggers {
		query := fmt.Sprintf("DROP TRIGGER IF EXISTS %s ON %s", trigger, getTableName(trigger))
		if _, err := db.ExecContext(ctx, query); err != nil {
			slog.Warn("failed to drop trigger", "trigger", trigger, "error", err)
		} else {
			slog.Info("dropped trigger", "trigger", trigger)
		}
	}

	slog.Info("triggers disabled for test environment")

	return nil
}

//nolint:goconst // fix later
func getTableName(triggerName string) string {
	switch triggerName {
	case "trg_create_envs_on_project_insert":
		return "projects"
	case "trg_create_params_on_feature_insert":
		return "features"
	case "trg_create_feature_params_on_env_insert":
		return "environments"
	case "trg_set_default_project_settings":
		return "projects"
	case "trg_init_project_safety_tags":
		return "projects"
	default:
		return ""
	}
}
