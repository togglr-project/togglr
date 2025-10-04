package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"path"
	"time"

	apibackend "github.com/togglr-project/togglr/internal/api/backend"
	"github.com/togglr-project/togglr/internal/api/backend/middlewares"
	apisdk "github.com/togglr-project/togglr/internal/api/sdk"
	wsapi "github.com/togglr-project/togglr/internal/api/ws"
	wsmiddlewares "github.com/togglr-project/togglr/internal/api/ws/middlewares"
	"github.com/togglr-project/togglr/internal/config"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	eventsbus "github.com/togglr-project/togglr/internal/event-bus"
	generatedsdk "github.com/togglr-project/togglr/internal/generated/sdkserver"
	generatedserver "github.com/togglr-project/togglr/internal/generated/server"
	natsmq "github.com/togglr-project/togglr/internal/infra/mq/nats"
	"github.com/togglr-project/togglr/internal/repository/auditlog"
	"github.com/togglr-project/togglr/internal/repository/categories"
	dashboardrepo "github.com/togglr-project/togglr/internal/repository/dashboard"
	environmentsrepo "github.com/togglr-project/togglr/internal/repository/environments"
	"github.com/togglr-project/togglr/internal/repository/errorreports"
	featureparamsrepo "github.com/togglr-project/togglr/internal/repository/feature_params"
	featuretagsrepo "github.com/togglr-project/togglr/internal/repository/feature_tags"
	"github.com/togglr-project/togglr/internal/repository/features"
	"github.com/togglr-project/togglr/internal/repository/featureschedules"
	"github.com/togglr-project/togglr/internal/repository/flagvariants"
	"github.com/togglr-project/togglr/internal/repository/guard_service"
	"github.com/togglr-project/togglr/internal/repository/ldapsynclogs"
	"github.com/togglr-project/togglr/internal/repository/ldapsyncstats"
	"github.com/togglr-project/togglr/internal/repository/licenses"
	"github.com/togglr-project/togglr/internal/repository/pending_changes"
	"github.com/togglr-project/togglr/internal/repository/productinfo"
	"github.com/togglr-project/togglr/internal/repository/project_approvers"
	"github.com/togglr-project/togglr/internal/repository/project_settings"
	"github.com/togglr-project/togglr/internal/repository/projects"
	"github.com/togglr-project/togglr/internal/repository/rbac"
	realtimerepo "github.com/togglr-project/togglr/internal/repository/realtime"
	ruleattributesrepo "github.com/togglr-project/togglr/internal/repository/ruleattributes"
	"github.com/togglr-project/togglr/internal/repository/rules"
	segmentsrepo "github.com/togglr-project/togglr/internal/repository/segments"
	"github.com/togglr-project/togglr/internal/repository/settings"
	"github.com/togglr-project/togglr/internal/repository/tags"
	"github.com/togglr-project/togglr/internal/repository/users"
	ratelimiter2fa "github.com/togglr-project/togglr/internal/services/2fa/ratelimiter"
	featuresprocessor "github.com/togglr-project/togglr/internal/services/features-processor"
	guardengine "github.com/togglr-project/togglr/internal/services/guard-engine"
	"github.com/togglr-project/togglr/internal/services/ldap"
	"github.com/togglr-project/togglr/internal/services/notification-channels/email"
	"github.com/togglr-project/togglr/internal/services/permissions"
	realtimechanges "github.com/togglr-project/togglr/internal/services/realtime-changes"
	ssoprovidermanager "github.com/togglr-project/togglr/internal/services/sso/provider-manager"
	samlprovider "github.com/togglr-project/togglr/internal/services/sso/saml"
	"github.com/togglr-project/togglr/internal/services/tokenizer"
	categoriesusecase "github.com/togglr-project/togglr/internal/usecases/categories"
	dashboardusecase "github.com/togglr-project/togglr/internal/usecases/dashboard"
	environmentsusecase "github.com/togglr-project/togglr/internal/usecases/environments"
	errorreportsusecase "github.com/togglr-project/togglr/internal/usecases/errorreports"
	featuretagsusecase "github.com/togglr-project/togglr/internal/usecases/feature-tags"
	featuresusecase "github.com/togglr-project/togglr/internal/usecases/features"
	featureschedulesusecase "github.com/togglr-project/togglr/internal/usecases/featureschedules"
	flagvariantsusecase "github.com/togglr-project/togglr/internal/usecases/flagvariants"
	ldapusecase "github.com/togglr-project/togglr/internal/usecases/ldap"
	pendingchangesusecase "github.com/togglr-project/togglr/internal/usecases/pending-changes"
	projectsettingsusecase "github.com/togglr-project/togglr/internal/usecases/project-settings"
	projectsusecase "github.com/togglr-project/togglr/internal/usecases/projects"
	rbacusecase "github.com/togglr-project/togglr/internal/usecases/rbac"
	ruleattributesusecase "github.com/togglr-project/togglr/internal/usecases/ruleattributes"
	rulesusecase "github.com/togglr-project/togglr/internal/usecases/rules"
	segmentsusecase "github.com/togglr-project/togglr/internal/usecases/segments"
	settingsusecase "github.com/togglr-project/togglr/internal/usecases/settings"
	tagsusecase "github.com/togglr-project/togglr/internal/usecases/tags"
	usersusecase "github.com/togglr-project/togglr/internal/usecases/users"
	"github.com/togglr-project/togglr/pkg/db"
	"github.com/togglr-project/togglr/pkg/httpserver"
	pkgmiddlewares "github.com/togglr-project/togglr/pkg/httpserver/middlewares"
	"github.com/togglr-project/togglr/pkg/passworder"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rom8726/di"
	"golang.org/x/sync/errgroup"
)

