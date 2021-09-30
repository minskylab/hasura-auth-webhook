package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/server"
)

func runServer(conf *config.Config, service server.Service, errCollector chan error) {
	serverURI := conf.HOST + ":" + conf.PORT

	router := mux.NewRouter()
	server.SetupRoutes(router, service)
	handler := cors.AllowAll().Handler(router)

	helpers.PrintLogo(serverURI)
	errCollector <- http.ListenAndServe(serverURI, handler)
}
