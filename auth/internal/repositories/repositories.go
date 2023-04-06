package repositories

import (
	"context"
	"database/sql"
)

type Datastore interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
