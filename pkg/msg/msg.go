package msg

import (
	"github.com/Dissociable/Couploan/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Type is a message type
type Type string

const (
	// TypeSuccess represents a success message type
	TypeSuccess Type = "success"

	// TypeInfo represents a info message type
	TypeInfo Type = "info"

	// TypeWarning represents a warning message type
	TypeWarning Type = "warning"

	// TypeDanger represents a danger message type
	TypeDanger Type = "danger"
)

const (
	// sessionName stores the name of the session which contains flash messages
	sessionName = "msg"
)

// Success sets a success flash message
func Success(ctx fiber.Ctx, store *session.Store, message string) {
	Set(ctx, store, TypeSuccess, message)
}

// Info sets an info flash message
func Info(ctx fiber.Ctx, store *session.Store, message string) {
	Set(ctx, store, TypeInfo, message)
}

// Warning sets a warning flash message
func Warning(ctx fiber.Ctx, store *session.Store, message string) {
	Set(ctx, store, TypeWarning, message)
}

// Danger sets a danger flash message
func Danger(ctx fiber.Ctx, store *session.Store, message string) {
	Set(ctx, store, TypeDanger, message)
}

// Set adds a new flash message of a given type into the session storage.
// Errors will be logged and not returned.
func Set(ctx fiber.Ctx, store *session.Store, typ Type, message string) {
	if sess, err := getSession(ctx, store); err == nil {
		flashesSess := sess.Get("flashes")
		var flashes map[string][]string
		if flashesSess == nil {
			flashes = make(map[string][]string)
		} else {
			flashes = flashesSess.(map[string][]string)
		}
		_, ok := flashes[string(typ)]
		if !ok {
			flashes[string(typ)] = []string{message}
		} else {
			flashes[string(typ)] = append(flashes[string(typ)], message)
		}
		sess.Set("flashes", flashes)
		save(ctx, sess)
	}
}

// Get gets flash messages of a given type from the session storage.
// Errors will be logged and not returned.
func Get(ctx fiber.Ctx, store *session.Store, typ Type) []string {
	var msgs []string

	if sess, err := getSession(ctx, store); err == nil {
		flashesSess := sess.Get("flashes")
		var flashes map[string][]string
		if flashesSess != nil {
			flashes = flashesSess.(map[string][]string)
		}
		if flashes == nil {
			return nil
		}
		if flashMessages, ok := flashes[string(typ)]; ok && len(flashMessages) > 0 {
			delete(flashes, string(typ))
			sess.Set("flashes", flashes)
			save(ctx, sess)
			msgs = append(msgs, flashMessages...)
		}
	}

	return msgs
}

// getSession gets the flash message session
func getSession(ctx fiber.Ctx, store *session.Store) (*session.Session, error) {
	l := util.GetLoggerFromFiberCtx(ctx)
	if store == nil {
		s := util.GetSessionStoreFromFiberCtx(ctx)
		if s == nil {
			l.Error("no session found in context")
			return nil, errors.New("no session found in context")
		}
		r, err := s.Get(ctx)
		if err != nil || r == nil {
			return nil, errors.New("no session found in context")
		}
		return r, nil
	}
	sess, err := store.Get(ctx)
	if err != nil {
		if l != nil {
			l.Error("cannot load flash message session", zap.Error(err))
		}
	}
	return sess, err
}

// save saves the flash message session
func save(ctx fiber.Ctx, sess *session.Session) {
	if err := sess.Save(); err != nil {
		l := util.GetLoggerFromFiberCtx(ctx)
		if l != nil {
			l.Error("cannot save flash message session", zap.Error(err))
		}
	}
}
