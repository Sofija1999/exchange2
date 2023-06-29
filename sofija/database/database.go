package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Row = sqlx.Row
type Rows = sqlx.Rows
type Tx = sqlx.Tx

type DB struct {
	db *sqlx.DB
}

func NewDB(cfg config.DatabaseConfig) (*DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name, cfg.Schema)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, errors.Wrap(err, "open postgres conn")
	}

	return &DB{db: db}, nil
}

func (db *DB) TxContext(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	ctx = NewContext(ctx, tx)

	defer func() {
		err = TxHandler(tx, err, recover())
	}()

	return fn(ctx)
}

func (db *DB) TxContextSerializable(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	_, err = tx.Exec(`set transaction isolation level SERIALIZABLE`)
	if err != nil {
		tx.Rollback()
		return
	}

	ctx = NewContext(ctx, tx)

	defer func() {
		err = TxHandler(tx, err, recover())
	}()

	return fn(ctx)
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	_, ok := FromContext(ctx)
	if ok {
		return nil, errors.New("transaction already started in context")
	}

	return db.db.BeginTxx(ctx, nil)
}

func (db *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.executor(ctx).GetContext(ctx, dest, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return db.executor(ctx).QueryxContext(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	return db.executor(ctx).QueryRowxContext(ctx, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.executor(ctx).ExecContext(ctx, query, args...)
}

func (db *DB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.executor(ctx).PrepareContext(ctx, query)
}

func (db *DB) executor(ctx context.Context) dbExecutor {
	if tx, ok := FromContext(ctx); ok {
		return tx
	}

	return db.db
}

// dbExecutor is an interface that is implemeneted
// both by *sqlx.Database and *sqlx.Tx so they can be used interchangeably.
type dbExecutor interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
