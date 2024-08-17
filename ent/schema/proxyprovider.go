package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProxyProvider holds the schema definition for the ProxyProvider entity.
type ProxyProvider struct {
	ent.Schema
}

// Fields of the ProxyProvider.
func (ProxyProvider) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("username"),
		field.String("password"),
		field.String("service_type"),
	}
}

// Edges of the ProxyProvider.
func (ProxyProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("proxy", Proxy.Type).
			Ref("proxyProvider"),
	}
}

// Indexes of the ProxyProvider.
func (ProxyProvider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("name", "service_type"),
		index.Fields("name", "username"),
		index.Fields("name", "service_type", "username", "password").Unique(),
	}
}
