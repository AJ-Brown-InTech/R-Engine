package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func InitDatabase(conn string) (*sqlx.DB, error) {
	//connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disabled",DBHost, DBPort, DBUser, DBPassword, DBName)
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


// Single record read
func SingleReadTransaction(db *sqlx.DB, dest interface{}, query string,  args ...any) error {
	err := db.Get(dest, query, args...)
	if err != nil{ 
		return err
	}
	return nil
}

// Read record for batch transaction
func BatchReadTransaction(tx sqlx.Tx, query string, dest, args []any)  error {
	err := tx.Select(dest, query, args...)
	if err != nil {
		txer := tx.Rollback()
		if txer != nil {
			return errors.Join(err,txer)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		txer := tx.Rollback()
		if txer != nil {
			return errors.Join(err,txer)
		}
		return err
	}
	return nil
}

 // Update record for batch transaction
func UpdateTransaction(ctx context.Context, tx *sqlx.Tx, table string, obj interface{}, conditionColumn string, conditionValue interface{}) error {
    columns, args, err := ObjectToDatabase(obj)
    if err != nil {
        return err
    }

    // Check if the condition column is in the object
    valid := false
    for _, v := range columns {
        if conditionColumn == v {
            valid = true
            break
        }
    }

    if !valid {
        return errors.New("condition column provided is not valid")
    }

    setClause := make([]string, len(columns))
    for i, column := range columns {
        setClause[i] = fmt.Sprintf("%s = $%d", column, i+1)
    }

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", table, strings.Join(setClause, ", "), conditionColumn, len(columns)+ 1)

    _, err = tx.ExecContext(ctx, query, append(args, conditionValue)...)
    if err != nil {
        txer := tx.Rollback()
        if txer != nil {
            return errors.Join(err, txer)
        }
        return err
    }

    err = tx.Commit()
    if err != nil {
        txer := tx.Rollback()
		if txer != nil {
			return errors.Join(err, txer)
		}
		return err
    }
    return nil
}

// Insert record with trasaction
func InsertTransaction(ctx context.Context, tx *sqlx.Tx, table string, obj interface{}) error {
	columns, args, err := ObjectToDatabase(obj)
	if err != nil {
		return err
	}

	// Convert slices to pq.Array if needed
	for i, arg := range args {
		if slice, isSlice := arg.([]string); isSlice {
			args[i] = pq.Array(slice)
		}
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err = tx.ExecContext(ctx, query, args...)
	
	if err != nil {
		return err
	}

	return nil
}

// Delete record for transaction
func DeleteTransaction(ctx context.Context, tx *sqlx.Tx, table string, conditionColumn string, conditionValue interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, conditionColumn)

	_, err := tx.ExecContext(ctx, query, conditionValue)
	if err != nil {
		txer := tx.Rollback()
		if txer != nil {
			return errors.Join(err, txer)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		txer := tx.Rollback()
		if txer != nil {
			return errors.Join(err, txer)
		}
		return err
	}

	return nil
}
