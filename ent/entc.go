//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	if err := entc.Generate(
		"./schema",
		&gen.Config{Features: []gen.Feature{gen.FeatureVersionedMigration, gen.FeatureUpsert}},
		entc.Extensions(&UserExtension{}),
	); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}

// UserExtension implements entc.Extension.
type UserExtension struct {
	entc.DefaultExtension
}

func (*UserExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(gen.NewTemplate("user_extension").ParseFiles("templates/user_extension.tmpl")),
	}
}
