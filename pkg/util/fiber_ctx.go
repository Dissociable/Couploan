package util

import (
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/pkg/context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.uber.org/zap"
)

func GetContainerFromFiberCtx(c fiber.Ctx) any {
	l := c.Locals("container")
	if l == nil {
		return nil
	}
	return l
}

func GetLoggerFromFiberCtx(c fiber.Ctx) *zap.Logger {
	l := c.Locals("logger")
	if l == nil {
		return nil
	}
	lz := l.(*zap.Logger)
	lz = lz.
		With(zap.String("ip", c.IP())).
		With(zap.String("method", c.Method())).
		With(zap.String("path", c.Path())).
		With(zap.String("request_id", requestid.FromContext(c)))

	return lz
}

func GetUserFromFiberCtx(c fiber.Ctx) *ent.User {
	u := c.Locals(context.AuthenticatedUserKey)
	if u == nil {
		return nil
	}
	return u.(*ent.User)
}

func GetSessionStoreFromFiberCtx(c fiber.Ctx) *session.Store {
	l := c.Locals("session_store")
	if l == nil {
		return nil
	}
	return l.(*session.Store)
}

func GetSessionFromFiberCtx(c fiber.Ctx) *session.Session {
	store := GetSessionStoreFromFiberCtx(c)
	if store == nil {
		return nil
	}
	s, err := store.Get(c)
	if err != nil {
		return nil
	}
	return s
}

func GetSessionIDFromFiberCtx(c fiber.Ctx) string {
	id := c.Locals("session_id")
	if id == nil {
		return ""
	}
	return id.(string)
}
