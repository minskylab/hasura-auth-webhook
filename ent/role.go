// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/minskylab/hasura-auth-webhook/ent/role"

	uuid "github.com/satori/go.uuid"
)

// Role is the model entity for the Role schema.
type Role struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Public holds the value of the "public" field.
	Public bool `json:"public,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RoleQuery when eager-loading is set.
	Edges RoleEdges `json:"edges"`
}

// RoleEdges holds the relations/edges for other nodes in the graph.
type RoleEdges struct {
	// Users holds the value of the users edge.
	Users []*User `json:"users,omitempty"`
	// Children holds the value of the children edge.
	Children []*Role `json:"children,omitempty"`
	// Parents holds the value of the parents edge.
	Parents []*Role `json:"parents,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// UsersOrErr returns the Users value or an error if the edge
// was not loaded in eager-loading.
func (e RoleEdges) UsersOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Users, nil
	}
	return nil, &NotLoadedError{edge: "users"}
}

// ChildrenOrErr returns the Children value or an error if the edge
// was not loaded in eager-loading.
func (e RoleEdges) ChildrenOrErr() ([]*Role, error) {
	if e.loadedTypes[1] {
		return e.Children, nil
	}
	return nil, &NotLoadedError{edge: "children"}
}

// ParentsOrErr returns the Parents value or an error if the edge
// was not loaded in eager-loading.
func (e RoleEdges) ParentsOrErr() ([]*Role, error) {
	if e.loadedTypes[2] {
		return e.Parents, nil
	}
	return nil, &NotLoadedError{edge: "parents"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Role) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case role.FieldPublic:
			values[i] = new(sql.NullBool)
		case role.FieldName:
			values[i] = new(sql.NullString)
		case role.FieldCreatedAt, role.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case role.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Role", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Role fields.
func (r *Role) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case role.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				r.ID = *value
			}
		case role.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				r.CreatedAt = value.Time
			}
		case role.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				r.UpdatedAt = value.Time
			}
		case role.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				r.Name = value.String
			}
		case role.FieldPublic:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field public", values[i])
			} else if value.Valid {
				r.Public = value.Bool
			}
		}
	}
	return nil
}

// QueryUsers queries the "users" edge of the Role entity.
func (r *Role) QueryUsers() *UserQuery {
	return (&RoleClient{config: r.config}).QueryUsers(r)
}

// QueryChildren queries the "children" edge of the Role entity.
func (r *Role) QueryChildren() *RoleQuery {
	return (&RoleClient{config: r.config}).QueryChildren(r)
}

// QueryParents queries the "parents" edge of the Role entity.
func (r *Role) QueryParents() *RoleQuery {
	return (&RoleClient{config: r.config}).QueryParents(r)
}

// Update returns a builder for updating this Role.
// Note that you need to call Role.Unwrap() before calling this method if this Role
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Role) Update() *RoleUpdateOne {
	return (&RoleClient{config: r.config}).UpdateOne(r)
}

// Unwrap unwraps the Role entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Role) Unwrap() *Role {
	tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Role is not a transactional entity")
	}
	r.config.driver = tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Role) String() string {
	var builder strings.Builder
	builder.WriteString("Role(")
	builder.WriteString(fmt.Sprintf("id=%v", r.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(r.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(r.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(r.Name)
	builder.WriteString(", public=")
	builder.WriteString(fmt.Sprintf("%v", r.Public))
	builder.WriteByte(')')
	return builder.String()
}

// Roles is a parsable slice of Role.
type Roles []*Role

func (r Roles) config(cfg config) {
	for _i := range r {
		r[_i].config = cfg
	}
}
