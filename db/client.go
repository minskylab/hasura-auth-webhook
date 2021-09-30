package db

import (
	"fmt"

	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/pkg/errors"
)

func connectionString(c *config.Config) (string, error) {
	host := c.DB_HOST
	port := c.DB_PORT
	user := c.DB_USER
	dbName := c.DB_DATABASE
	password := c.DB_PASS

	if host == "" || port == "" || user == "" || dbName == "" || password == "" {
		return "", ErrConfigDB
	}

	config := "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"

	return fmt.Sprintf(config, host, port, user, dbName, password), nil
}

func openPostgresClient(c *config.Config) (*ent.Client, error) {
	source, err := connectionString(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ent.Open("postgres", source)
}

// OpenPostgresDB create a new Ent Client with Postgres Connection
func OpenDBClient(c *config.Config) (*ent.Client, error) {
	return openPostgresClient(c)
}
