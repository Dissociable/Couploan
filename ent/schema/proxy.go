package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Proxy holds the schema definition for the Proxy entity.
type Proxy struct {
	ent.Schema
}

// Fields of the Proxy.
func (Proxy) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("HTTP", "HTTPS", "SOCKS4", "SOCKS5", "SOCKS5H"),
		field.String("ip"),
		field.Uint16("port"),
		field.String("username").Nillable(),
		field.String("password").Nillable(),
		field.Bool("rotating").Default(false),
	}
}

// Edges of the Proxy.
func (Proxy) Edges() []ent.Edge {
	return nil
}

func (Proxy) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ip", "port"),
		index.Fields("ip", "port", "username", "password").Unique(),
	}
}
