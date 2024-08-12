package services

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/logger"
	"github.com/Dissociable/Couploan/pkg/funcmap"
	"github.com/Dissociable/Couploan/proxstore"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	// Required by ent
	"ariga.io/atlas-go-sdk/atlasexec"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests
type Container struct {
	// Validator stores a validator
	Validator *Validator

	// Web stores the web framework
	Web *fiber.App

	// WebSession stores the web's session store
	WebSession *session.Store

	// Config stores the application configuration
	Config *config.Config

	// Cache contains the cache client
	Cache *CacheClient

	// Database stores the connection to the database
	Database *sql.DB

	// ORM stores a client to the ORM
	ORM *ent.Client

	// Auth stores an authentication client
	Auth *AuthClient

	// TemplateRenderer stores a service to easily render and cache templates
	TemplateRenderer *TemplateRenderer

	// Tasks stores the Task client
	Tasks *TaskClient

	Logger *zap.Logger

	ProxyStore *proxstore.ProxStore[tls_client.HttpClient]
}

// NewContainer creates and initializes a new Container
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initLogger()
	c.initValidator()
	c.initWeb()
	c.initCache()
	c.initDatabase()
	c.initORM()
	c.initAuth()
	c.initTemplateRenderer()
	c.initTasks()
	c.initProxyStore()
	return c
}

// Shutdown shuts the Container down and disconnects all connections
func (c *Container) Shutdown() error {
	if c.Tasks != nil {
		if err := c.Tasks.Close(); err != nil {
			return err
		}
	}
	if c.Cache != nil {
		if err := c.Cache.Close(); err != nil {
			return err
		}
	}
	if c.ORM != nil {
		if err := c.ORM.Close(); err != nil {
			return err
		}
	}
	if c.Database != nil {
		if err := c.Database.Close(); err != nil {
			return err
		}
	}
	if c.Logger != nil {
		_ = c.Logger.Sync()
	}

	return nil
}

// initConfig initializes configuration
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg
}

// initValidator initializes the validator
func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

// initWeb initializes the web framework
func (c *Container) initWeb() {
	prxyHeader := c.Config.HTTP.ProxyHeader
	c.Web = fiber.New(
		fiber.Config{
			AppName:                 c.Config.App.Name,
			JSONEncoder:             json.Marshal,
			JSONDecoder:             json.Unmarshal,
			ProxyHeader:             prxyHeader,
			EnableTrustedProxyCheck: true,
			TrustedProxies:          c.Config.HTTP.TrustedProxies,
			ErrorHandler: func(ctx fiber.Ctx, err error) error {
				if !strings.HasPrefix(ctx.Route().Path, "/api/") {
					errPage := Error{TemplateRenderer: c.TemplateRenderer}
					errPage.Page(err, ctx)
					return nil
				}
				code := fiber.StatusInternalServerError
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}
				return ctx.Status(code).SendString(fiber.ErrInternalServerError.Message)
			},
		},
	)

	c.Web.Use(
		func(fctx fiber.Ctx) error {
			fctx.Locals("container", c)
			return fctx.Next()
		},
	)

	c.Web.Use(
		func(fctx fiber.Ctx) error {
			fctx.Locals("logger", c.Logger)
			return fctx.Next()
		},
	)

	// Configure logging
	switch c.Config.App.Environment {
	case config.EnvProduction:
		log.SetLevel(log.LevelWarn)
	default:
		log.SetLevel(log.LevelDebug)
		// c.Web.Logger.SetLevel(log.DEBUG)
	}
}

func (c *Container) StartServer(gracefulShutdownContext context.Context) error {
	c.Web.Use(
		func(c fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
		},
	)
	cfg := fiber.ListenConfig{GracefulContext: gracefulShutdownContext}
	if c.Config.HTTP.TLS.Enabled {
		cfg.CertFile = c.Config.HTTP.TLS.Certificate
		cfg.CertKeyFile = c.Config.HTTP.TLS.Key
	}
	return c.Web.Listen(
		fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port), cfg,
	)
}

// initCache initializes the cache
func (c *Container) initCache() {
	var err error
	if c.Cache, err = NewCacheClient(c.Config); err != nil {
		panic(err)
	}
}

