package db

import "github.com/pkg/errors"

// ErrConfigDB ...
var ErrConfigDB = errors.New("invalid db config, please verify if your config file is defined correctly")
