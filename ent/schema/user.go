package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	uuid "github.com/satori/go.uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.NewV4).Unique(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),

		field.String("email"),
		field.String("password").Sensitive(),
		field.String("token").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
