package db

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/steven.t/doppler-example/config"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrExists         = errors.New("exists")
	ErrNoUpdate       = errors.New("nothing to update")
	ErrNoPrimaryEmail = errors.New("no primary email was provided")
)

type Instance struct {
	connStr      string
	ensureSchema bool
	suffix       string
	db           *pgxpool.Pool
}

func Configure(ensureSchema bool, suffix string, conf *config.Config) *Instance {
	i := Instance{
		ensureSchema: ensureSchema,
		suffix:       suffix,
	}
	ctx := context.Background()
	i.connStr = conf.Secrets["ROACH_CONN"]
	conn, err := pgxpool.New(ctx, i.connStr)
	if err != nil {
		log.Fatal("unable to establist connection to cockroachdb")
	}
	i.db = conn

	if ensureSchema {
		if err := i.ensureSchemas(ctx, i.suffix); err != nil {
			log.Fatal("unable to run the schema check: %w", err)
		}
	}
	return &i
}

func (i *Instance) SetNewConnection(ctx context.Context, connStr string) error {
	i.db.Close()
	i.connStr = connStr
	conn, err := pgxpool.New(ctx, i.connStr)
	if err != nil {
		return err
	}
	i.db = conn

	if i.ensureSchema {
		if err := i.ensureSchemas(ctx, i.suffix); err != nil {
			return err
		}
	}
	return nil
}
