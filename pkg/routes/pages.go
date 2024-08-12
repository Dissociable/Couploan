package routes

import (
	"github.com/Dissociable/Couploan/pkg/page"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/Dissociable/Couploan/templates"
	"github.com/gofiber/fiber/v3"
	"html/template"
)

const (
	routeNameHome = string(templates.PageHome)
)

type (
	menu struct {
		uiNames []string
	}

	Pages struct {
		*services.TemplateRenderer
		menu *menu
	}

	post struct {
		Title string
		Body  string
	}

	aboutData struct {
		ShowCacheWarning bool
		FrontendTabs     []aboutTab
		BackendTabs      []aboutTab
	}

	aboutTab struct {
		Title string
		Body  template.HTML
	}
)

func init() {
	public := new(Pages)
	public.menu = &menu{
		uiNames: []string{
			string(templates.UiNameMap[templates.PageHome]),
		},
	}
	Register(public)
}

func (h *Pages) Init(c *services.Container) error {
	h.TemplateRenderer = c.TemplateRenderer
	return nil
}

func (h *Pages) Routes(g fiber.Router) {
	g.Get("/", h.Home).Name(routeNameHome)
}

func (h *Pages) Home(ctx fiber.Ctx) error {
	p := page.New(ctx)
	p.Layout = templates.LayoutMain
	p.Name = templates.PageHome
	p.Metatags.Description = "Welcome to the homepage."
	p.Metatags.Keywords = []string{}
	if h.menu != nil {
		p.Menu = h.menu.uiNames
	}

	return h.RenderPage(ctx, p)
}
