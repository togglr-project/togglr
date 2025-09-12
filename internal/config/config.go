package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	prefix = ""
)

type Config struct {
	Logger        Logger   `envconfig:"LOGGER"`
	APIServer     Server   `envconfig:"API_SERVER"`
	TechServer    Server   `envconfig:"TECH_SERVER"`
	MigrationsDir string   `envconfig:"MIGRATIONS_DIR" default:"./migrations"`
	Postgres      Postgres `envconfig:"POSTGRES"`
}

type Logger struct {
	Lvl string `envconfig:"LEVEL" default:"info"`
}

func (l *Logger) Level() slog.Level {
	switch l.Lvl {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("invalid logger level " + l.Lvl)
	}
}

type Server struct {
	Addr         string        `envconfig:"ADDR" required:"true"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
}

type Postgres struct {
	User            string        `envconfig:"USER" required:"true"`
	Password        string        `envconfig:"PASSWORD" required:"true"`
	Host            string        `envconfig:"HOST" required:"true"`
	Port            string        `envconfig:"PORT" default:"5432"`
	Database        string        `envconfig:"DATABASE" required:"true"`
	MaxIdleConnTime time.Duration `envconfig:"MAX_IDLE_CONN_TIME" default:"5m"`
	MaxConns        int           `envconfig:"MAX_CONNS" default:"20"`
	ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"10m"`
}

func (db *Postgres) ConnString() string {
	var user *url.Userinfo

	if db.User != "" {
		var pass string

		if db.Password != "" {
			pass = db.Password
		}

		user = url.UserPassword(db.User, pass)
	}

	params := url.Values{}
	params.Set("sslmode", "disable")
	params.Set("connect_timeout", "10")

	uri := url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     net.JoinHostPort(db.Host, db.Port),
		Path:     db.Database,
		RawQuery: params.Encode(),
	}

	return uri.String()
}

func (db *Postgres) ConnStringWithPoolSize() string {
	connString := db.ConnString()

	return connString + fmt.Sprintf("&pool_max_conns=%d", db.MaxConns)
}

func New(filePath string) (*Config, error) {
	cfg := &Config{}

	if filePath != "" {
		if err := godotenv.Load(filePath); err != nil {
			return nil, fmt.Errorf("error loading env file: %w", err)
		}
	}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func MustNew(filePath string) *Config {
	cfg, err := New(filePath)
	if err != nil {
		panic(err)
	}

	return cfg
}
