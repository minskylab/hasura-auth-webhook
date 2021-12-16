package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/pkg/errors"
)

func NewConfig(filepath string) (*Config, error) {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.TagName = "yaml"
	})

	if err := config.LoadFiles(filepath); err != nil {
		return nil, errors.WithStack(err)
	}

	c := new(Config)

	return c, config.BindStruct("", c)
}

func NewConfigV2(filepath string) (*Config2, error) {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.TagName = "yaml"
	})

	if err := config.LoadFiles(filepath); err != nil {
		return nil, errors.WithStack(err)
	}

	conf := new(Config2)

	return conf, config.BindStruct("", conf)
}
