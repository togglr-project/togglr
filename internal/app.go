package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/rom8726/etoggl/internal/api/rest"
	"github.com/rom8726/etoggl/internal/config"
	generatedserver "github.com/rom8726/etoggl/internal/generated/server"
	"github.com/rom8726/etoggl/pkg/httpserver"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rom8726/di"
	"github.com/rom8726/etoggl/pkg/db"
	"golang.org/x/sync/errgroup"
)

const (
	ctxTimeout = 10 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool *pgxpool.Pool

	APIServer *httpserver.Server

	container *di.Container
	diApp     *di.App
}

func NewApp(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	pgPool, err := newPostgresConnPool(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	container := di.New()
	diApp := di.NewApp(container)

	app := &App{
		Config:       cfg,
		Logger:       logger,
		container:    container,
		diApp:        diApp,
		PostgresPool: pgPool,
	}

	app.registerComponents()
	app.APIServer, err = app.newAPIServer()
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	return app, nil
}

func (app *App) RegisterComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

func (app *App) ResolveComponent(target any) error {
	return app.container.Resolve(target)
}

func (app *App) ResolveComponentsToStruct(target any) error {
	return app.container.ResolveToStruct(target)
}

func (app *App) Run(ctx context.Context) error {
	techServer, err := app.newTechServer()
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start API server")

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return app.APIServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return techServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.diApp.Run(groupCtx) })

	return group.Wait()
}

func (app *App) Close() {

	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}

}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

func (app *App) registerComponents() {
	app.registerComponent(rest.New)

	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)
	// TODO: register service components
}

func (app *App) newAPIServer() (*httpserver.Server, error) {
	cfg := app.Config.APIServer

	var restAPI generatedserver.Handler
	if err := app.container.Resolve(&restAPI); err != nil {
		return nil, fmt.Errorf("resolve REST API service component: %w", err)
	}

	genServer, err := generatedserver.NewServer(restAPI)

	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      genServer,
	}, nil
}

func (app *App) newTechServer() (*httpserver.Server, error) {
	cfg := app.Config.TechServer
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	router := httprouter.New()
	router.Handle(http.MethodGet, "/health",
		func(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("ok"))
		},
	)

	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	router.HandlerFunc(http.MethodGet, "/debug/pprof", pprof.Index)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)
	router.Handler(http.MethodGet, "/debug/pprof/allocs", pprof.Handler("allocs"))
	router.Handler(http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))
	router.Handler(http.MethodGet, "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handler(http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))
	router.Handler(http.MethodGet, "/debug/pprof/mutex", pprof.Handler("mutex"))
	router.Handler(http.MethodGet, "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      router,
	}, nil
}

func newPostgresConnPool(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.ConnStringWithPoolSize())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgCfg.MaxConnLifetime = cfg.ConnMaxLifetime
	pgCfg.MaxConnLifetimeJitter = time.Second * 5
	pgCfg.MaxConnIdleTime = cfg.MaxIdleConnTime
	pgCfg.HealthCheckPeriod = time.Second * 5

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}
