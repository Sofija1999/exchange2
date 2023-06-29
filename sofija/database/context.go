package database

import (
	"context"
)

// ctxKey is an unexported type for context keys defined in this package.
// This prevents collisions with context keys defined in other packages.
type ctxKey int

// txKey is the key for *sqlx.Tx values in Contexts. It is unexported;
// clients use database.NewContext and database.FromContext instead of using this key directly.
var txKey ctxKey

// NewContext returns a new Context with provided *sqlx.Tx.
func NewContext(ctx context.Context, tx *Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// FromContext returns the *sqlx.Tx value stored in ctx, if any.
func FromContext(ctx context.Context) (*Tx, bool) {
	tx, ok := ctx.Value(txKey).(*Tx)

	return tx, ok
}
