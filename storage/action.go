package storage

import (
	"context"
	"fmt"
	"strings"
)

// Read executes a read operation.
func (d *DB) Read(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	return d.Db.GetContext(ctx, dest, query, args...)
}

// BatchRead executes a batch read operation.
func (d *DB) BatchRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	return d.Db.SelectContext(ctx, dest, query, args...)
}

// Write executes a write operation.
func (d *DB) Write(ctx context.Context, table string, obj interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err = d.Db.ExecContext(ctx, query, args...)
	return err
}

// Update executes an update operation.
func (d *DB) Update(ctx context.Context, table string, obj interface{}, conditionColumn string, conditionValue interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	setClause := make([]string, len(columns))
	for i, column := range columns {
		setClause[i] = fmt.Sprintf("%s = $%d", column, i+1)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", table, strings.Join(setClause, ", "), conditionColumn, len(columns)+1)
	_, err = d.Db.ExecContext(ctx, query, append(args, conditionValue)...)
	return err
}

// Delete executes a delete operation.
func (d *DB) Delete(ctx context.Context, table string, conditionColumn string, conditionValue interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, conditionColumn)
	_, err := d.Db.ExecContext(ctx, query, conditionValue)
	return err
}
