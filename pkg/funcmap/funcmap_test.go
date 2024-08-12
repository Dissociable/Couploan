package funcmap

import (
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/gofiber/fiber/v3"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFuncMap(t *testing.T) {
	f := NewFuncMap(fiber.New())
	assert.NotNil(t, f["hasField"])
	assert.NotNil(t, f["link"])
	assert.NotNil(t, f["file"])
	assert.NotNil(t, f["url"])
}

func TestHasField(t *testing.T) {
	type example struct {
		name string
	}
	var e example
	f := new(funcMap)
	assert.True(t, f.hasField(e, "name"))
	assert.False(t, f.hasField(e, "abcd"))
}

func TestLink(t *testing.T) {
	f := new(funcMap)

	link := string(f.link("/abc", "Text", "/abc"))
	expected := `<a class="is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)

	link = string(f.link("/abc", "Text", "/abc", "first", "second"))
	expected = `<a class="first second is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)

	link = string(f.link("/abc", "Text", "/def"))
	expected = `<a class="" href="/abc">Text</a>`
	assert.Equal(t, expected, link)
}

func TestFile(t *testing.T) {
	f := new(funcMap)

	file := f.file("test.png")
	expected := fmt.Sprintf("/%s/test.png?v=%s", config.StaticPrefix, CacheBuster)
	assert.Equal(t, expected, file)
}