const (
	ctxTimeout = 10 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool *pgxpool.Pool
	Bus          *natsmq.NATSMq

	APIServer *httpserver.Server
	SDKServer *httpserver.Server
	WSServer  *httpserver.Server

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

	bus, err := natsmq.New(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("init message queue: %w", err)
	}

	container := di.New()
	diApp := di.NewApp(container)

	app := &App{
		Config:       cfg,
		Logger:       logger,
		container:    container,
		diApp:        diApp,
		PostgresPool: pgPool,
		Bus:          bus,
	}

	app.registerComponents()

	app.APIServer, err = app.newAPIServer()
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	app.SDKServer, err = app.newSDKServer()
	if err != nil {
		return nil, fmt.Errorf("create SDK server: %w", err)
	}

	app.WSServer, err = app.newWSServer()
	if err != nil {
		return nil, fmt.Errorf("create WS server: %w", err)
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
	// Check and create superuser if needed
	if app.Config.AdminEmail != "" {
		if err := app.ensureSuperuser(ctx); err != nil {
			app.Logger.Error("Failed to ensure superuser exists", "error", err)
		}
	}

	techServer, err := app.newTechServer()
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start API server")

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return app.APIServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.SDKServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.WSServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return techServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.diApp.Run(groupCtx) })

	return group.Wait()
}

func (app *App) Close() {
	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}
	if app.Bus != nil {
		app.Bus.Close()
	}
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

