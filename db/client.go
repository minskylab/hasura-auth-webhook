package db

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

func openEntClient(c *config.Config) (*ent.Client, error) {
	db, err := dburl.Open(c.DB.URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var driver *entsql.Driver = nil

	if strings.HasPrefix(c.DB.URL, "postgres") {
		driver = entsql.OpenDB(dialect.Postgres, db)
	} else if strings.HasPrefix(c.DB.URL, "sqlite") {
		driver = entsql.OpenDB(dialect.SQLite, db)
	}

	if driver == nil {
		return nil, fmt.Errorf("invalid sql dialect. '%s' url not supported", c.DB.URL)
	}

	return ent.NewClient(ent.Driver(driver)), nil
}

// OpenDBClient create a new Ent Client with Postgres or Sqlite Connection
func OpenDBClient(c *config.Config) (*ent.Client, error) {
	return openEntClient(c)
}
