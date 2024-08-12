package templates

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type (
	Layout          string
	Page            string
	UiName          string
	UiNameMapType   map[Page]UiName
	PageNameMapType map[UiName]Page
)

const (
	LayoutMain Layout = "main"
	LayoutAuth Layout = "auth"
	LayoutHTMX Layout = "htmx"
)

const (
	PageError Page = "error"
	PageHome  Page = "home"
)

var UiNameMap UiNameMapType
var PageNameMap PageNameMapType

func init() {
	UiNameMap = make(UiNameMapType)
	UiNameMap[PageHome] = "Index"

	PageNameMap = make(PageNameMapType)
	for k, v := range UiNameMap {
		PageNameMap[v] = k
	}
}

//go:embed *
var templates embed.FS

// Get returns a file system containing all templates via embed.FS
func Get() embed.FS {
	return templates
}

// GetOS returns a file system containing all templates which will load the files directly from the operating system.
// This should only be used for local development in order to facilitate live reloading.
func GetOS() fs.FS {
	// Gets the complete templates directory path
	// This is needed in case this is called from a package outside of main, such as within tests
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	p := filepath.Join(filepath.Dir(d), "templates")
	return os.DirFS(p)
}

// GetOSPath returns the path to the templates directory
func GetOSPath() string {
	// Gets the complete templates directory path
	// This is needed in case this is called from a package outside of main, such as within tests
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), "templates")
}
