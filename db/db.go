package db

import "github.com/jmoiron/sqlx"

type Postgres struct {
	database *sqlx.DB
}

func NewDb(connect string) (*Postgres, error) {
	database, err := sqlx.Connect("pqx", connect)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		database: database,
	}, nil
}

func (postgres *Postgres) GetPs() *sqlx.DB {
	return postgres.database
}

func (postgres *Postgres) Close() error {
	err := postgres.database.Close()
	return err
}