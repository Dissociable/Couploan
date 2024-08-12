//go:build ignore

package main

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	config "github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/ent"
	"github.com/Dissociable/Couploan/ent/migrate"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"time"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	_ "entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/Dissociable/Couploan/ent"
)

const (
	dir = "ent/migrate/migrations"
)

func main() {
	ctx := context.Background()
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	// Create a local migration directory able to understand Atlas migration file format for replay.
	if err = os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("creating migration directory: %v", err)
	}
	dir, err := atlas.NewLocalDir(dir)
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	}
	if len(os.Args) < 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	testDb := false
	if len(os.Args) == 3 && os.Args[2] == "test" {
		testDb = true
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	conn := fmt.Sprintf(
		"postgres://postgres:root@localhost:5437/%s?search_path=public&sslmode=disable",
		cfg.Database.Database,
	)
	if testDb {
		conn = fmt.Sprintf(
			"postgres://postgres:root@localhost:5437/%s?search_path=public&sslmode=disable",
			cfg.Database.TestDatabase,
		)
	}
	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	drv := entsql.OpenDB("postgres", db)
	client := ent.NewClient(ent.Driver(drv))
	// Install uuid-ossp extension in PostgreSQL so we can generate UUID in postgresql
	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		panic(fmt.Sprintf("failed to install postgresql extensions: %v", err))
	}
	err = client.Schema.NamedDiff(ctx, os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
	time.Sleep(1 * time.Second)
}
