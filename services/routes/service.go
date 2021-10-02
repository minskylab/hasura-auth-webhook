package routes

import (
	"github.com/minskylab/hasura-auth-webhook/engine"
	"github.com/minskylab/hasura-auth-webhook/server"
)

type service struct {
	engine *engine.Engine
}

func NewService(engine *engine.Engine) server.Service {
	return service{
		engine: engine,
	}
}
