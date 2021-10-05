package engine

import (
	"context"
	"log"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/server"
)

// CreateNewEngine creates a new Engine instance.
func CreateNewEngine(client *ent.Client, authInstance *auth.AuthManager, conf *config.Config, createSchema bool) (*Engine, error) {
	if createSchema {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	engine := &Engine{
		PublicService:   server.NewPublicServer(client, authInstance, conf),
		InternalService: server.NewInternalServer(client, authInstance, conf),
		Auth:            authInstance,
	}

	return engine, nil
}
