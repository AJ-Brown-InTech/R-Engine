
package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func InitDatabase(conn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Info("Established a successful database connection.")
	return db, nil
}

// Convert struct to a slice of column names strings and column values strings
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

// CRUD Operations

// Read operation
func Read(ctx context.Context, db *sqlx.DB, query string, dest interface{}, args ...interface{}) error {
	err := db.GetContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Batch Read operation
func BatchRead(ctx context.Context, db *sqlx.DB, query string, dest interface{}, args ...interface{}) error {
	err := db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Write operation
func Write(ctx context.Context, db *sqlx.DB, table string, obj interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update operation
func Update(ctx context.Context, db *sqlx.DB, table string, obj interface{}, conditionColumn string, conditionValue interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	setClause := make([]string, len(columns))
	for i, column := range columns {
		setClause[i] = fmt.Sprintf("%s = $%d", column, i+1)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", table, strings.Join(setClause, ", "), conditionColumn, len(columns)+1)

	_, err = db.ExecContext(ctx, query, append(args, conditionValue)...)
	if err != nil {
		return err
	}

	return nil
}

// Delete operation
func Delete(ctx context.Context, db *sqlx.DB, table string, conditionColumn string, conditionValue interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, conditionColumn)

	_, err := db.ExecContext(ctx, query, conditionValue)
	if err != nil {
		return err
	}

	return nil
}

// Batch Write operation
func BatchWrite(ctx context.Context, db *sqlx.DB, table string, objs []interface{}) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		err := Write(ctx, tx, table, obj)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Transaction Operations

// StartTransaction starts a new transaction
func StartTransaction(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CommitTransaction commits a transaction
func CommitTransaction(tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// RollbackTransaction rolls back a transaction
func RollbackTransaction(tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}
