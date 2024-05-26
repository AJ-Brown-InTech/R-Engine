package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Tx represents a database transaction.
type Tx struct {
	tx *sqlx.Tx
}

// Commit commits the transaction.
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction.
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// Write executes a write operation within a transaction.
func (t *Tx) Write(ctx context.Context, table string, obj interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err = t.tx.ExecContext(ctx, query, args...)
	return err
}

// BatchWrite executes a batch write operation within a transaction.
func (t *Tx) BatchWrite(ctx context.Context, table string, objs []interface{}) error {
	for _, obj := range objs {
		err := t.Write(ctx, table, obj)
		if err != nil {
			t.Rollback()
			return err
		}
	}
	return t.Commit()
}

// BatchRead executes a batch read operation within a transaction.
func (t *Tx) BatchRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := t.tx.SelectContext(ctx, dest, query, args...)
	if err != nil {
		t.Rollback()
		return err
	}
	return nil
}