// initDatabase initializes the database
// If the environment is set to test, the test database will be used and will be dropped, recreated and migrated
func (c *Container) initDatabase() {
	var err error

	getAddr := func(dbName string) string {
		return fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s",
			c.Config.Database.User,
			c.Config.Database.Password,
			c.Config.Database.Hostname,
			c.Config.Database.Port,
			dbName,
		)
	}

	c.Database, err = sql.Open("pgx", getAddr(c.Config.Database.Database))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Check if this is a test environment
	if c.Config.App.Environment == config.EnvTest {
		// Drop the test database, ignoring errors in case it doesn't yet exist
		_, _ = c.Database.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", c.Config.Database.TestDatabase))

		// Create the test database
		if _, err = c.Database.Exec(fmt.Sprintf("CREATE DATABASE %s", c.Config.Database.TestDatabase)); err != nil {
			panic(fmt.Sprintf("failed to create test database: %v", err))
		}

		// Connect to the test database
		if err = c.Database.Close(); err != nil {
			panic(fmt.Sprintf("failed to close database connection: %v", err))
		}
		c.Database, err = sql.Open("pgx", getAddr(c.Config.Database.TestDatabase))

		if _, err := c.Database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
			panic("failed to add uuid extension")
		}

		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
	}
}

// initORM initializes the ORM
func (c *Container) initORM() {
	drv := entsql.OpenDB(dialect.Postgres, c.Database)
	c.ORM = ent.NewClient(ent.Driver(drv))

	// Install uuid-ossp extension in PostgreSQL so we can generate UUID in postgresql
	_, err := c.Database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		panic(fmt.Sprintf("failed to install postgresql extensions: %v", err))
	}

	dbName := c.Config.Database.Database
	// Check if this is a test environment
	if c.Config.App.Environment == config.EnvTest {
		dbName = c.Config.Database.TestDatabase
	}
	u, err := url.Parse(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?search_path=public&sslmode=disable",
			c.Config.Database.User,
			c.Config.Database.Password,
			c.Config.Database.Hostname,
			c.Config.Database.Port,
			dbName,
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	// Define the execution context, supplying a migration directory
	// and potentially an `atlas.hcl` configuration file using `atlasexec.WithHCL`.
	migrationsDir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			ent.EmbeddedMigrations,
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize migrations directory: %v", err))
	}
	// atlasexec works on a temporary directory, so we need to close it
	defer migrationsDir.Close()

	dirPath := migrationsDir.Path()
	// Initialize the client.
	client, err := atlasexec.NewClient(dirPath, "atlas")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}
	// Run `atlas migrate apply` on a SQLite database under /tmp.
	res, err := client.MigrateApply(
		context.Background(), &atlasexec.MigrateApplyParams{
			Env:             "",
			ConfigURL:       "",
			Context:         nil,
			DirURL:          "file://migrations/migrate/migrations/",
			AllowDirty:      true,
			URL:             u.String(),
			RevisionsSchema: "",
			BaselineVersion: "",
			TxMode:          "none",
			ExecOrder:       "",
			Amount:          0,
			DryRun:          false,
			Vars:            nil,
		},
	)
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations\n", len(res.Applied))
}

// initAuth initializes the authentication client
func (c *Container) initAuth() {
	c.Auth = NewAuthClient(c.Config, c.ORM)
}

// initTemplateRenderer initializes the template renderer
func (c *Container) initTemplateRenderer() {
	c.TemplateRenderer = NewTemplateRenderer(c.Config, c.Cache, funcmap.NewFuncMap(c.Web))
}

// initTasks initializes the Task client
func (c *Container) initTasks() {
	c.Tasks = NewTaskClient(c.Config)
}

// initLogger initializes the logger
func (c *Container) initLogger() {
	if c.Config.App.LogLevel != "" {
		logger.LogLevel = strings.ToLower(c.Config.App.LogLevel)
	} else {
		logger.LogLevel = "debug"
	}
	logger.LokiConfig = c.Config.Loki
	c.Logger = logger.Get()
}

func (c *Container) initProxyStore() {
	tls_client.DefaultTimeoutSeconds = 20
	options := proxstore.Options{}
	optionsCreateHttpClient := proxstore.OptionsCreateHttpClient[tls_client.HttpClient]{
		Creator: func(proxy *proxstore.Proxy[tls_client.HttpClient]) (hc tls_client.HttpClient, err error) {
			opts := []tls_client.HttpClientOption{
				tls_client.WithTimeoutSeconds(20),
				tls_client.WithClientProfile(profiles.Chrome_124),
			}
			if !proxy.IsEmpty() && !proxy.IsDirect() {
				opts = append(opts, tls_client.WithProxyUrl(proxy.String()))
			}
			hc, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), opts...)
			return
		},
	}
	c.ProxyStore = proxstore.NewWithOptions[tls_client.HttpClient](&options, &optionsCreateHttpClient)
}
