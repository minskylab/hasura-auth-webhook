package main

import (
	log "github.com/sirupsen/logrus"

	haw "github.com/minskylab/hasura-auth-webhook"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
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

	engine, err := haw.CreateNewEngine(client, nil, false)
	if err != nil {
		log.Panicf("%+v", err)
	}

	var service server.Service
	{
		service = routes.NewService(engine)
	}

	errorCollector := make(chan error)
	runServer(conf, service, errorCollector)
	log.Errorln("exit: ", <-errorCollector)
}
