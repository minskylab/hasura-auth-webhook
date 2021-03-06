package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/pkg/errors"
)

func initConfig(filepaths ...string) error {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.TagName = "yaml"
	})

	if err := config.LoadFiles(filepaths...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewConfig(filepaths ...string) (*Config, error) {
	if err := initConfig(filepaths...); err != nil {
		return nil, errors.WithStack(err)
	}

	conf := new(Config)

	return conf, config.BindStruct("", conf)
}
