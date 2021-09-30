package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
	"github.com/RussellLuo/kok/pkg/httpoption"
	"github.com/rs/cors"

	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/server"
)

func runServer(conf *config.Config, service server.Service, errCollector chan error) {

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCollector <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		serverURI := conf.HOST + ":" + conf.PORT
		helpers.PrintLogo(serverURI)
		s := server.NewHTTPRouter(service, httpcodec.NewDefaultCodecs(nil), httpoption.RequestValidators())
		serv := cors.AllowAll().Handler(s)
		errCollector <- http.ListenAndServe(serverURI, serv)
	}()

}
