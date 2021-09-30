package routes

import (
	haw "github.com/minskylab/hasura-auth-webhook"
	"github.com/minskylab/hasura-auth-webhook/server"
)

type service struct {
	engine *haw.Engine
}

func NewService(engine *haw.Engine) server.Service {
	return service{
		engine: engine,
	}
}
