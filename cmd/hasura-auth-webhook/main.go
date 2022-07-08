package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
	"github.com/minskylab/hasura-auth-webhook/engine"
)

const (
	defaultConfigFilename = "config.yaml"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	configFile := flag.String("config-file", "", "config file path")

	flag.Parse()

	if debug != nil && *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var conf *config.Config

	if *configFile == "" {
		conf, _ = config.NewConfig(defaultConfigFilename)
	} else {
		conf, _ = config.NewConfig(*configFile)
	}

	client, err := db.OpenDBClient(conf)
	if err != nil {
		logrus.Panicf("%+v", err)
	}
	defer client.Close()

	secretSource := auth.RawSecret(
		[]byte(conf.Providers.Email.JWT.AccessSecret),
		[]byte(conf.Providers.Email.JWT.RefreshSecret),
	)

	// var anonymous *auth.AnonymousRole
	// if conf.Anonymous != nil {
	// 	anonymous = &auth.AnonymousRole{
	// 		Name: conf.Anonymous.Name,
	// 	}
	// }

	// authManager, err := auth.New(secretSource, anonymous)
	authManager, err := auth.New(secretSource, conf)
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
