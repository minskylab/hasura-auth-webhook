package engine

import (
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/ent"
)

// Engine wraps all related to the platform (it's the core engine).
type Engine struct {
	*ent.Client
	Auth *auth.AuthManager
}
