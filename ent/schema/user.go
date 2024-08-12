package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Annotations(entsql.Default("uuid_generate_v4()")),
		field.String("name").NotEmpty(),
		field.Text("contact").NotEmpty(),
		field.Enum("role").Values("user", "admin").Default("user"),
		field.String("key").Unique().Nillable().MinLen(32).MaxLen(32),
		// balance field, in cents
		field.Int("balance").Default(0),
		// expires_at field
		field.Time("expires_at").Nillable().Optional().Default(nil),
		field.Time("created_at").Default(time.Now).Annotations(entsql.Default("now()")),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// index on name
		index.Fields("name"),

		// index on the key
		index.Fields("key"),

		// index on the key and balance
		index.Fields("key", "balance"),

		// index on the key and expires_at
		index.Fields("key", "expires_at"),

		// index on the key, expires_at, balance
		index.Fields("key", "expires_at", "balance"),

		// index on contact
		index.Fields("contact"),
	}
}
