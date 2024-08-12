package middleware

import (
	"fmt"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/pkg/context"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/Dissociable/Couploan/pkg/util"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// LoadAuthenticatedUser loads the authenticated user, if one, and stores in context
func LoadAuthenticatedUser(authClient *services.AuthClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		l := util.GetLoggerFromFiberCtx(c)
		u, err := authClient.GetAuthenticatedUser(c)
		switch err.(type) {
		case *ent.NotFoundError:
			l.Warn("auth user not found")
		case services.NotAuthenticatedError:
		case nil:
			c.Locals(context.AuthenticatedUserKey, u)
			// l.Debug("auth user loaded in to context", zap.String("user_id", u.ID.String()))
		default:
			l.Warn("error querying for authenticated user", zap.Error(err))
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				fmt.Sprintf("error querying for authenticated user: %v", err),
			)
		}

		return c.Next()
	}
}

func RequireAdminUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		u := c.Locals(context.AuthenticatedUserKey)
		if u == nil {
			return fiber.ErrUnauthorized
		}
		if !u.(*ent.User).IsAdmin() {
			return fiber.ErrForbidden
		}
		return c.Next()
	}
}

// RequireAuthentication requires that the user be authenticated in order to proceed
func RequireAuthentication() fiber.Handler {
	return func(c fiber.Ctx) error {
		if u := c.Locals(context.AuthenticatedUserKey); u == nil {
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}
}

// RequireNoAuthentication requires that the user not be authenticated in order to proceed
func RequireNoAuthentication() fiber.Handler {
	return func(c fiber.Ctx) error {
		if u := c.Locals(context.AuthenticatedUserKey); u != nil {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
