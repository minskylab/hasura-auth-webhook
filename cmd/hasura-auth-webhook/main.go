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
	defaultConfigfileV1 = "config.yaml"
	defaultConfigfileV2 = "config.v2.yaml"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	configVersion := flag.Int("config-version", 2, "enable file config v2")
	configFile := flag.String("config", "", "config file path")

	flag.Parse()

	if debug != nil && *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var conf *config.Config

	if *configVersion == 1 {
		if *configFile == "" {
			conf, _ = config.NewConfig(defaultConfigfileV1)
		} else {
			conf, _ = config.NewConfig(*configFile)
		}
	} else if *configVersion == 2 {
		var confV2 *config.Config2
		if *configFile == "" {
			confV2, _ = config.NewConfigV2(defaultConfigfileV2)
		} else {
			confV2, _ = config.NewConfigV2(*configFile)
		}

		conf = config.ConfigV2ToConfigV1(confV2)
	}

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
