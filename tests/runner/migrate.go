package runner

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func upMigrations(connStr, migrationsDir string) error {
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
	if err := disableTriggersForTests(db); err != nil {
		return fmt.Errorf("disable triggers: %w", err)
	}

	return nil
}

func disableTriggersForTests(db *sql.DB) error {
	slog.Info("disabling triggers for test environment...")

	triggers := []string{
		"trg_create_envs_on_project_insert",
		"trg_create_params_on_feature_insert",
		"trg_create_feature_params_on_env_insert",
	}

	for _, trigger := range triggers {
		query := fmt.Sprintf("DROP TRIGGER IF EXISTS %s ON %s", trigger, getTableName(trigger))
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("drop trigger %s: %w", trigger, err)
		}
	}

	slog.Info("triggers disabled for test environment")

	return nil
}

func getTableName(triggerName string) string {
	switch triggerName {
	case "trg_create_envs_on_project_insert":
		return "projects"
	case "trg_create_params_on_feature_insert":
		return "features"
	case "trg_create_feature_params_on_env_insert":
		return "environments"
	default:
		return ""
	}
}
