package hasuraauthwebhook

import (
	"context"
	"log"

	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

// Engine wraps all related to fairpay platform (it's the core engine).
type Engine struct {
	*ent.Client
	Auth *union.Union
}

// CreateNewEngine creates a new Fairpay Engine instance.
func CreateNewEngine(client *ent.Client, authInstance *union.Union, createSchema bool) (*Engine, error) {
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
