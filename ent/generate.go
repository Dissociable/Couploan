package ent

import "embed"

//go:generate go run entc.go

//go:embed migrate/migrations/*
var EmbeddedMigrations embed.FS
