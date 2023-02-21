package storage

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
)

type PsqlStorage struct {
	db *pgxpool.Pool
}

func NewPsqlStorage(cfg *configs.Config) (*PsqlStorage, error) {
	st := &PsqlStorage{}

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	st.db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = st.doMigrate(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return st, nil
} // todo: use another name when create context with timeout in functions

func (st *PsqlStorage) doMigrate(dsn string) error {
	m, err := migrate.New("file://migrations/postgres", dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
