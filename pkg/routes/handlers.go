package routes

import (
	"fmt"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"net/http"
)

var handlers []Handler

// Handler handles one or more HTTP routes
type Handler interface {
	// Routes allows for self-registration of HTTP routes on the router
	Routes(g fiber.Router)

	// Init provides the service container to initialize
	Init(*services.Container) error
}

// Register registers a handler
func Register(h Handler) {
	handlers = append(handlers, h)
}

// GetHandlers returns all handlers
func GetHandlers() []Handler {
	return handlers
}

// fail is a helper to fail a request by returning a 500 error and logging the error
func fail(err error, log string) error {
	// The error handler will handle logging
	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("%s: %v", log, err))
}
