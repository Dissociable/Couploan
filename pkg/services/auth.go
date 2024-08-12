package services

import (
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/ent/user"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	// authSessionKeyUserID stores the key used to store the user ID in the session
	authSessionKeyUserID = "user_id"

	// authSessionKeyAuthenticated stores the key used to store the authentication status in the session
	authSessionKeyAuthenticated = "authenticated"
)

// NotAuthenticatedError is an error returned when a user is not authenticated
type NotAuthenticatedError struct{}

// Error implements the error interface.
func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}

// AuthClient is the client that handles authentication requests
type AuthClient struct {
	config *config.Config
	orm    *ent.Client
}

// NewAuthClient creates a new authentication client
func NewAuthClient(cfg *config.Config, orm *ent.Client) *AuthClient {
	return &AuthClient{
		config: cfg,
		orm:    orm,
	}
}

// Login logs in a user of a given ID
func (c *AuthClient) Login(ctx fiber.Ctx, userID uuid.UUID) error {
	ctx.Locals(authSessionKeyUserID, userID)
	ctx.Locals(authSessionKeyAuthenticated, true)
	return nil
}

// Logout logs the requesting user out
func (c *AuthClient) Logout(ctx fiber.Ctx) error {
	ctx.Locals(authSessionKeyUserID, nil)
	ctx.Locals(authSessionKeyAuthenticated, false)
	return nil
}

// GetAuthenticatedUserID returns the authenticated user's ID, if the user is logged in
func (c *AuthClient) GetAuthenticatedUserID(ctx fiber.Ctx) (uuid.UUID, error) {
	i := ctx.Locals(authSessionKeyUserID)
	if i != nil {
		return i.(uuid.UUID), nil
	}

	return uuid.UUID{}, NotAuthenticatedError{}
}

// GetAuthenticatedUser returns the authenticated user if the user is logged in
func (c *AuthClient) GetAuthenticatedUser(ctx fiber.Ctx) (*ent.User, error) {
	if userID, err := c.GetAuthenticatedUserID(ctx); err == nil {
		return c.orm.User.Query().
			Where(user.ID(userID)).
			Only(ctx.Context())
	}

	return nil, NotAuthenticatedError{}
}
