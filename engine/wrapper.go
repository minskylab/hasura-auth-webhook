package engine

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/minskylab/hasura-auth-webhook/ent/role"

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

		adminRole := config.Role{
			Name:   "admin",
			Users:  conf.Admin.Users,
			Public: false,
		}

		anonymousRole := config.Role{
			Name:   conf.Anonymous.Name,
			Users:  []config.User{},
			Public: true,
		}

		for _, r := range conf.Roles {
			adminRole.Children = append(adminRole.Children, r.Name)
		}

		mapEntityRoles := make(map[string]*ent.Role)

		allRoles := append(conf.Roles, adminRole, anonymousRole)

		for _, r := range allRoles {
			rolEntity, err := client.Role.Create().SetName(r.Name).SetPublic(r.Public).Save(context.Background())
			if err != nil {
				// log.Warnf("failed creating role: %v", err)
				// continue
				rolEntity, err = client.Role.Query().Where(role.Name(r.Name)).First(context.Background())
				if err != nil {
					log.Warnf("failed creating role: %v", err)
					continue
				}
			}

			mapEntityRoles[r.Name] = rolEntity

			for _, u := range r.Users {
				hashed, err := helpers.HashPassword(u.Password)
				if err != nil {
					log.Warnf("failed hashing password: %v", err)
					continue
				}

				_, err = client.User.Create().SetEmail(u.Email).SetHashedPassword(hashed).AddRoles(rolEntity).Save(context.Background())
				if err != nil {
					log.Warnf("failed creating role: %v", err)
				}
			}

			// spew.Dump(rolEntity)
		}

		for _, r := range allRoles {
			parents := []*ent.Role{}
			children := []*ent.Role{}

			rolEntity := mapEntityRoles[r.Name]

			if r.Parent != nil {
				if parentRol, ok := mapEntityRoles[*r.Parent]; !ok {
					log.Warnf("unknown parent")
				} else {
					parents = append(parents, parentRol)
				}
			}
			for _, p := range r.Parents {
				if parentRol, ok := mapEntityRoles[p]; !ok {
					log.Warnf("unknown parent")
				} else {
					parents = append(parents, parentRol)
				}
			}

			if r.Child != nil {
				if childRol, ok := mapEntityRoles[*r.Child]; !ok {
					log.Warnf("unknown child")
				} else {
					children = append(children, childRol)
				}
			}
			for _, c := range r.Children {
				if childRol, ok := mapEntityRoles[c]; !ok {
					log.Warnf("unknown child")
				} else {
					children = append(children, childRol)
				}
			}

			spew.Dump(rolEntity.Name)

			spew.Dump("PARENTS: ", parents)
			spew.Dump("CHILDREN: ", children)
			fmt.Println()

			// rolEntity, err := rolEntity.Update().AddParents(parents...).AddChildren(children...).Save(context.Background())
			//rolEntity, err := rolEntity.Update().AddParents(parents...).Save(context.Background())
			//if err != nil {
			//	log.Warnf("failed updating roles: %v", err)
			//}

			_, err := rolEntity.Update().AddChildren(children...).Save(context.Background())
			if err != nil {
				log.Warnf("failed updating roles: %v", err)
			}

			if rolEntity != nil {
				rolEntity, _ = client.Role.Query().WithParents().WithChildren().Where(role.ID(rolEntity.ID)).Only(context.Background())
				spew.Dump(rolEntity)
				spew.Dump(rolEntity.Edges)
				mapEntityRoles[r.Name] = rolEntity
			}

			fmt.Println("----------------------------------\n")
		}
	}

	engine := &Engine{
		PublicService:   server.NewPublicServer(client, authInstance, conf),
		InternalService: server.NewInternalServer(client, authInstance, conf),
		Auth:            authInstance,
	}

	return engine, nil
}
