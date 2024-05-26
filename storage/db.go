package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" 
)

// DB represents a database connection.
type DB struct {
	Db *sqlx.DB
}

// NewDB initializes a new database connection.
func NewDB(conn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	return &DB{Db: db}, nil
}

// NewTx starts a new transaction.
func (d *DB) NewTx(ctx context.Context) (*Tx, error) {
	tx, err := d.Db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.Db.Close()
}
