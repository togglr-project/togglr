package runner

import (
	"fmt"
	"os"
)

func NewEnv() Env {
	return Env{
		collidedVars: make(map[string]bool),
		redefinedVars: map[string]string{
			"ENVIRONMENT":           "test",
			"LOGGER_LEVEL":          "info",
			"FRONTEND_URL":          "http://localhost:3000",
			"SECRET_KEY":            "secret_key123456",
			"ADMIN_EMAIL":           "admin@togglr.tech",
			"ADMIN_TMP_PASSWORD":    "password543210",
			"API_SERVER_ADDR":       ":8080",
			"TECH_SERVER_ADDR":      ":8081",
			"SDK_SERVER_ADDR":       ":8090",
			"POSTGRES_HOST":         "localhost",
			"POSTGRES_PORT":         "5432",
			"POSTGRES_USER":         "user",
			"POSTGRES_PASSWORD":     "password",
			"POSTGRES_DATABASE":     "test_db",
			"MIGRATIONS_DIR":        "../migrations",
			"JWT_SECRET_KEY":        "secret_key123456",
			"ACCESS_TOKEN_TTL":      "3h",
			"REFRESH_TOKEN_TTL":     "168h",
			"RESET_PASSWORD_TTL":    "8h",
			"EMAIL_SMTP_HOST":       "togglr-mailhog",
			"EMAIL_SMTP_PORT":       "1025",
			"EMAIL_FROM":            "noreply@togglr.local",
			"EMAIL_FROM_NAME":       "Togglr",
			"LDAP_ENABLED":          "false",
			"MAILER_ADDR":           "smtp.togglr.tech:4655",
			"MAILER_USER":           "u1-otdx2j3z0l",
			"MAILER_PASSWORD":       "Togglr123!!",
			"MAILER_FROM":           "noreply@togglr.tech",
			"MAILER_ALLOW_INSECURE": "true",
			"MAILER_USE_TLS":        "false",
		},
	}
}

type Env struct {
	redefinedVars map[string]string
	collidedVars  map[string]bool
}

func (e *Env) SetUp() {
	var err error
	for key, value := range e.redefinedVars {
		if envVar := os.Getenv(key); envVar != "" {
			e.redefinedVars[key] = envVar
			e.collidedVars[key] = true

			continue
		}
		if err = os.Setenv(key, value); err != nil {
			err = fmt.Errorf("can't clear ENV %s: %w", key, err)
			panic(err)
		}
	}
}

func (e *Env) CleanUp() {
	var err error
	for key := range e.redefinedVars {
		if _, ok := e.collidedVars[key]; ok {
			continue
		}
		err = os.Unsetenv(key)
		if err != nil {
			err = fmt.Errorf("can't clear ENV %s: %w", key, err)
			panic(err)
		}
	}
}

func (e *Env) Set(key, val string) {
	e.redefinedVars[key] = val
}

func (e *Env) Get(key string) string {
	return e.redefinedVars[key]
}

func (*Env) SetMock(key, server string) {
	if err := os.Setenv(key, server); err != nil {
		err = fmt.Errorf("can't set mock ENV %s: %w", key, err)
		panic(err)
	}
}
