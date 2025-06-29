package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBTx is an interface that both *sqlx.DB and *sqlx.Tx satisfy
type DBTx interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	Rebind(query string) string
}

// DB is an interface that represents a database connection with transaction capabilities
type DB interface {
	DBTx
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}
