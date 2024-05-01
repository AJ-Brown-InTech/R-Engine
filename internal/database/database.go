package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"github.com/jmoiron/sqlx"
)

// DB represents a database connection.
type DB struct {
	db *sqlx.DB
}

// Tx represents a database transaction.
type Tx struct {
	tx *sqlx.Tx
}

// NewDB initializes a new database connection.
func NewDB(conn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	
	return &DB{db: db}, nil
}

// NewTx starts a new transaction.
func (d *DB) NewTx(ctx context.Context) (*Tx, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// ObjectToDatabase converts a struct to column names and values.
func ObjectToDatabase(obj interface{}) ([]string, []interface{}, error) {
	var columns []string
	var values []interface{}

	objValue := reflect.ValueOf(obj)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	objType := objValue.Type()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		tag := field.Tag
		dbTag := tag.Get("db")

		if dbTag != "" && dbTag != "-" {
			fieldValue := objValue.Field(i)

			if fieldValue.Kind() != reflect.Ptr || !fieldValue.IsNil() {
				columns = append(columns, dbTag)
				values = append(values, fieldValue.Interface())
			}
		}
	}
	return columns, values, nil
}

// Read executes a read operation.
func (d *DB) Read(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := d.db.GetContext(ctx, dest, query, args...)
	return err
}

// BatchRead executes a batch read operation.
func (d *DB) BatchRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := d.db.SelectContext(ctx, dest, query, args...)
	return err
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

	_, err = d.db.ExecContext(ctx, query, args...)
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

	_, err = d.db.ExecContext(ctx, query, append(args, conditionValue)...)
	return err
}

// Delete executes a delete operation.
func (d *DB) Delete(ctx context.Context, table string, conditionColumn string, conditionValue interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, conditionColumn)

	_, err := d.db.ExecContext(ctx, query, conditionValue)
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

	err := t.Commit()
	return err
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

// BatchRead executes a batch read operation within a transaction.
func (t *Tx) BatchRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := t.tx.SelectContext(ctx, dest, query, args...)
	if err != nil {
		t.Rollback()
		return err
	}
	return nil
}
