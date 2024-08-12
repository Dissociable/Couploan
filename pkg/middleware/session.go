package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func SessionToContext(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("session_store", store)
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Locals("session_id", sess.ID())
		return c.Next()
	}
}
