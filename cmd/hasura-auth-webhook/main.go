package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
	"github.com/minskylab/hasura-auth-webhook/engine"
)

func main() {
	log.SetLevel(log.DebugLevel)

	conf := config.NewDefaultConfig()

	client, err := db.OpenDBClient(conf)
	if err != nil {
		log.Panicf("%+v", err)
	}
	defer client.Close()

	secretSource := auth.RawSecret(
		[]byte(conf.JWT.AccessSecret),
		[]byte(conf.JWT.RefreshSecret),
	)

	authManager, err := auth.New(secretSource)
	if err != nil {
		log.Panicf("%+v", err)
	}

	engine, err := engine.CreateNewEngine(client, authManager, conf, true)
	if err != nil {
		log.Panicf("%+v", err)
	}

	signalErr := engine.Launch()

	if err := <-signalErr; err != nil {
		log.Panicf("%+v", err)
	}
}
