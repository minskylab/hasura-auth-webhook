package engine

import (
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/services"
)

// Engine wraps all related to the platform (it's the core engine).
type Engine struct {
	Auth *auth.AuthManager

	PublicService   services.PublicService
	InternalService services.InternalService
}
