package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
	"github.com/minskylab/hasura-auth-webhook/engine"
	"github.com/minskylab/hasura-auth-webhook/server"
	"github.com/minskylab/hasura-auth-webhook/services/routes"
)

func main() {
	log.SetLevel(log.DebugLevel)

	conf := config.NewConfig()

	client, err := db.OpenDBClient(conf)
	if err != nil {
		log.Panicf("%+v", err)
	}
	defer client.Close()

	secretSource := auth.RawSecret([]byte(conf.JwtAccessKeySecret), []byte(conf.JwtRefreshKeySecret))
	authManager, err := auth.New(secretSource)
	if err != nil {
		log.Panicf("%+v", err)
	}

	engine, err := engine.CreateNewEngine(client, authManager, true)
	if err != nil {
		log.Panicf("%+v", err)
	}

	var service server.Service
	{
		service = routes.NewService(engine)
	}

	errorCollector := make(chan error)

	go checkForSignals(errorCollector)
	go runServer(conf, service, errorCollector)

	log.Errorln("exit: ", <-errorCollector)
}
