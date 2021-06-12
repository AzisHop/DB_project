package db

import (
	"github.com/jackc/pgx"
	_ "github.com/jmoiron/sqlx"
)

type Postgres struct {
	database *pgx.ConnPool
}

func NewDb() (*Postgres, error) {
	conf := pgx.ConnConfig{
		User:                 "postgres",
		Database:             "postgres",
		Password:             "admin",
		PreferSimpleProtocol: false,
	}

	poolConf := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	database, err := pgx.NewConnPool(poolConf)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		database: database,
	}, nil
}

func (postgres *Postgres) GetPs() *pgx.ConnPool {
	return postgres.database
}

func (postgres *Postgres) Close() error {
	postgres.database.Close()
	return nil
}