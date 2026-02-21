package tx

import (
	"context"
	"database/sql"
)

// Runner runs a function inside a transaction; rollback on error.
type Runner struct {
	DB *sql.DB
}

// Run executes fn in a transaction. If fn returns an error, the transaction is rolled back.
func (r *Runner) Run(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