func (app *App) registerComponents() {
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)
	app.registerComponent(func() *natsmq.NATSMq {
		return app.Bus
	})

	// Register repositories
	app.registerComponent(projects.New).Arg(app.PostgresPool)
	app.registerComponent(users.New).Arg(app.PostgresPool)
	app.registerComponent(ldapsyncstats.New).Arg(app.PostgresPool)
	app.registerComponent(categories.New).Arg(app.PostgresPool)
	app.registerComponent(tags.New).Arg(app.PostgresPool)
	app.registerComponent(ldapsynclogs.New).Arg(app.PostgresPool)
	app.registerComponent(settings.New).Arg(app.PostgresPool)
	app.registerComponent(licenses.New).Arg(app.PostgresPool)
	app.registerComponent(productinfo.New).Arg(app.PostgresPool)
	app.registerComponent(features.New).Arg(app.PostgresPool)
	app.registerComponent(flagvariants.New).Arg(app.PostgresPool)
	app.registerComponent(rules.New).Arg(app.PostgresPool)
	app.registerComponent(featureschedules.New).Arg(app.PostgresPool)
	app.registerComponent(segmentsrepo.New).Arg(app.PostgresPool)
	app.registerComponent(ruleattributesrepo.New).Arg(app.PostgresPool)
	app.registerComponent(auditlog.New).Arg(app.PostgresPool)
	app.registerComponent(featuretagsrepo.New).Arg(app.PostgresPool)
	app.registerComponent(pending_changes.New).Arg(app.PostgresPool)
	app.registerComponent(project_approvers.New).Arg(app.PostgresPool)
	app.registerComponent(project_settings.New).Arg(app.PostgresPool)
	app.registerComponent(guard_service.New).Arg(app.PostgresPool)
	app.registerComponent(environmentsrepo.New).Arg(app.PostgresPool)
	app.registerComponent(featureparamsrepo.New).Arg(app.PostgresPool)
	app.registerComponent(dashboardrepo.New).Arg(app.PostgresPool)
	app.registerComponent(realtimerepo.New).Arg(app.PostgresPool)
	app.registerComponent(errorreports.New).Arg(app.PostgresPool)

	// Register RBAC repositories
	app.registerComponent(rbac.NewRoles).Arg(app.PostgresPool)
	app.registerComponent(rbac.NewPermissions).Arg(app.PostgresPool)
	app.registerComponent(rbac.NewMemberships).Arg(app.PostgresPool)

	// Register permissions service
	app.registerComponent(permissions.New)
	// Register feature processor service
	app.registerComponent(featuresprocessor.New).Arg(time.Second * 3)
	// Register events bus
	app.registerComponent(eventsbus.New)

	// Register use cases
	app.registerComponent(projectsusecase.New)
	app.registerComponent(ldapusecase.New)
	app.registerComponent(settingsusecase.New).Arg(app.Config.SecretKey)
	app.registerComponent(categoriesusecase.New)
	app.registerComponent(tagsusecase.New)
	app.registerComponent(featuretagsusecase.New)
	app.registerComponent(featuresusecase.New)
	app.registerComponent(flagvariantsusecase.New)
	app.registerComponent(rulesusecase.New)
	app.registerComponent(featureschedulesusecase.New)
	app.registerComponent(segmentsusecase.New)
	app.registerComponent(ruleattributesusecase.New)
	app.registerComponent(pendingchangesusecase.New)
	app.registerComponent(projectsettingsusecase.New)
	app.registerComponent(environmentsusecase.New)
	app.registerComponent(dashboardusecase.New)
	app.registerComponent(realtimechanges.New)
	app.registerComponent(rbacusecase.New)
	app.registerComponent(errorreportsusecase.New)

	app.registerComponent(email.New).Arg(&email.Config{
		SMTPHost:      app.Config.Mailer.Addr,
		Username:      app.Config.Mailer.User,
		Password:      app.Config.Mailer.Password,
		CertFile:      app.Config.Mailer.CertFile,
		KeyFile:       app.Config.Mailer.KeyFile,
		AllowInsecure: app.Config.Mailer.AllowInsecure,
		UseTLS:        app.Config.Mailer.UseTLS,
		BaseURL:       app.Config.FrontendURL,
		From:          app.Config.Mailer.From,
	})

	// Register LDAP service
	app.registerComponent(ldap.New)

	var ldapService contract.LDAPService
	if err := app.container.Resolve(&ldapService); err != nil {
		panic(err)
	}

	// Initialize SSO provider manager
	app.registerComponent(ssoprovidermanager.New)

	// Initialize SAML provider
	app.registerComponent(samlprovider.New).Arg(&samlprovider.SAMLParams{
		Name:        domain.SSOProviderNameADSaml,
		DisplayName: "Sign in with Active Directory",
		IconURL:     "",
		Config: &domain.SAMLConfig{
			Enabled:          app.Config.SAML.Enabled,
			EntityID:         app.Config.SAML.EntityID,
			CertificatePath:  app.Config.SAML.CertificatePath,
			PrivateKeyPath:   app.Config.SAML.PrivateKeyPath,
			IDPMetadataURL:   app.Config.SAML.IDPMetadataURL,
			AttributeMapping: app.Config.SAML.AttributeMapping,
			CallbackURL:      path.Join(app.Config.FrontendURL, "/api/v1/auth/sso/callback"),
			PublicRootURL:    app.Config.FrontendURL,
			SkipTLSVerify:    app.Config.SAML.SkipTLSVerify,
		},
	})

	var samlProvider *samlprovider.SAMLProvider
	if err := app.container.Resolve(&samlProvider); err != nil {
		panic(err)
	}

	app.registerComponent(usersusecase.New).Arg([]usersusecase.AuthProvider{
		ldap.NewAuthService(ldapService.(*ldap.Service)),
	})

	// Register services
	app.registerComponent(tokenizer.New).Arg(&tokenizer.ServiceParams{
		SecretKey:        []byte(app.Config.JWTSecretKey),
		AccessTTL:        app.Config.AccessTokenTTL,
		RefreshTTL:       app.Config.RefreshTokenTTL,
		ResetPasswordTTL: app.Config.ResetPasswordTTL,
	})
	app.registerComponent(ratelimiter2fa.New)

	// Register guard engine service
	app.registerComponent(guardengine.New)

	// Register API components
	app.registerComponent(apibackend.NewSecurityHandler)
	app.registerComponent(apibackend.New).Arg(app.Config)

	// Register SDK API components
	app.registerComponent(apisdk.NewSecurityHandler)
	app.registerComponent(apisdk.New)
}

