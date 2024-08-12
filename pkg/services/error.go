package services

import (
	"errors"
	"github.com/Dissociable/Couploan/pkg/context"
	"github.com/Dissociable/Couploan/pkg/page"
	"github.com/Dissociable/Couploan/pkg/util"
	"github.com/Dissociable/Couploan/templates"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"net/http"
)

type Error struct {
	*TemplateRenderer
}

func (e *Error) Page(err error, ctx fiber.Ctx) {
	if ctx.Response().SkipBody || context.IsCanceledError(err) {
		return
	}

	// Determine the error status code
	code := http.StatusInternalServerError
	var he *fiber.Error
	if errors.As(err, &he) {
		code = he.Code
	}

	// Log the error
	logger := util.GetLoggerFromFiberCtx(ctx)
	switch {
	case code >= 500:
		logger.Error(err.Error())
	case code >= 400:
		logger.Warn(err.Error())
	}

	// Render the error page
	p := page.New(ctx)
	p.Layout = templates.LayoutMain
	p.Name = templates.PageError
	p.Title = http.StatusText(code)
	p.StatusCode = code

	if err = e.RenderPage(ctx, p); err != nil {
		logger.Error(
			"failed to render error page",
			zap.Error(err),
		)
	}
}
