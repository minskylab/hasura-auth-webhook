package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/db"
	"github.com/minskylab/hasura-auth-webhook/engine"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	configV2 := flag.Bool("config2", false, "enable file config v2")

	flag.Parse()

	if debug != nil && *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var conf *config.Config

	if configV2 != nil && *configV2 {
		logrus.Info("Using config v2")
		confV2, _ := config.NewConfigV2("config.v2.yaml")
		conf = config.ConfigV2ToConfigV1(confV2)
	} else {
		conf, _ = config.NewConfig("config.yaml")
	}
	// helpers.Log(conf)

	client, err := db.OpenDBClient(conf)
	if err != nil {
		logrus.Panicf("%+v", err)
	}
	defer client.Close()

	secretSource := auth.RawSecret(
		[]byte(conf.JWT.AccessSecret),
		[]byte(conf.JWT.RefreshSecret),
	)

	// var anonymous *auth.AnonymousRole
	// if conf.Anonymous != nil {
	// 	anonymous = &auth.AnonymousRole{
	// 		Name: conf.Anonymous.Name,
	// 	}
	// }

	// authManager, err := auth.New(secretSource, anonymous)
	authManager, err := auth.New(secretSource)
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
