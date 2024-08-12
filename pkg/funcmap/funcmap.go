package funcmap

import (
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/templates"
	sprig "github.com/go-task/slim-sprig"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/teris-io/shortid"
	"html/template"
	"math/rand/v2"
	"reflect"
	"strings"
)

var (
	sid, _ = shortid.New(1, shortid.DefaultABC, rand.Uint64())
	// CacheBuster stores a random string used as a cache buster for static files.
	CacheBuster = sid.MustGenerate()
)

type funcMap struct {
	web *fiber.App
}

// NewFuncMap provides a template function map
func NewFuncMap(web *fiber.App) template.FuncMap {
	fm := &funcMap{web: web}

	// See http://masterminds.github.io/sprig/ for all provided funcs
	funcs := sprig.FuncMap()

	// Include all the custom functions
	funcs["hasField"] = fm.hasField
	funcs["file"] = fm.file
	funcs["link"] = fm.link
	funcs["url"] = fm.url
	funcs["urlUiName"] = fm.urlUiName
	funcs["YN"] = fm.YN
	funcs["JsonIndent"] = fm.JsonIndent

	return funcs
}

// hasField checks if an interface contains a given field
func (fm *funcMap) hasField(v any, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// file appends a cache buster to a given filepath so it can remain cached until the app is restarted
func (fm *funcMap) file(filepath string) string {
	return fmt.Sprintf("/%s/%s?v=%s", config.StaticPrefix, filepath, CacheBuster)
}

// link outputs HTML for a link element, providing the ability to dynamically set the active class
func (fm *funcMap) link(url, text, currentPath string, classes ...string) template.HTML {
	if currentPath == url {
		classes = append(classes, "is-active")
	}

	html := fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, strings.Join(classes, " "), url, text)
	return template.HTML(html)
}

// url generates a URL from a given route name and optional parameters
func (fm *funcMap) url(routeName string, params ...any) string {
	r := fm.web.GetRoute(routeName)
	var paramsMap fiber.Map
	paramsMap = make(fiber.Map)
	for i, param := range r.Params {
		paramsMap[param] = fmt.Sprintf("%v", params[i])
	}
	result, err := fiber.NewDefaultCtx(fm.web).GetRouteURL(routeName, paramsMap)
	if err != nil {
		return ""
	}
	return result
}

// url generates a URL from a given route name and optional parameters
func (fm *funcMap) urlUiName(uiName string, params ...any) string {
	routeName := string(templates.PageNameMap[templates.UiName(uiName)])
	r := fm.web.GetRoute(routeName)
	var paramsMap fiber.Map
	paramsMap = make(fiber.Map)
	for i, param := range r.Params {
		paramsMap[param] = fmt.Sprintf("%v", params[i])
	}
	result, err := fiber.NewDefaultCtx(fm.web).GetRouteURL(routeName, paramsMap)
	if err != nil {
		return ""
	}
	return result
}

// YN returns "Yes" or "No"
func (fm *funcMap) YN(v any) string {
	if v == nil {
		return ""
	}

	vNew := false

	switch vc := v.(type) {
	case bool:
		vNew = vc
	case *bool:
		vNew = *vc
	}

	if vNew {
		return "Yes"
	}
	return "No"
}

func (fm *funcMap) JsonIndent(v any) string {
	r, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}

	return string(r)
}