func (app *App) newAPIServer() (*httpserver.Server, error) {
	cfg := app.Config.APIServer

	var restAPI generatedserver.Handler
	if err := app.container.Resolve(&restAPI); err != nil {
		return nil, fmt.Errorf("resolve REST API service component: %w", err)
	}

	var securityHandler generatedserver.SecurityHandler
	if err := app.container.Resolve(&securityHandler); err != nil {
		return nil, fmt.Errorf("resolve API security handler component: %w", err)
	}

	genServer, err := generatedserver.NewServer(restAPI, securityHandler)
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	var tokenizerSrv contract.Tokenizer
	if err := app.container.Resolve(&tokenizerSrv); err != nil {
		return nil, fmt.Errorf("resolve tokenizer service component: %w", err)
	}

	var usersSrv contract.UsersUseCase
	if err := app.container.Resolve(&usersSrv); err != nil {
		return nil, fmt.Errorf("resolve users service component: %w", err)
	}

	var permService contract.PermissionsService
	if err := app.container.Resolve(&permService); err != nil {
		return nil, fmt.Errorf("resolve permissions service component: %w", err)
	}

	// Middleware chain:
	// CORS → RAW → RequestID → Actor → Auth → API implementation
	handler := pkgmiddlewares.CORSMdw(
		middlewares.WithRawRequest(
			middlewares.RequestIDMdw(
				middlewares.ActorMdw(
					middlewares.AuthMiddleware(tokenizerSrv, usersSrv)(
						genServer,
					),
					domain.AuditActorUser,
				),
			),
		),
	)

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
	}, nil
}

func (app *App) newSDKServer() (*httpserver.Server, error) {
	cfg := app.Config.SDKServer

	var restAPI generatedsdk.Handler
	if err := app.container.Resolve(&restAPI); err != nil {
		return nil, fmt.Errorf("resolve SDK API service component: %w", err)
	}

	var securityHandler generatedsdk.SecurityHandler
	if err := app.container.Resolve(&securityHandler); err != nil {
		return nil, fmt.Errorf("resolve SDK API security handler component: %w", err)
	}

	genServer, err := generatedsdk.NewServer(restAPI, securityHandler)
	if err != nil {
		return nil, fmt.Errorf("create SDK API server: %w", err)
	}

	// Middleware chain:
	// CORS → RequestID → Actor → API implementation
	handler := pkgmiddlewares.CORSMdw(
		middlewares.RequestIDMdw(
			middlewares.ActorMdw(
				genServer,
				domain.AuditActorSDK,
			),
		),
	)

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
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

// ensureSuperuser checks if a user with the admin email exists and creates one if not.
func (app *App) ensureSuperuser(ctx context.Context) error {
	app.Logger.Info("Checking if superuser exists")

	var usersRepo contract.UsersRepository
	if err := app.container.Resolve(&usersRepo); err != nil {
		return fmt.Errorf("resolve users repository: %w", err)
	}

	// Check if user with admin email already exists
	_, err := usersRepo.GetByEmail(ctx, app.Config.AdminEmail)
	if err == nil {
		// User already exists
		app.Logger.Info("Superuser already exists")

		return nil
	}

	// Extract username from email (part before @)
	username := app.Config.AdminEmail
	for i, c := range username {
		if c == '@' {
			username = username[:i]

			break
		}
	}

	// Hash the temporary password
	passwordHash, err := passworder.PasswordHash(app.Config.AdminTmpPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Create the superuser
	userDTO := domain.UserDTO{
		Username:      username,
		Email:         app.Config.AdminEmail,
		PasswordHash:  passwordHash,
		IsSuperuser:   true,
		IsTmpPassword: true,
		IsExternal:    false,
	}

	user, err := usersRepo.Create(ctx, userDTO)
	if err != nil {
		return fmt.Errorf("create superuser: %w", err)
	}

	app.Logger.Info("Created superuser", "id", user.ID, "username", user.Username)

	return nil
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

// newWSServer creates a dedicated WebSocket server on a separate port.
func (app *App) newWSServer() (*httpserver.Server, error) {
	cfg := app.Config.WSServer

	var tokenizerSrv contract.Tokenizer
	if err := app.container.Resolve(&tokenizerSrv); err != nil {
		return nil, fmt.Errorf("resolve tokenizer service component: %w", err)
	}

	var usersSrv contract.UsersUseCase
	if err := app.container.Resolve(&usersSrv); err != nil {
		return nil, fmt.Errorf("resolve users service component: %w", err)
	}

	var rtSvc *realtimechanges.Service
	if err := app.container.Resolve(&rtSvc); err != nil {
		return nil, fmt.Errorf("resolve realtime service component: %w", err)
	}

	wsHandler := pkgmiddlewares.CORSMdw(
		middlewares.WithRawRequest(
			middlewares.RequestIDMdw(
				middlewares.ActorMdw(
					wsmiddlewares.WSAuthMiddleware(tokenizerSrv, usersSrv)(wsapi.New(rtSvc.Broadcaster())),
					domain.AuditActorUser,
				),
			),
		),
	)

	router := httprouter.New()
	router.Handler(http.MethodGet, "/api/ws", wsHandler)

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      router,
	}, nil
}
