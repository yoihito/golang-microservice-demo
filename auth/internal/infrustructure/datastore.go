package infrustructure

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Datastore struct {
	db *sql.DB
}

func NewDatastore(connStr string) (*Datastore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Datastore{db: db}, nil
}

func (d *Datastore) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}
