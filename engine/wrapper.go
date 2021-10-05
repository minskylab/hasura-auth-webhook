package engine

import (
	"context"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/server"
	log "github.com/sirupsen/logrus"
)

// CreateNewEngine creates a new Engine instance.
func CreateNewEngine(client *ent.Client, authInstance *auth.AuthManager, conf *config.Config, bootstrap bool) (*Engine, error) {
	if bootstrap {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		for _, r := range conf.Roles {
			role, err := client.Role.Create().SetName(r.Name).Save(context.Background())
			if err != nil {
				log.Warnf("failed creating role: %v", err)
				continue
			}

			for _, u := range r.Users {
				hashed, err := helpers.HashPassword(u.Password)
				if err != nil {
					log.Warnf("failed hashing password: %v", err)
					continue
				}

				_, err = client.User.Create().SetEmail(u.Email).SetHashedPassword(hashed).AddRoles(role).Save(context.Background())
				if err != nil {
					log.Warnf("failed creating role: %v", err)
				}

				// user
			}

		}
	}

	engine := &Engine{
		PublicService:   server.NewPublicServer(client, authInstance, conf),
		InternalService: server.NewInternalServer(client, authInstance, conf),
		Auth:            authInstance,
	}

	return engine, nil
}
