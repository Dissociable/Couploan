package tests

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Dissociable/Couploan/ent"
	"github.com/brianvoe/gofakeit/v7"
)

// NewContextTest creates a new Echo context for tests using an HTTP test request and response recorder
func NewContextTest(a *fiber.App, url string, tests func(ctx fiber.Ctx)) (*http.Response, error) {
	if a == nil {
		a = fiber.New()
		a.Get(
			"/*", func(c fiber.Ctx) error {
				tests(c)
				return nil
			},
		)
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)
	return a.Test(req, -1)
}

// NewContextTestWithHeaders creates a new Echo context for tests using an HTTP test request and response recorder
func NewContextTestWithHeaders(
	a *fiber.App, url string, headers http.Header, tests func(ctx fiber.Ctx),
) (*http.Response, error) {
	if a == nil {
		a = fiber.New()
		a.Get(
			"/*", func(c fiber.Ctx) error {
				tests(c)
				return nil
			},
		)
	}

	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header = headers
	return a.Test(req, -1)
}

// CreateUser creates a random user entity
func CreateUser(orm *ent.Client) (*ent.User, error) {
	seed := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), rand.Intn(1000000))
	return orm.User.
		Create().
		SetKey(gofakeit.Password(true, true, true, false, false, 32)).
		SetName(fmt.Sprintf("Test User %s", seed)).
		SetContact(fmt.Sprintf("testuser-%s@localhost.localhost", seed)).
		Save(context.Background())
}
