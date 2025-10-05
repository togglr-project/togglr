package runner

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/rom8726/pgfixtures"
	"github.com/rom8726/testy"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v3"

	"github.com/togglr-project/togglr/internal"
	"github.com/togglr-project/togglr/internal/config"
	"github.com/togglr-project/togglr/pkg/crypt"
)

type Config struct {
	CasesDir string
	UsesOTP  bool

	BeforeReq func(app *internal.App) error
	AfterReq  func(app *internal.App) error
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

	// NATS ----------------------------------------------------------------
	natsC, natsDown := startNATS(t)
	defer natsDown()

	natsPort, err := natsC.MappedPort(t.Context(), "4222")
	require.NoError(t, err)
	env.Set("NATS_URL", fmt.Sprintf("nats://127.0.0.1:%d", natsPort.Int()))

	// MailHog ----------------------------------------------------------------
	mailC, mailDown := startMailHog(t)
	defer mailDown()

	mailPort, err := mailC.MappedPort(t.Context(), "1025")
	require.NoError(t, err)
	env.Set("MAILER_ADDR", "localhost:"+mailPort.Port())

	mailPort, _ = mailC.MappedPort(t.Context(), "8025")
	t.Log("MailHog UI: http://localhost:" + mailPort.Port())

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

	if err := upMigrations(t.Context(), connStr, cfg.MigrationsDir); err != nil {
		t.Fatal(err)
	}

	app, err := internal.NewApp(t.Context(), cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if testCfg.UsesOTP {
		modifiedFixtures := setValidFASecretsInFixtures(t, "./fixtures")

		defer func() {
			// --- reset modified fixtures ---
			if len(modifiedFixtures) > 0 {
				for _, filePath := range modifiedFixtures {
					resetFixtureFile(t, filePath)
				}
			}
		}()
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

func RunSDK(t *testing.T, testCfg *Config) {
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

	// NATS ----------------------------------------------------------------
	natsC, natsDown := startNATS(t)
	defer natsDown()

	natsPort, err := natsC.MappedPort(t.Context(), "4222")
	require.NoError(t, err)
	env.Set("NATS_URL", fmt.Sprintf("nats://127.0.0.1:%d", natsPort.Int()))

	// MailHog ----------------------------------------------------------------
	mailC, mailDown := startMailHog(t)
	defer mailDown()

	mailPort, err := mailC.MappedPort(t.Context(), "1025")
	require.NoError(t, err)
	env.Set("MAILER_ADDR", "localhost:"+mailPort.Port())

	mailPort, _ = mailC.MappedPort(t.Context(), "8025")
	t.Log("MailHog UI: http://localhost:" + mailPort.Port())

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

	if err := upMigrations(t.Context(), connStr, cfg.MigrationsDir); err != nil {
		t.Fatal(err)
	}

	app, err := internal.NewApp(t.Context(), cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if testCfg.UsesOTP {
		modifiedFixtures := setValidFASecretsInFixtures(t, "./fixtures")

		defer func() {
			// --- reset modified fixtures ---
			if len(modifiedFixtures) > 0 {
				for _, filePath := range modifiedFixtures {
					resetFixtureFile(t, filePath)
				}
			}
		}()
	}

	testyCfg := testy.Config{
		Handler:     app.SDKServer.Handler,
		DBType:      dbType,
		CasesDir:    testCfg.CasesDir,
		FixturesDir: "./fixtures",
		ConnStr:     connStr,
	}
	testy.Run(t, &testyCfg)
}

// Postgres -----------------------------------------------------------------.
func startPostgres(t *testing.T) (*postgres.PostgresContainer, func()) {
	t.Helper()

	container, err := postgres.Run(t.Context(),
		"timescale/timescaledb:latest-pg17",
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

// MailHog ----------------------------------------------------------------.
func startMailHog(t *testing.T) (testcontainers.Container, func()) {
	t.Helper()

	container, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "togglr-mailhog-test",
			Image:        "mailhog/mailhog:latest",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   wait.ForListeningPort("1025"),
		},
		Started: true,
	})
	require.NoError(t, err)

	return container, func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate mailhog: %v", err)
		}
	}
}

// NATS ----------------------------------------------------------------.
func startNATS(t *testing.T) (testcontainers.Container, func()) {
	t.Helper()

	container, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "togglr-nats-test",
			Image:        "nats:latest",
			ExposedPorts: []string{"4222/tcp"},
			WaitingFor:   wait.ForListeningPort("4222"),
		},
		Started: true,
	})
	require.NoError(t, err)

	return container, func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate nats: %v", err)
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

func setValidFASecretsInFixtures(t *testing.T, fixturesDir string) []string {
	t.Helper()

	var modifiedFiles []string

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		t.Fatal("JWT_SECRET_KEY is not set")
	}

	files, err := filepath.Glob(filepath.Join(fixturesDir, "*.yml"))
	if err != nil {
		t.Fatalf("failed to list fixture files: %v", err)
	}

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file %s: %v", filePath, err)
		}

		var data map[string]any
		if err := yaml.Unmarshal(content, &data); err != nil {
			t.Fatalf("failed to parse YAML in %s: %v", filePath, err)
		}

		modified := false

		usersRaw, ok := data["public.users"]
		if !ok {
			continue
		}

		users, ok := usersRaw.([]any)
		if !ok {
			t.Fatalf("invalid format for public.users in %s", filePath)
		}

		for _, u := range users {
			user, ok := u.(map[string]any)
			if !ok {
				continue
			}

			twoFA, enabled := user["two_fa_enabled"].(bool)
			email, hasEmail := user["email"].(string)

			if enabled && twoFA && hasEmail && email != "" {
				enc := generateValid2FASecret(t, secretKey)
				user["two_fa_secret"] = enc
				modified = true
			}
		}

		if modified {
			var buf bytes.Buffer
			encoder := yaml.NewEncoder(&buf)
			encoder.SetIndent(2)

			if err := encoder.Encode(data); err != nil {
				t.Fatalf("failed to marshal YAML for %s: %v", filePath, err)
			}

			_ = encoder.Close()

			if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
				t.Fatalf("failed to write updated file %s: %v", filePath, err)
			}

			modifiedFiles = append(modifiedFiles, filePath)
		}
	}

	return modifiedFiles
}

func generateValid2FASecret(t *testing.T, secret string) string {
	t.Helper()

	generatedSecret := findOTPSecret(t, "123456")

	encSecret, err := crypt.EncryptAESGCM([]byte(generatedSecret), []byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(encSecret)
}

func findOTPSecret(t *testing.T, targetCode string) string {
	t.Helper()

	now := time.Now()

	for attempt := range 10_000_000 {
		seed := fmt.Sprintf("SECRET%d", attempt)
		secret := base32.StdEncoding.EncodeToString([]byte(seed))
		secret = strings.TrimRight(secret, "=")

		code, err := totp.GenerateCodeCustom(secret, now, totp.ValidateOpts{
			Period:    30,
			Skew:      0,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			panic(err)
		}

		if code == targetCode && totp.Validate(targetCode, secret) {
			t.Logf(">>> Found valid 2FA OTP secret [iteration: %d]\n", attempt+1)

			return secret
		}
	}

	t.Fatal("Failed to find a valid 2FA OTP secret")

	return ""
}

func resetFixtureFile(t *testing.T, filePath string) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "git", "checkout", "--", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	t.Logf("⏪ Resetting fixture: %s\n", filePath)

	_ = cmd.Run()
}
