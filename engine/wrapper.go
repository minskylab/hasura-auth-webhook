package engine

import (
	"context"

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

		// anonymousRole := config.Role{
		// 	Name:   conf.Anonymous.Name,
		// 	Users:  []config.User{},
		// 	Public: true,
		// }

		for _, r := range conf.Roles {
			adminRole.Children = append(adminRole.Children, r.Name)
		}

		mapEntityRoles := make(map[string]*ent.Role)

		// allRoles := append(conf.Roles, adminRole, anonymousRole)
		allRoles := append(conf.Roles, adminRole) //, anonymousRole)

		{
			mapRoles := make(map[string]*config.Role)

			for i := 0; i < len(allRoles); i++ {
				mapRoles[allRoles[i].Name] = &allRoles[i]
			}

			for _, v := range mapRoles {
				for _, c := range v.Children {
					if !contains(mapRoles[c].Parents, v.Name) {
						mapRoles[c].Parents = append(mapRoles[c].Parents, v.Name)
					}
				}
			}
		}

		for _, r := range allRoles {
			rolEntity, err := client.Role.Create().SetName(r.Name).SetPublic(r.Public).Save(context.Background())
			if err != nil {
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
		}

		for _, r := range allRoles {
			var parents []*ent.Role

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

			rolEntity, err := rolEntity.Update().AddParents(parents...).Save(context.Background())
			if err != nil {
				log.Warnf("failed updating roles: %v", err)
			}

			if rolEntity != nil {
				mapEntityRoles[r.Name] = rolEntity
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
