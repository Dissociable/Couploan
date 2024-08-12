package routes

import (
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent/user"
	"github.com/Dissociable/Couploan/pkg/middleware"
	"github.com/Dissociable/Couploan/pkg/middleware/keyauth"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.uber.org/zap"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// BuildRouter builds the router
func BuildRouter(c *services.Container) {
	// Static files with proper cache control
	// funcmap.File() should be used in templates to append a cache key to the URL in order to break cache
	// after each server restart
	c.Web.Static(
		config.StaticPrefix, config.StaticDir, fiber.Static{
			Compress:      true,
			CacheDuration: 24 * time.Hour, // refresh in 1 day.
		},
	)

	// Initialize default config
	// check if favicon file exists
	if _, err := os.Stat(config.StaticDir + "/favicon.png"); err == nil {
		c.Web.Use(
			favicon.New(
				favicon.Config{
					File: "./static/favicon.png",
					URL:  "/favicon.png",
				},
			),
		)
	}

	recoverConfig := recover.ConfigDefault
	recoverConfig.StackTraceHandler = func(fctx fiber.Ctx, e any) {
		c.Logger.Error(
			"Crash recovering, stack trace:",
			zap.Any("error", e),
			zap.String("request_id", requestid.FromContext(fctx)),
			zap.Any("request", fctx.Request()),
			zap.ByteString("stack", debug.Stack()),
		)
	}
	recoverConfig.EnableStackTrace = true
	c.Web.Use(recover.New(recoverConfig))

	corsCfg := cors.ConfigDefault
	corsCfg.Next = func(ctx fiber.Ctx) bool {
		return false
	}
	c.Web.Use(cors.New(corsCfg))

	c.Web.Use(helmet.New())

	c.Web.Use(
		compress.New(
			compress.Config{
				Level: compress.LevelBestSpeed, // 1
			},
		),
	)

	// Non-static file route group
	g := c.Web.Group("")

	g.Use(requestid.New())

	// Error handler
	// err := errorHandler{Controller: ctr}
	// c.Web.HTTPErrorHandler = err.Get

	gManage := g.Group("/manage")
	store := session.New()
	c.WebSession = store
	gManage.Use(middleware.SessionToContext(c.WebSession))
	g.Use(middleware.SessionToContext(c.WebSession))
	gManage.Use(csrf.New())
	gApi := g.Group("/api")
	gApiV1 := gApi.Group("/v1")

	// Auth middleware
	keyAuthMwConfig := keyauth.ConfigDefault
	keyAuthMwConfig.KeyLookup = "header:Authorization|query:key"
	keyAuthMwConfig.ErrorHandler = func(ctx fiber.Ctx, err error) error {
		if !strings.HasPrefix(ctx.Route().Path, "/api/") {
			return ctx.Next()
		}
		return keyauth.ConfigDefault.ErrorHandler(ctx, err)
	}
	keyAuthMwConfig.Validator = func(ctx fiber.Ctx, key string) (bool, error) {
		u, err := c.ORM.User.Query().
			Where(
				user.And(
					user.Key(key),
					user.BalanceGTE(c.Config.Pricing.PricePerRegister),
				),
			).First(ctx.Context())
		if err != nil {
			return false, err
		}
		err = c.Auth.Login(ctx, u.ID)
		if err != nil {
			c.Logger.Error("failed to login the user via AuthClient", zap.String("uuid", u.ID.String()), zap.Error(err))
			return false, err
		}
		return true, nil
	}
	gApiV1.Use(keyauth.New(keyAuthMwConfig))
	gApiV1.Use(middleware.LoadAuthenticatedUser(c.Auth))
	gApiV1.Use(middleware.RequireAuthentication())
	g.Use(keyauth.New(keyAuthMwConfig))
	g.Use(middleware.LoadAuthenticatedUser(c.Auth))
	// Example routes
	navRoutes(c, g)
	userRoutes(c, gApiV1)
}

func navRoutes(c *services.Container, g fiber.Router) {
	// Initialize and register all handlers
	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			c.Logger.Panic("failed to init handler", zap.Error(err))
		}

		h.Routes(g)
	}
	// g.Get(
	// 	"/",
	// 	func(ctx fiber.Ctx) error {
	// 		return ctx.SendString("Hello, World!")
	// 	},
	// )
}

func userRoutes(c *services.Container, g fiber.Router) {
	g.Get(
		"/",
		func(ctx fiber.Ctx) error {
			return ctx.SendString("Hello, World!")
		},
	)
	// g.Get("/ping", v1.Pong)
}
