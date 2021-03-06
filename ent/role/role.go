// Code generated by entc, DO NOT EDIT.

package role

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	// Label holds the string label denoting the role type in the database.
	Label = "role"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldPublic holds the string denoting the public field in the database.
	FieldPublic = "public"
	// EdgeUsers holds the string denoting the users edge name in mutations.
	EdgeUsers = "users"
	// EdgeChildren holds the string denoting the children edge name in mutations.
	EdgeChildren = "children"
	// EdgeParents holds the string denoting the parents edge name in mutations.
	EdgeParents = "parents"
	// Table holds the table name of the role in the database.
	Table = "roles"
	// UsersTable is the table that holds the users relation/edge. The primary key declared below.
	UsersTable = "user_roles"
	// UsersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UsersInverseTable = "users"
	// ChildrenTable is the table that holds the children relation/edge. The primary key declared below.
	ChildrenTable = "role_children"
	// ParentsTable is the table that holds the parents relation/edge. The primary key declared below.
	ParentsTable = "role_children"
)

// Columns holds all SQL columns for role fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
	FieldPublic,
}

var (
	// UsersPrimaryKey and UsersColumn2 are the table columns denoting the
	// primary key for the users relation (M2M).
	UsersPrimaryKey = []string{"user_id", "role_id"}
	// ChildrenPrimaryKey and ChildrenColumn2 are the table columns denoting the
	// primary key for the children relation (M2M).
	ChildrenPrimaryKey = []string{"role_id", "parent_id"}
	// ParentsPrimaryKey and ParentsColumn2 are the table columns denoting the
	// primary key for the parents relation (M2M).
	ParentsPrimaryKey = []string{"role_id", "parent_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultPublic holds the default value on creation for the "public" field.
	DefaultPublic bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
