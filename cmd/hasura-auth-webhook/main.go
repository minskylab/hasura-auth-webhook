package main

import (
	"github.com/sirupsen/logrus"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
	"github.com/minskylab/hasura-auth-webhook/engine"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	conf := config.NewDefaultConfig()

	client, err := db.OpenDBClient(conf)
	if err != nil {
		logrus.Panicf("%+v", err)
	}
	defer client.Close()

	secretSource := auth.RawSecret(
		[]byte(conf.JWT.AccessSecret),
		[]byte(conf.JWT.RefreshSecret),
	)

	// authManager, err := auth.New(secretSource, c.isAnonymousAllowed)
	var anonymous *auth.AnonymousRole
	for _, r := range conf.Roles {
		if r.IsAnonymous {
			anonymous = &auth.AnonymousRole{
				Name: r.Name,
			}
			break
		}
	}
	authManager, err := auth.New(secretSource, anonymous)
	if err != nil {
		logrus.Panicf("%+v", err)
	}

	engine, err := engine.CreateNewEngine(client, authManager, conf, true)
	if err != nil {
		logrus.Panicf("%+v", err)
	}

	signalErr := engine.Launch()

	if err := <-signalErr; err != nil {
		logrus.Panicf("%+v", err)
	}
}
