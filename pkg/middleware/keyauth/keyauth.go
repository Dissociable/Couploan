// Special thanks to Echo: https://github.com/labstack/echo/blob/master/middleware/key_auth.go
package keyauth

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/url"
	"strings"
)

// When there is no request of the key thrown ErrMissingOrMalformedAPIKey
var ErrMissingOrMalformedAPIKey = errors.New("missing or malformed API Key")

const (
	query  = "query"
	form   = "form"
	param  = "param"
	cookie = "cookie"
	header = "header"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Init config
	cfg := configDefault(config...)

	// Initialize
	var extractor []func(c fiber.Ctx) (string, error)
	multiParts := strings.Split(cfg.KeyLookup, "|")
	for _, multiPart := range multiParts {
		parts := strings.Split(multiPart, ":")
		switch parts[0] {
		case header:
			extractor = append(extractor, keyFromHeader(parts[1], cfg.AuthScheme))
		case query:
			extractor = append(extractor, keyFromQuery(parts[1]))
		case form:
			extractor = append(extractor, keyFromForm(parts[1]))
		case param:
			extractor = append(extractor, keyFromParam(parts[1]))
		case cookie:
			extractor = append(extractor, keyFromCookie(parts[1]))
		}
	}

	// Return middleware handler
	return func(c fiber.Ctx) error {
		// Filter request to skip middleware
		if (cfg.Next != nil && cfg.Next(c)) || len(extractor) == 0 {
			return c.Next()
		}

		// Extract and verify key
		var key string
		var err error
		for _, extractor := range extractor {
			key, err = extractor(c)
			if err != nil && !errors.Is(err, ErrMissingOrMalformedAPIKey) {
				return cfg.ErrorHandler(c, err)
			}
			if err == nil && len(key) > 0 {
				break
			}
		}
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		valid, err := cfg.Validator(c, key)

		if err == nil && valid {
			c.Locals(cfg.ContextKey, key)
			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}

// keyFromHeader returns a function that extracts api key from the request header.
func keyFromHeader(header, authScheme string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > 0 && l == 0 {
			return auth, nil
		}
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrMissingOrMalformedAPIKey
	}
}

// keyFromQuery returns a function that extracts api key from the query string.
func keyFromQuery(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.Query(param)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromForm returns a function that extracts api key from the form.
func keyFromForm(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.FormValue(param)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromParam returns a function that extracts api key from the url param string.
func keyFromParam(param string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key, err := url.PathUnescape(c.Params(param))
		if err != nil {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}

// keyFromCookie returns a function that extracts api key from the named cookie.
func keyFromCookie(name string) func(c fiber.Ctx) (string, error) {
	return func(c fiber.Ctx) (string, error) {
		key := c.Cookies(name)
		if key == "" {
			return "", ErrMissingOrMalformedAPIKey
		}
		return key, nil
	}
}
