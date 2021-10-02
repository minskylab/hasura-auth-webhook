package engine

import (
	"context"
	"log"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/ent"
)

// CreateNewEngine creates a new Engine instance.
func CreateNewEngine(client *ent.Client, authInstance *auth.AuthManager, createSchema bool) (*Engine, error) {
	if createSchema {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	engine := &Engine{
		Client: client,
		Auth:   authInstance,
	}

	return engine, nil
}
