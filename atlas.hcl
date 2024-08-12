env {
    name = atlas.env
    url = var.url
    extension "uuid-ossp" {
      comment = "UUID functions"
    }
    migration {
        // URL where the migration directory resides.
        dir = "file://ent/migrate/migrations"
    }
}