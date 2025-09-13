package runner

import (
	"context"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rom8726/pgfixtures"
	"github.com/rom8726/testy"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/rom8726/etoggle/internal"
	"github.com/rom8726/etoggle/internal/config"
)

type Config struct {
	CasesDir string
}

func Run(t *testing.T, testCfg *Config) {
	t.Helper()

	env := NewEnv()

	var connStr string
	var err error
	// Postgres ---------------------------------------------------------------
	dbType := pgfixtures.PostgreSQL
	pgContainer, pgDown := startPostgres(t)
	defer pgDown()

	connStr, err = pgContainer.ConnectionString(t.Context(), "sslmode=disable")
	require.NoError(t, err)
	env.Set("POSTGRES_PORT", extractPort(connStr))

	// Config and App initialization ------------------------------------------
	env.SetUp()
	defer env.CleanUp()

	cfg, err := config.New("")
	if err != nil {
		t.Fatal(err)
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &cfg.Logger,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	time.Sleep(time.Second * 3)
	app, err := internal.NewApp(t.Context(), cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if err := upMigrations(connStr, cfg.MigrationsDir); err != nil {
		t.Fatal(err)
	}

	testyCfg := testy.Config{
		Handler:     app.APIServer.Handler,
		DBType:      dbType,
		CasesDir:    testCfg.CasesDir,
		FixturesDir: "./fixtures",
		ConnStr:     connStr,
	}
	testy.Run(t, &testyCfg)
}

// Postgres -----------------------------------------------------------------
func startPostgres(t *testing.T) (*postgres.PostgresContainer, func()) {
	t.Helper()

	container, err := postgres.Run(t.Context(),
		"postgres:16",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	require.NoError(t, err)

	return container, func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate postgres: %v", err)
		}
	}
}

func extractPort(connStr string) string {
	// Check if it's a MySQL connection string (user:password@tcp(host:port)/dbname)
	if strings.Contains(connStr, "@tcp(") {
		start := strings.Index(connStr, "@tcp(")
		if start == -1 {
			return ""
		}
		start += 5 // Skip "@tcp("

		end := strings.Index(connStr[start:], ")")
		if end == -1 {
			return ""
		}

		hostPort := connStr[start : start+end]
		_, port, _ := net.SplitHostPort(hostPort)
		return port
	}

	// Otherwise, assume it's a PostgreSQL connection string (postgres://user:password@host:port/dbname)
	u, err := url.Parse(connStr)
	if err != nil {
		return ""
	}

	host, port, _ := net.SplitHostPort(u.Host)
	if host == "" {
		return ""
	}

	return port
}
