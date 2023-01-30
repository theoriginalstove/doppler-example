package db

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/avocagrow/doppler-example/config"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrExists         = errors.New("exists")
	ErrNoUpdate       = errors.New("nothing to update")
	ErrNoPrimaryEmail = errors.New("no primary email was provided")
)

type Instance struct {
	db *pgxpool.Pool
}

func Configure(setupSchema bool, suffix string, conf *config.Config) *Instance {
	i := Instance{}
	ctx := context.Background()
	connString := conf.Secrets["ROACHCONNSTR"]
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("unable to establist connection to cockroachdb")
	}
	i.db = conn

	if setupSchema {
		if err := i.ensureSchemas(ctx, suffix); err != nil {
			log.Fatal("unable to run the schema check")
		}
	}
	return &i
}
